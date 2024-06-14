import json
import traceback
import uuid
from threading import Thread
from flask import Flask,jsonify,make_response
import redis
import logging
from logging.handlers import HTTPHandler
from typing import List
from dynaconf import Dynaconf
import os
import requests
import time
def init_logger(pod,task_id,node_name):
    logger=logging.getLogger(f"{pod}/{task_id}")
    logger.setLevel(logging.DEBUG)
    fmt=logging.Formatter(fmt='%(asctime)s %(levelname)s %(module)s %(funcName)s [%(lineno)d] %(message)s',datefmt='%Y-%m-%d %H:%M:%S')
    cerebellum_addr="cerebellum-svc-"+node_name+":3000"
    http=HTTPHandler(cerebellum_addr,"/api/v1/report/log","POST",False)
    http.setLevel(logging.INFO)
    straem_handler=logging.StreamHandler()
    straem_handler.setFormatter(fmt)
    straem_handler.setLevel(logging.DEBUG)
    logger.addHandler(straem_handler)
    logger.addHandler(http)
    return logger



class BaseMetadata():
    def __init__(self,ts:float,equipment_id:str,message_id:str) -> None:
        self.ts=ts
        self.equipment_id=equipment_id
        self.message_id=message_id
class Message():
    def __init__(self,filepath:str="",metadata:List[BaseMetadata]=[],length:int=0) -> None:
        self.filepath=filepath
        self.metadata=metadata
        self.length=length
    def serialize(self)->str:
        m=[vars(t) for t in self.metadata]
        return json.dumps({"file":self.filepath,"metadata":json.dumps(m),"length":self.length})
    def deserialize(self,json_str:str):
        t=json.loads(json_str)
        self.filepath=t['file']
        self.metadata=[]
        meta=json.loads(t['metadata'])
        for i in meta:
            self.metadata.append(BaseMetadata(**i))
        self.length=t['length']
    def __str__(self) -> str:
        return f"File at {self.filepath} with length of {self.length}"



class RedisQueue():
    def __init__(self,host="redis-svc",port=6379,prefix="rawData",task_id="Default") -> None:
        self.r=redis.StrictRedis(host=host,port=port,db=1)
        self.prefix=prefix
        self.ps=self.r.pubsub()
        self.task_id=task_id
    
    def put(self,msg:Message,prefix=None):
        return self.r.publish(f"/{self.task_id}/{self.prefix[0] if not prefix else prefix}",msg.serialize())
    
    def put_all(self,msg:Message):
        for item in self.prefix:
            self.r.publish(f"/{self.task_id}/{item}",msg.serialize())
        return True

    def subscribe(self):
        for item in self.prefix:
            self.ps.subscribe(f"/{self.task_id}/{item}")
    
    def get_listener(self):
        return self.ps.listen()

class Settings():
    def __init__(self) -> None:
        settings=Dynaconf(load_dotenv=True)
        self.type=settings.type
        self.input=settings.input #list
        self.output=settings.output #list
        self.base_path=settings.base_path
        self.task_id=settings.task_id
        self.port=settings.port
        self.node_name=settings.node
        self.redis_host="redis-svc-"+self.node_name
        self.pod_name=settings.pod
        try:
            self.custom=json.loads(settings.custom)
        except:
            self.custom={}

class EdgeComputing:
    def __init__(self) -> None:
        self._settings=Settings()
        self.logger=init_logger(self._settings.type,self._settings.task_id,self._settings.node_name)
        self._receive=RedisQueue(host=self._settings.redis_host,prefix=self._settings.input,task_id=self._settings.task_id)
        self._receive.subscribe()
        self._output=RedisQueue(host=self._settings.redis_host,prefix=self._settings.output,task_id=self._settings.task_id)
        self._t=Thread(target=self._listen)
        self.custom=self._settings.custom
        self._app=Flask(__name__)
        self._app.add_url_rule("/api/ping","_ping",self._ping)
    def _on_trace(self,type,message_id):
        cerebellum_addr="http://cerebellum-svc-"+self._settings.node_name+":3000"
        Thread(target=requests.post,args=(cerebellum_addr+"/api/v1/report/trace",),kwargs={"json":{"type":type,"message_id":message_id,"task":self._settings.task_id,"pod":self._settings.pod_name,"time":time.time_ns()}}).start()
    def _on_receive(self,i):
        if i['type']=="subscribe" or i['type']=="pong":
            return
        msg=Message()
        msg.deserialize(i['data'].decode())
        self.logger.debug("获取数据成功:" + str(msg))
        self._on_trace("receive",msg.metadata[0].message_id)
        fileout=f"{self._settings.base_path}{str(uuid.uuid4())}.csv"
        try:
            self.forward(filein=msg.filepath,fileout=fileout)
            msg_out=Message(filepath=fileout,metadata=msg.metadata,length=msg.length)
            self._output.put(msg_out)
            self._on_trace("publish",msg.metadata[0].message_id)
            self._after_publish_hook(msg.filepath,fileout)
        except:
            self.logger.error("故障:"+traceback.format_exc())
    def _after_publish_hook(self,filein,fileout):
        self.logger.debug("running after publish hook")
        try:
            os.remove(filein)
        except:
            self.logger.error("故障:"+traceback.format_exc())

    def _listen(self):
        for i in self._receive.get_listener():
            self._on_receive(i)
    def _ping(self):
        if self._t.is_alive():
            return jsonify(status=0,msg='ok',data={"type":self._settings.type,"task":self._settings.task_id,"input":self._settings.input,"output":self._settings.output,"base_path":self._settings.base_path})
        else:
            return make_response(jsonify(state=-1,msg="thread dead"),502)
    def forward(self,filein:str,fileout:str):
        self.logger.warning("forward function not implemented, use a dummy one")
        self.logger.debug("input:%s,output:%s",filein,fileout)
    def run(self):
        self._t.start()
        self._app.run(host="0.0.0.0",port=self._settings.port,debug=False)
    
