FROM node:alpine

RUN apk update && apk add --no-cache python3 py3-pip

WORKDIR /root/

COPY ./ /root/

RUN pip3 install flask && npm install web3@1.10.1 -g && npm install fs-extra -g

EXPOSE 5001

# run application
CMD ["python3", "app.py"]