import traceback
from EdgeSDK import EdgeComputing,RedisQueue,init_logger,Settings,Message,BaseMetadata
from threading import Thread
import requests
from flask import Flask
import pandas as pd
class Compute(EdgeComputing):
    ''' 继承自SDK '''
    def __init__(self) -> None:
        ''' 
        初始化模型、初始化环境等工作，必须添加 super().__init__()
        '''
        self._settings=Settings()
        self.logger=init_logger(self._settings.type,self._settings.task_id,self._settings.node_name)
        self._receive=RedisQueue(host=self._settings.redis_host,prefix=self._settings.input,task_id=self._settings.task_id)
        self._t=Thread(target=self._listen)
        self.custom=self._settings.custom
        self._app=Flask(__name__)
        self._app.add_url_rule("/api/ping","_ping",self._ping)
        self.logger.info("start up success")
    def _on_receive(self,i):
        if i['type']=="subscribe" or i['type']=="pong":
            return
        msg=Message()
        msg.deserialize(i['data'].decode())
        self._on_trace("receive",msg.metadata[0].message_id)
        self.logger.debug("获取数据成功:" + str(msg))
        try:
            self.send(msg)
            self._on_trace("publish",msg.metadata[0].message_id)
        except:
            self.logger.error("故障:"+traceback.format_exc())
        self._after_publish_hook(msg.filepath,None)
    def send(self,msg:Message):
        df=pd.read_csv(msg.filepath)
        data=df.to_dict("records")
        result=[]
        print(len(msg.metadata),len(data))
        for m,d in zip(msg.metadata,data):
            result.append(
                {
                    "timestamp":str(m.ts),
                    "values":d
                }
            )
        req={
            "equipment_id":msg.metadata[0].equipement_id,
            "task_id":self._settings.task_id,
            "items":result
        }
        r=requests.post(f"http://cerebellum-svc-{self._settings.node_name}:3000/api/v1/report/value",json=req)
        self.logger.debug(r.text)

# 运行
Compute().run()


#     def _on_receive(self,i):
#         if i['type']=="subscribe" or i['type']=="pong":
#             return
        
#         msg=Message()
#         msg.deserialize(i['data'].decode())
#         self.logger.debug("获取数据成功:" + msg)
#         fileout=f"{self._settings.base_path}{str(uuid.uuid4())}.csv"
#         try:
#             self.forward(filein=msg.filepath,fileout=fileout)
#             msg_out=Message(filepath=fileout,metadata=msg.metadata,length=msg.length)
#             self._output.put(msg_out)
#             self._after_publish_hook(msg.filepath,fileout)
#         except:
#             self.logger.error("故障:"+traceback.format_exc())
#     def _after_publish_hook(self,filein,fileout):
#         self.logger.debug("running after publish hook")
#         try:
#             os.remove(filein)
#         except:
#             self.logger.error("故障:"+traceback.format_exc())

#     def _listen(self):
#         for i in self._receive.subscribe():
#             self._on_receive(i)
#     def _ping(self):
#         if self._t.is_alive():
#             return jsonify(status=0,msg='ok',data={"type":self._settings.type,"task":self._settings.task_id,"input":self._settings.input,"output":self._settings.output,"base_path":self._settings.base_path})
#         else:
#             return make_response(jsonify(state=-1,msg="thread dead"),502)
#     def forward(self,filein:str,fileout:str):
#         self.logger.warning("forward function not implemented, use a dummy one")
#         self.logger.debug("input:%s,output:%s",filein,fileout)
#     def run(self):
#         self._t.start()
#         self._app.run(host="0.0.0.0",port=self._settings.port,debug=False)
    
