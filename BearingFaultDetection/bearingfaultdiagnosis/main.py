import pickle
from EdgeSDK import EdgeComputing
import pandas as pd
class Compute(EdgeComputing):
    ''' 继承自SDK '''
    def __init__(self) -> None:
        ''' 
        初始化模型、初始化环境等工作，必须添加 super().__init__()
        '''
        super().__init__()
        with open('paderborn.pickle', 'rb') as f:
            self.model = pickle.load(f)
    def forward(self, filein: str, fileout: str):
        '''
        运算函数，执行运算操作，传入输入输出文件参数
        可使用 self.logger 进行日志、使用self.custom 读取自定义配置项
        '''
        test = pd.read_csv(filein)
        test_x = test.iloc[:, :43].values
        res_predict = self.model.predict(test_x)
        data_df = pd.DataFrame(res_predict)
        data_df.to_csv(fileout, index=False, float_format='%.1f')
        self.logger.debug("The result is:" + str(data_df))
# 运行
Compute().run()