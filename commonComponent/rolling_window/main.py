import os
import traceback
import uuid
from EdgeSDK import EdgeComputing,Message
from collections import defaultdict
import csv
class window:
    def __init__(self,window_length:int,slide_step:int) -> None:
        if slide_step>window_length:
            raise Exception("slide_step:%d is larger than window_length:%d"%(slide_step,window_length))
        self.data=[]
        self.window_length=window_length
        self.slide_step=slide_step
        self.slide_counter=0
    def append(self,element)-> bool:
        if len(self.data)==self.window_length:
            self.data.pop(0)
            self.slide_counter+=1
        self.data.append(element)
        if self.slide_counter==self.slide_step:
            self.slide_counter=0
            return True
        elif self.slide_counter==0 and len(self.data)==self.window_length:
            return True
        else:
            return False

def get_window(length,slide):
    def return_window():
        return window(length,slide)
    return return_window

class Compute(EdgeComputing):
    ''' 继承自SDK '''
    def __init__(self) -> None:
        ''' 
        初始化模型、初始化环境等工作，必须添加 super().__init__()
        '''
        super().__init__()
        print(self._settings.base_path)
        self.queue=defaultdict(get_window(int(self.custom['window_length']),int(self.custom['slide_step'])))
    def _on_receive(self,incoming):
        if incoming['type']=="subscribe" or incoming['type']=="pong":
            return
        try:
            msg=Message()
            msg.deserialize(incoming['data'].decode())
            self.logger.debug("获取数据成功:" + msg.filepath)
            self._on_trace("receive",msg.metadata[0].message_id)
            with open(msg.filepath,'r') as fin:
                reader=csv.reader(fin)
                for i,m in zip(reader,msg.metadata):
                    if self.queue[m.equipement_id].append((i,m)):
                        fileout=f"{self._settings.base_path}{str(uuid.uuid4())}.csv"
                        with open(fileout,'w') as fout:
                            writer=csv.writer(fout)
                            writer.writerows([d for d,_ in self.queue[m.equipement_id].data])
                        if self.custom['metadata_aggregation']=="one":
                            msgout=Message(fileout,[self.queue[m.equipement_id].data[0][1]],1)
                        else:
                            msgout=Message(fileout,[metadata for _,metadata in self.queue[m.equipement_id].data],len(self.queue[m.equipement_id].data))
                        self._output.put(msgout)
                        self._on_trace("publish",(self.queue[m.equipement_id].data[0][1]).message_id)
                        self.logger.debug("导出数据成功:" + fileout)
                        
                    else:
                        pass
                        # self.logger.debug("数据数量不够")
        except:
            self.logger.error("故障:"+traceback.format_exc())
        self._after_publish_hook(msg.filepath,None)
    def _after_publish_hook(self,filein,fileout):
        self.logger.debug("running after publish hook")
        try:
            os.remove(filein)
        except:
            self.logger.error("故障:"+traceback.format_exc())
Compute().run()





    

