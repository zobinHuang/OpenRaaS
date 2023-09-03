#!/usr/bin/env bash 


cd ../dockertools

# 创建镜像
# build-wine.sh

# 删除编号为 1 的 dcwine
# ./rm-wine.sh 1


cd ..

# 用 webdav 为编号为 1 的 instance 挂载远端磁盘
expect ./auto-mount.exp 1 davfs 192.168.10.189:7189 /public_hdd/game/PC/dcwine kb109 ******

# 开启编号为 1 的 dcwine，并指定游戏为 spider
sh ./run-wine.sh dcwine 1 /apps/spider sol.exe spider game 480 320 192.168.10.138 30 h.264
# sh ./run-wine.sh dcwine 1 /apps/gputest GpuTest.exe gputest game 800 600 192.168.10.138 30 h.264
# sh ./run-wine.sh dcwine 1 /apps/PotPlayer PotPlayerMini64.exe videoplayer game 800 600 192.168.10.138 30 h.264 test.flv

cd ./networktools

# 设置 docker 的网络特性
./set-dockernet.sh