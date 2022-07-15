#!/bin/bash

RED_COLOR='\E[1;31m'  #红
GREEN_COLOR='\E[1;32m' #绿
YELLOW_COLOR='\E[1;33m' #黄
RES='\E[0m'

echo -e "${GREEN_COLOR}感谢您选择PilotGo${RES}"
echo "-------------------开始执行部署安装-------------------"
echo "开始创建部署安装目录..."

# 创建项目部署目录和日志目录
if [ ! -d "/opt/PilotGo/" ];then
    mkdir /opt/PilotGo/
    mkdir /opt/PilotGo/log
else
    rm -rf /opt/PilotGo/*
    mkdir /opt/PilotGo/log
fi

# 创建prometheus targets文件目录
if [ ! -d "/opt/PilotGo/monitor_target/" ];then
    mkdir /opt/PilotGo/monitor_target/
fi

# 创建prometheus 配置文件目录
if [ ! -d "/etc/prometheus/" ];then
    mkdir /etc/prometheus/
fi

# 复制文件
cp config_agent.yaml /opt/PilotGo/config_agent.yaml
cp config_server.yaml /opt/PilotGo/config_server.yaml
cp alert.rules /etc/prometheus/alert.rules
cp server /opt/PilotGo/
cp agent /opt/PilotGo/

echo -e "${GREEN_COLOR}启动前准备完成，请前往/opt/PilotGo下执行二进制文件以启动PilotGo${RES}"