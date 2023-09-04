FROM python:3.9.1
RUN apt-get update

RUN mkdir -p /daemon
WORKDIR /daemon

ADD . /daemon
RUN mkdir -p /storage
RUN pip install -r requirements.txt
