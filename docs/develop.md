# 开发文档

## 开发环境
golang: v1.15
nodejs: v14
OS：openEuler、kylinOS

## 代码结构

pkg/: golang源码，包含server、agent及公共代码等
pkg/app/server: server端代码，包含server入口及server特有模块包
pkg/app/server: agent端代码，包含agent入口及agent特有模块包
src/: web端源码