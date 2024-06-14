
import time
from util import RedisQueue, trace
from flask import request,jsonify
from util import settings,Message,BaseMetadata
import csv
import pandas as pd
import json
import uuid
from kafka import KafkaConsumer

def getHandleIncomingOne(output:RedisQueue):
    def handleIncomingOne():
        '''
        {"value":[],"equipment_id":"",ts:xxx.xx}
        '''
        r=request.json
        saveRoute = settings.base_path + str(uuid.uuid4())+".csv"
        data=newhttpmapper(r,settings.custom["mapper"])
        # print(data)
        if type(data['value'])!=list:
            value=[data['value']]
        else:
            value=data['value']
        with open(saveRoute,'w') as fout:
            csv.writer(fout).writerow(value)
        msg=Message(saveRoute,[BaseMetadata(data['ts'],data['equipment_id'],uuid.uuid4())],1)
        output.put(msg.serialize())
        trace("receive",msg.metadata[0].message_id)
        return jsonify(code=0, msg="uploadRawData successfully", data=None)
    return handleIncomingOne

def getHandleIncomingMany(output:RedisQueue):
    def handleIncomingMany():
        file = request.files['file']
        t=settings.custom.get("metadata_field_name")
        metadata_origin=request.form.get("metadata" if t is None else t)
        # print(metadata_origin,file)
        if file is None or metadata_origin is None:
            return jsonify(code=1, msg="failed to uploadRawData", data=None)
        data = pd.read_csv(file,header=None)
        print(data)
        metadata=json.loads(metadata_origin)
        saveRoute = settings.base_path + str(uuid.uuid4())+".csv"
        if len(metadata)!=len(data):
            # print(len(metadata),len(data))
            return jsonify(code=1, msg="failed to uploadRawData", data=None)
        metadata_after_map=[]
        first_id=""
        for index,i in enumerate(metadata):
            temp=newhttpmapper(i,settings.custom['mapper'])
            t_id=uuid.uuid4()
            if index==0:
                first_id=t_id
            temp["message_id"]=t_id
            metadata_after_map.append(temp)
        data.to_csv(saveRoute, sep=',', header=None, index=False)
        msg=json.dumps({"file":saveRoute,"metadata":json.dumps(metadata_after_map),"length":len(data)})
        output.put(msg)
        trace("receive",first_id)
        return jsonify(code=0, msg="uploadRawData successfully", data=None)
    return handleIncomingMany


def getHandleIncomingKafka(output:RedisQueue):
    def handleIncomingOne():
        '''
        {"value":[],"equipment_id":"",ts:xxx.xx}
        '''
        consumer = KafkaConsumer(settings.custom['kafka-topic'],
                         bootstrap_servers=[settings.custom['kafka']],
                         value_deserializer=lambda m: json.loads(m.decode()))
        for message in consumer:
            saveRoute = settings.base_path + str(uuid.uuid4())+".csv"
            data=newkafkamapper(message,settings.custom["mapper"])
            if type(data['value'])!=list:
                value=[data['value']]
            else:
                value=data['value']
            with open(saveRoute,'w') as fout:
                csv.writer(fout).writerow(value)
            msg=Message(saveRoute,[BaseMetadata(data['ts'],data['equipment_id']+"-"+data["station_id"],uuid.uuid4())],1)
            output.put(msg.serialize())
            trace("receive",msg.metadata[0].message_id)
    return handleIncomingOne

def mapper(original:dict,mapconf:dict)->dict:
    result={}
    for k,v in mapconf.items():
        t=original.get(k)
        if t:
            result[v["name"]]=t
        else:
            if v['default']=="ts":
                result[v["name"]]=time.time()
            else:
                result[v["name"]]=v["default"]
    return result

def kafkaMapper(kafkaMessage,mapconf:dict)->dict:
    result={}
    for k,v in mapconf.items():
        t=kafkaMessage.value.get(k)
        if t:
            result[v["name"]]=t
        else:
            if v['default']=="kafkats":
                result[v['name']]=kafkaMessage.timestamp/1000
            elif v['default']=="kafkatopic":
                result[v['name']]=kafkaMessage.topic
            else:
                result[v["name"]]=v["default"]
    return result

"""
"mapper": {
        "ts": { "type": "kafkats", "value": "" },
        "equipment_id": { "type": "kafkatopic", "value": "" },
        "value": { "type": "custom", "value": "test" }
      },
"""

def newkafkamapper(kafkaMessage,mapconf:dict)->dict:
    result={}
    for k,v in mapconf.items():
        if v['type']=="custom":
            result[k]=kafkaMessage.value.get(v['value'])
        elif v['type']=="kafkats":
            result[k]=kafkaMessage.timestamp/1000
        elif v['type']=="kafkatopic":
            result[k]=kafkaMessage.topic
        elif v['type']=="ts":
            result[k]=time.time()
        else:
            result[k]=v['value']
    return result
def newhttpmapper(msg:dict,mapconf:dict)->dict:
    result={}
    # print(mapconf)
    for k,v in mapconf.items():
        if v['type']=="custom":
            result[k]=msg.get(v['value'])
        elif v['type']=="ts":
            result[k]=time.time()
        else:
            result[k]=v['value']
    return result

