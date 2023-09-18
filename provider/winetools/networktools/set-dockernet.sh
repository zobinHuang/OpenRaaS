#!/usr/bin/bash

# 该脚本的执行会阻塞终端

. ./dockernet.cfg

# pumba netem --duration 5m --tc-image gaiadocker/iproute2 delay --time 30 --jitter 0 --distribution normal provider-streamer-1
# sudo tc qdisc add dev enp6s18 root netem delay 30ms
# sudo tc qdisc change dev enp6s18 root netem delay 30ms

for ((i=1; i<=${docker_number}; i++))
do
    index=${i}-1
    # 延迟
    pumba netem --duration ${duration} --tc-image gaiadocker/iproute2 delay --time ${delay_time}${index} --jitter ${delay_jitter}${index} --distribution ${delay_distribution}${index} appvm${i} &
    # 丢包
    pumba netem --duration ${duration} --tc-image gaiadocker/iproute2 loss --percent ${loss_percent}${index} appvm${i} &
    # 带宽
    pumba netem --duration ${duration} --tc-image gaiadocker/iproute2 rate --rate ${rate}${index} appvm${i} &
done