import traceback
from EdgeSDK import EdgeComputing
from sklearn.ensemble import IsolationForest
import numpy as np
import pandas as pd
import csv
class Compute(EdgeComputing):
    ''' 继承自SDK '''
    def __init__(self) -> None:
        ''' 
        初始化模型、初始化环境等工作，必须添加 super().__init__()
        '''
        super().__init__()
        self.model="123"
    def forward(self, filein: str, fileout: str):
        '''
        运算函数，执行运算操作，传入输入输出文件参数
        可使用 self.logger 进行日志、使用self.custom 读取自定义配置项
        '''
        self.logger.debug("fin:%s,fout:%s",filein,fileout)
        params = {'n_estimators': 100,
                'max_samples': 'auto',
                'contamination': 0.1,
                'max_features': 1.0,
                'path': filein,
                'opath': fileout
                }

        self.logger.debug("path:" + params['path'])
        self.logger.debug("opath:" + params['opath'])
        try:
            tn = pd.read_csv(params['path'],header=None)
            tn.dropna(inplace=True)
            train = np.array(tn)
            train_x = train[:, :-1]
            train_y = train[:, -1]
            np.array(train_y)

            train_x = np.array(train_x)
            clf = IsolationForest(n_estimators=params['n_estimators'],
                                max_samples=params['max_samples'],
                                contamination=params['contamination'],
                                max_features=params['max_features'],
                                bootstrap=False, n_jobs=1, random_state=None,
                                verbose=0).fit(train_x)
            pred = clf.predict(train_x)

            self.logger.debug("The number of raw data rows is:" + str(pred.size))
            df = pd.DataFrame(pd.read_csv(params['path'],header=None))[0:pred.size]
            df['pred'] = pred
            df2 = df[-df.pred.isin([-1])]
            df2 = df2.drop(['pred'], axis=1)
            data_out = df2.iloc[:, :].values
            csvfile2 = open(params['opath'], 'w', newline='')
            writer = csv.writer(csvfile2)
            m = len(data_out)
            self.logger.debug("The number of data rows after preprocessing is:" + str(m))
            for i in range(m):
                writer.writerow(data_out[i])
        except Exception as e:
            self.logger.error("get exception:"+traceback.format_exc())
# 运行
Compute().run()