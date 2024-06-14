import csv
import json
import uuid
from util import BaseMetadata, Message
from util import RedisQueue,init_logger,settings
from flask import Flask,jsonify,request
import pandas as pd
logger=init_logger("intake",settings.task_id,settings.node_name)
output=RedisQueue(host=settings.redis_host,prefix=settings.output,task_id=settings.task_id)
app=Flask(__name__)


@app.route('/api/uploadRawDataMany',methods=['POST'])
def handleIncomingMany():
    file = request.files['file']
    metadata_origin=request.form.get("metadata")
    print(metadata_origin,file)
    if file is None or metadata_origin is None:
        return jsonify(code=1, msg="failed to uploadRawData", data=None)
    data = pd.read_csv(file,header=None)
    metadata=json.loads(metadata_origin)
    saveRoute = settings.base_path + str(uuid.uuid4())+".csv"
    if len(metadata)!=len(data):
        print(len(metadata),len(data))
        return jsonify(code=1, msg="failed to uploadRawData", data=None)
    data.to_csv(saveRoute, sep=',', header=None, index=False)
    msg=json.dumps({"file":saveRoute,"metadata":metadata_origin,"length":len(data)})
    output.put(msg)
    return jsonify(code=0, msg="uploadRawData successfully", data=None)

@app.route('/api/uploadRawDataOne',methods=['POST'])
def handleIncomingOne():
    '''
    {"value":[],"equipment_id":"",ts:xxx.xx}
    '''
    r=request.json
    saveRoute = settings.base_path + str(uuid.uuid4())+".csv"
    with open(saveRoute,'w') as fout:
        csv.writer(fout).writerow(r['value'])
    msg=Message(saveRoute,[BaseMetadata(r['ts'],r['equipment_id'])],1)
    output.put(msg.serialize())
    return jsonify(code=0, msg="uploadRawData successfully", data=None)

@app.route("/api/ping",methods=['GET'])
def ping():
    return jsonify(status=0,msg='ok',data={"type":settings.type,"task":settings.task_id,"input":settings.input,"output":settings.output,"base_path":settings.base_path})

app.run(host="0.0.0.0",port=settings.port,debug=False)