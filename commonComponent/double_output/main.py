from EdgeSDK import EdgeComputing,RedisQueue
import uuid
import os
class Compute(EdgeComputing):
    ''' 继承自SDK '''
    def __init__(self) -> None:
        ''' 
        初始化模型、初始化环境等工作，必须添加 super().__init__()
        '''
        super().__init__()
        self.output2=RedisQueue(prefix=self.custom['output2'],task_id=self._settings.task_id)
    def forward(self, filein: str, fileout: str):
        fileout2=f"{self._settings.base_path}{str(uuid.uuid4())}.csv"
        os.system(f"cp {filein} {fileout}")
        os.system(f"cp {filein} {fileout2}")
        self.output2.put(fileout2)
# 运行
Compute().run()