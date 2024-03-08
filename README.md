# PilotGo

#### 介绍

PilotGo是一个openEuler社区原生的运维管理平台。

#### 软件架构
开发工具：golang 1.15

系统支持：openEuler、麒麟操作系统

PilotGo项目后端采用golang语言开发，使用到以下开源库：

​        web框架：https://github.com/gin-gonic/gin

​        websocket：https://github.com/gorilla/websocket

​        日志框架：https://github.com/sirupsen/logrus

​        文件监控：https://github.com/fsnotify/fsnotify

​        配置解析：https://github.com/spf13/viper

​        mock测试：https://github.com/golang/mock

以及系统自带库：

​        net/http

​        os

​        time 等

前端代码主要使用到以下技术：

​        JavaScript技术：https://www.javascript.com

​        vue框架：https://cn.vuejs.org

​        element组件：https://element.eleme.cn

可在该网站方便查询第三方库及系统库的API文档：

​        https://pkg.go.dev/

​        https://pkg.go.dev/std

#### 安装、启动教程

源码编译部署安装
```bash
# Download pilotgo source code：
    git clone https://gitee.com/openeuler/PilotGo.git
# Quick build:
    cd PilotGo
    chmod +x build.sh
    ./build.sh
# Quick install:
    cd /root
    tar -xzvf pilotgo-xxx.tar.gz
    cd pilotgo-xxx
    chmod +x install.sh
    ./install.sh
# Go into /opt/PilotGo
    cd /opt/PilotGo
# Modify configuration file：
   vim config_server.yaml
   vim config_agent.yaml
# Warn: There are two options for log-driver in config_server.yaml or config_agent.yaml.
    stdout: Terminal console output log
    file: Output log to specified file
# Start-up
    server:
       nohup ./server &
    agent:
       nohup ./agent &
```

二次开发部署
```bash
# Required before Pilotgo application deployment：
go >=1.17;  nodejs >=14

# Front end deployment：
1. Set NPM source address to Taobao source
    npm install -g cnpm --registry=https://registry.npm.taobao.org
2. Download npm dependency packages
    cnpm install
3. Modify the basepath to the server address and port under config/index.js, and run it directly：
    npm run dev

# Server deployment：
1. Rename the config_server.yaml.templete to config_server.yaml, and configuration
2. Rename the config_agent.yaml.templete to config_agent.yaml, and configuration
3. PilotGo server with hot reload at localhost:8888
    go run pkg/app/server/main.go
4. PilotGo agent with hot reload at localhost:8879
    go run pkg/app/agent/main.go
```
登录首页
```bash
# web页面访问  ip:8888
初始用户：admin@123.com
密码：123456
```
![](./docs/images/login.png)

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
