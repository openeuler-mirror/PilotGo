# 安装部署PilotGo

## 1. 确认服务器架构
确认将要部署PilotGo的服务器的架构，选择合适的tar.gz进行解压。
```bash
$ tar -xzvf pilotgo-1.0.0.xxx.tar.gz
$ cd pilotgo-1.0.0
```
## 2. 一键安装
执行脚本：
```bash
$ ./install.sh
```
### 注意：
如果执行sh脚本出现 : 没有那个文件或目录  
解决方案:

1.执行
```bash
$ vim build.sh
```
2.不点击 `i` 修改文件的情况下,执行
```bash
$ :set ff=unix 
```
3.再执行
```bash
$ :wq    #问题解决
```
## 3. 服务端部署
进入部署目录
```bash
$ cd /opt/PilotGo
``` 
修改yaml文件,按照实际部署情况进行配置：
```bash
$ vim config_server.yaml
``` 
### 注意：
在yaml中有log模块，driver字段是可选模式：  
stdout：输出到终端控制台；  
file：输出到path下的指定文件。   

```bash
# server端启动
$ ./server &
# web页面访问  ip:端口
初始用户：admin@123.com
密码：123456
```
![](./images/login.png)

## 4. agent端部署
注册新机器流程（原则上只需配置`config_agent.yaml`和`agent`二进制即可注册，但是考虑到日志输出，建议先执行`./install.sh`）：  
进入部署目录
```bash
$ cd /opt/PilotGo
``` 
修改yaml文件,按照实际部署情况进行配置：
```bash
$ vim config_agent.yaml
``` 
### 注意：
在yaml中有log模块，driver字段是可选模式：  
stdout：输出到终端控制台；  
file：输出到path下的指定文件。   

```bash
# agent端启动
$ ./agent &
```