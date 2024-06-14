from EdgeSDK import EdgeComputing,RedisQueue
import uuid
import os
import pandas
class Compute(EdgeComputing):
    ''' 继承自SDK '''
    def __init__(self) -> None:
        ''' 
        初始化模型、初始化环境等工作，必须添加 super().__init__()
        '''
        super().__init__()
    def forward(self, filein: str, fileout: str):
       df=pandas.read_csv(filein,header=None)
       assert df.shape[0]==1 or df.shape[1]==1 , "wrong csv format"
       df=df.transpose()
       df.to_csv(fileout,header=None,index=False)
# 运行
Compute().run()