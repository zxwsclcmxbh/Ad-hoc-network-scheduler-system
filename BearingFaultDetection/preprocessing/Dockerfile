FROM python:3.8
RUN mkdir /workspace
WORKDIR /workspace
COPY . /workspace
RUN pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
ENTRYPOINT [ "python3","main.py" ]