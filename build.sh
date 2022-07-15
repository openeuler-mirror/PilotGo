#!/bin/bash

RED_COLOR='\E[1;31m'  #红
GREEN_COLOR='\E[1;32m' #绿
YELLOW_COLOR='\E[1;33m' #黄
RES='\E[0m'

echo -e "${GREEN_COLOR}感谢您选择PilotGo${RES}"
echo "正在检查前端编译环境..."

# 判断是否安装了NodeJS
if ! type node >/dev/null 2>&1; then
    yum install -y nodejs
    NodeJS=`node -v | grep -oP '\d*\.\d*.\d+'`
    if [ ${NodeJS: 0: 2} -lt 14 ]; then
        echo -e "${YELLOW_COLOR}警告: nodejs版本略低,请升级到 14.0 以上${RES}"
        exit 1
    fi
    exit 1
else
    NodeJS=`node -v | grep -oP '\d*\.\d*.\d+'`
    if [ ${NodeJS: 0: 2} -lt 14 ]; then
        echo -e "${YELLOW_COLOR}警告: nodejs版本略低,请升级到 14.0 以上${RES}"
        exit 1
    fi

    # 判断是否安装了NPM
    if ! type npm >/dev/null 2>&1; then
        echo -e "${RED_COLOR}错误: 请先安装 6.0 以上版本的npm${RES}"
        exit 1;
    fi
fi

echo -e "${GREEN_COLOR}前端运行依赖检查完成${RES}\n"
echo "正在检查服务端编译环境..."

# 判断是否安装了golang
if ! type go >/dev/null 2>&1; then
    yum install -y golang
    GoLang=`go version |awk '{print $3}' | grep -oP '\d*\.\d*.\d+'`
    if [ ${GoLang: 2: 2} -lt 17 ]; then
        echo -e "${YELLOW_COLOR}警告: golang版本略低,请升级到 1.17.0 以上${RES}"
        exit 1
    fi
    exit 1
else
    GoLang=`go version |awk '{print $3}' | grep -oP '\d*\.\d*.\d+'`
    if [ ${GoLang: 2: 2} -lt 17 ]; then
        echo -e "${YELLOW_COLOR}警告: golang版本略低,请升级到 1.17.0 以上${RES}"
        exit 1
    fi
fi
echo -e "${GREEN_COLOR}服务端运行依赖检查完成${RES}\n"

echo -e "${GREEN_COLOR}正在下载前端依赖，请稍候...${RES}"
npm install -g cnpm --registry=https://registry.npm.taobao.org
cnpm install
echo -e "\n${GREEN_COLOR}正在编译前端资源，请稍候...${RES}"
npm run build

cp ./dist/index.html ./resource/index.html
cp ./dist/static/pilotgo.ico ./resource/pilotgo.ico
cp -r ./dist/static/css ./resource/
cp -r ./dist/static/fonts ./resource/
cp -r ./dist/static/img ./resource/
cp -r ./dist/static/js ./resource/

echo "正在编译server端二进制文件..."
go build -o server ./pkg/app/server/main.go

echo "正在编译agent端二进制文件..."
go build -o agent pkg/app/agent/main.go

# 创建文件包路径
if [ ! -d "./pilotgo-1.0.0" ];then
    mkdir ./pilotgo-1.0.0
else
    rm -rf ./pilotgo-1.0.0/*
fi

# 复制配置文件
cp server ./pilotgo-1.0.0/
cp agent ./pilotgo-1.0.0/
cp config_agent.yaml.templete ./pilotgo-1.0.0/config_agent.yaml
cp config_server.yaml.templete ./pilotgo-1.0.0/config_server.yaml
cp alert.rules.templete ./pilotgo-1.0.0/alert.rules
cp install.sh ./pilotgo-1.0.0/

echo "正在压缩文件..."
tar -czvf pilotgo-1.0.0.$(arch).tar.gz ./pilotgo-1.0.0/*
rm -rf ./pilotgo-1.0.0
mv pilotgo-1.0.0.$(arch).tar.gz /root/

echo -e "${GREEN_COLOR}编译打包完成${RES}"
