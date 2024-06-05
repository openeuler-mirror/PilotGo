# 开发文档

## 开发环境
golang: v1.17  
nodejs: v14  
OS：openEuler、kylinOS

## 代码结构

src/: golang源码，包含server、agent及公共代码等  
src/app/server: server端代码，包含server入口及server特有模块包  
src/app/agent: agent端代码，包含agent入口及agent特有模块包  
frontend/: web端源码