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

```
# Required before startup
go >=1.15;  nodejs >=14
# install npm dependencies
npm install
# vue server with hot reload at localhost:8080(Modify the IP addresses of go server and web server under config/index.js)
npm run dev
# Rename the config.yaml.templete to config.yaml

# PilotGo server with hot reload at localhost:8083
go run pkg/app/server/main.go
# PilotGo agent with hot reload at localhost:8083
go build -o agent pkg/app/agent/main.go
./agent

```

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
