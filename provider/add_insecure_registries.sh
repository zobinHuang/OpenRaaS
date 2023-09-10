#!/bin/bash

# 指定要添加的不安全 Registry 地址，将它们存储在数组中
INSECURE_REGISTRIES=("192.168.0.109:10960" "192.168.182.156:10960" "192.168.0.189:10960" "192.168.222.216:10960")

# 检查是否已经存在 daemon.json 文件
if [ ! -f /etc/docker/daemon.json ]; then
    echo "{}" > /etc/docker/daemon.json
fi

for REGISTRY in "${INSECURE_REGISTRIES[@]}"; do
    # 检查是否已经存在 insecure-registries 配置项
    if grep -q 'insecure-registries' /etc/docker/daemon.json; then
        # 如果已经存在，添加新指定的值
        sed -i "s|\"insecure-registries\":.*|\"insecure-registries\": [\n    \"$REGISTRY\",|" /etc/docker/daemon.json
    else
        # 如果不存在，添加配置项
        jq '. += {"insecure-registries": ["'$REGISTRY'"]}' /etc/docker/daemon.json > ~/daemon.json.tmp
        mv ~/daemon.json.tmp /etc/docker/daemon.json
    fi
    echo "已添加 $REGISTRY 到 /etc/docker/daemon.json"
done

# 重启 Docker 守护程序以应用更改
systemctl restart docker