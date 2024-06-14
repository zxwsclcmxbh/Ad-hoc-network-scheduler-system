import time
from kafka import KafkaProducer
from kafka.errors import KafkaError
import json
producer = KafkaProducer(bootstrap_servers=['cjl.hfzhang.wang:9092'],value_serializer=lambda m: json.dumps(m).encode())
print(producer.bootstrap_connected())
future=producer.send("test",{"duration":[1,2,3],"eid":""})
producer.flush()
try:
    record_metadata = future.get(timeout=10)
except KafkaError as e:
    # Decide what to do if produce request failed...
    print(e)
    pass

# Successful result returns assigned partition and offset
print (record_metadata.topic)
print (record_metadata.partition)
print (record_metadata.offset)

