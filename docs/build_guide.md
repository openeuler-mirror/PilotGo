# 通过源码构建PilotGo

## 1. 检查golang、nodejs构建环境
为了构建PilotGo，需保证已经安装了golang语言环境和nodejs软件。golang的推荐版本为1.17.0及其之后的版本,nodejs的推荐版本为14.0及其之后的版本，否则编译可能失败。


## 2. 前端资源编译打包
编译后端代码之前，需要先将前端编译，以将前端资源打包进后端二进制中，实现启动后端二进制的同时，前端页面也可正常访问。 


## 3. 服务端构建
在项目根目录下找到build.sh脚本，执行`chmod +x pilotgo.sh`赋予该脚本执行权限。  
运行脚本：
```bash
$ ./build.sh
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

脚本执行过程中会检查golang、nodejs是否满足编译和执行的版本要求；将前端和后端server、agent打包成二进制；将二进制文件以及启动所需的必要文件`config_agent.yml`、`config_server.yml`等文件进行打包，生成tar.gz压缩包文件。