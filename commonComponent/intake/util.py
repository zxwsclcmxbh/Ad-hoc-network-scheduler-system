import json
from logging.handlers import HTTPHandler
import redis
import logging
from dynaconf import Dynaconf

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
    def __init__(self,ts:float,equipement_id:str) -> None:
        self.ts=ts
        self.equipement_id=equipement_id
class Message():
    def __init__(self,filepath:str,metadata:list[BaseMetadata],length:int) -> None:
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
    
    def put(self,msg:str):
        return self.r.publish(f"/{self.task_id}/{self.prefix}",msg)
    
    def subscribe(self):
        self.ps.subscribe(f"/{self.task_id}/{self.prefix}")
        return self.ps.listen()


class Settings():
    def __init__(self) -> None:
        settings=Dynaconf(load_dotenv=True)
        self.type=settings.type
        self.input=settings.input
        self.output=settings.output
        self.base_path=settings.base_path
        self.task_id=settings.task_id
        self.port=settings.port
        self.node_name=settings.node
        self.redis_host=settings.redis_host
settings=Settings()