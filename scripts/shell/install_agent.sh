#!/bin/bash

echo "thanks for choosing PilotGo"
echo "------------------- start to install -------------------"

echo "creating directories ..."
# 创建项目部署目录和日志目录
if [ ! -d "/opt/PilotGo/agent" ];then
    mkdir -p /opt/PilotGo/agent
fi

# 复制文件
cp config_agent.yaml /opt/PilotGo/agent/config_agent.yaml
cp pilotgo-agent /opt/PilotGo/agent/

echo "启动前准备完成，请前往/opt/PilotGo下执行二进制文件以启动PilotGo"