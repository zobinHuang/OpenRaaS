#!/usr/bin/env bash

apt-get update
apt-get install qperf expect davfs2

# pumba 安装
# docker 网络控制工具
curl -L https://github.com/alexei-led/pumba/releases/download/0.9.0/pumba_linux_amd64 --output /usr/local/bin/pumba
chmod +x /usr/local/bin/pumba && pumba --help

# docker 安装
# 如果本地已有 docker 环境，不建议安装
curl -fsSL https://get.docker.com -o /tmp/get-docker.sh
sh /tmp/get-docker.sh

# Nvidia runtime 工具
apt-get install nvidia-container-runtime