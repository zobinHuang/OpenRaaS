FROM python:3.9.1
RUN apt-get update

RUN mkdir -p /daemon
WORKDIR /daemon

ADD ./requirements.txt /daemon/requirements.txt
RUN pip install -r requirements.txt
