#!/bin/bash

# run in docker: cpu version

# docker run -v $PWD:/workspace -w /workspace daisukekobayashi/darknet:darknet_yolo_v4_pre-cpu \
# darknet detector train ./data ./yolov2-tiny.cfg

# To use OpenRaaS, fist you should tag the image in the form of `192.168.0.109:10960/darknet:darknet_yolo_v4_pre-cpu` and then push

# 从 macOS 挂载，直接在 backup 文件夹写入会有 “Input/output error”，所以先在上级目录（“..”）写入，完成后再 cp
darknet detector train ./data ./yolov4-tiny.cfg

# 默认 100 步保存一次 weights
# 在 cfg 设置 max_batches=100 即可在第一次保存梯度时完成测试

cp ../yolov4-tiny_last.weights ./backup/yolov4-tiny_last.weights
sleep 1
# 从 macos 上挂载，第一次写入会有 “Input/output error”，所以写入两次 
cp ../yolov4-tiny_last.weights ./backup/yolov4-tiny_last.weights

echo "开始等待"
sleep 10
echo "完成"
