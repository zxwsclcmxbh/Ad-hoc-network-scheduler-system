import json
from kafka import KafkaConsumer
consumer = KafkaConsumer('tt1-current',
                         bootstrap_servers=['cjl.hfzhang.wang:9092'],
                         value_deserializer=lambda m: json.loads(m.decode()))
for message in consumer:
    # message value and key are raw bytes -- decode if necessary!
    # e.g., for unicode: `message.value.decode('utf-8')`
    print ("%s:%d:%d: key=%s value=%s ts=%s" % (message.topic, message.partition,
                                          message.offset, message.key,
                                          message.value,message.timestamp))