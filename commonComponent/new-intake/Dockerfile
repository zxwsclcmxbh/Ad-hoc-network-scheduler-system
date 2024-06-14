FROM --platform=$TARGETPLATFORM 10.112.245.161:5000/iai/python-base:3.7
RUN mkdir /workspace
WORKDIR /workspace
COPY . /workspace
ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" = "linux/arm/v7" ] ; then pip3 install -r requirements.txt  ; else pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple; fi
ENTRYPOINT [ "python3","main.py" ]