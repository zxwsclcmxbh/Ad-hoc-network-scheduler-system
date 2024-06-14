from threading import Thread
from util import BaseMetadata, Message
from util import RedisQueue,init_logger,settings
from flask import Flask,jsonify,request
import pandas as pd
import func
logger=init_logger("intake",settings.task_id,settings.node_name)
output=RedisQueue(host=settings.redis_host,prefix=settings.output,task_id=settings.task_id)
app=Flask(__name__)


@app.route("/api/ping",methods=['GET'])
def ping():
    return jsonify(status=0,msg='ok',data={"type":settings.type,"task":settings.task_id,"input":settings.input,"output":settings.output,"base_path":settings.base_path})

if settings.custom['mode']=="http-file":
    app.add_url_rule("/api/uploadRawDataMany",view_func=func.getHandleIncomingMany(output),methods=["POST"],endpoint="many")
elif settings.custom['mode']=="http-json":
    app.add_url_rule("/api/uploadRawDataOne",view_func=func.getHandleIncomingOne(output),methods=["POST"],endpoint="one")
elif settings.custom['mode']=="kafka":
    t=Thread(target=func.getHandleIncomingKafka(output))
    t.start()
app.run(host="0.0.0.0",port=settings.port,debug=False)