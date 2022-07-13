# 通过源码构建PilotGo

## 1. 检查golang、nodejs构建环境
为了构建PilotGo，需保证已经安装了golang语言环境和nodejs软件。golang的推荐版本为1.17.0及其之后的版本,nodejs的推荐版本为14.0及其之后的版本，否则编译可能失败。


## 2. 前端资源编译打包
为了保证后端推送前端告警功能的正常运行，建议将根目录下config/index.js下的basePath修改为可正常访问的后端ip和端口后再打包。  
前端资源编译打包：
```bash
$ npm run build
```

## 3. 服务端构建
在项目根目录下找到pilotgo.sh脚本，使用`chmod +x pilotgo.sh`赋予该脚本执行权限。  
运行脚本：
```bash
$ ./pilotgo.sh
```
脚本执行过程中会检查golang、nodejs是否满足编译和执行的版本要求。脚本执行完毕，在/opt/PilotGo目录下存在`agent`和`server`的二进制执行文件，以及`config_agent.yml`、`config_server.yml`以及生成`log`等文件夹。用户只需按照服务器环境和对应需求按实修改yml配置文件内容即可。
```bash
# server端启动
$ ./server
# agent端启动
$ ./agent
# web页面访问  ip:端口
初始用户：admin@123.com
密码：123456
```
![](./images/login.png)