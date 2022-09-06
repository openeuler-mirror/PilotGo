#!/bin/bash

echo "thanks for choosing PilotGo"
echo "------------------- start to install -------------------"

echo "creating directories ..."
# 创建项目部署目录和日志目录
if [ ! -d "/opt/PilotGo/server" ];then
    mkdir -p /opt/PilotGo/server
    # used for record monitor targets
    mkdir -p /opt/PilotGo/monitor
fi

# 复制文件
cp config_server.yaml /opt/PilotGo/server/config_server.yaml
cp pilotgo-server /opt/PilotGo/server/

echo "启动前准备完成，请前往/opt/PilotGo下执行二进制文件以启动PilotGo"
