#!/usr/bin/env bash

# 传入参数: 1.容器编号(第几个wine容器)

container_id=$1
container_name="appvm${container_id}"
docker rm -f ${container_name}