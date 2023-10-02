#!/bin/bash

# 本脚本仅用于三种工作节点的配置, 后台程序的配置需要自行修改: 
# 1. "./backstage/docker-compose.yaml"
# 2. "./nginx/conf.d/static.conf"
# 3. "./web/ant-client-page/src/Configurations/APIConfig.json"

# 注意: 要在本脚本所在目录执行, 且执行前自行配置目录下三种节点的配置内容

# 获取当前脚本所在目录的绝对路径
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 拷贝文件
# provider
cp "${script_dir}/provider_streamer_env" "${script_dir}/../provider/streamer/.env.dev"
cp "${script_dir}/provider_serverd_config" "${script_dir}/../provider/serverd/config"
# depository
cp "${script_dir}/scheduler_config.yaml" "${script_dir}/../depository/daemon/scheduler_config.yaml"
cp "${script_dir}/depository_config.yaml" "${script_dir}/../depository/daemon/config.yaml"
# filestore
cp "${script_dir}/scheduler_config.yaml" "${script_dir}/../filestore/daemon/scheduler_config.yaml"
cp "${script_dir}/filestore_config.yaml" "${script_dir}/../filestore/daemon/config.yaml"
cp "${script_dir}/app_scheduler_config.yaml" "${script_dir}/../filestore/app_online/scheduler_config.yaml"
cp "${script_dir}/app_config.yaml" "${script_dir}/../filestore/app_online/config.yaml"

echo "配置文件下发完毕"