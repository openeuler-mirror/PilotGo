# PilotGo

#### 介绍

PilotGo 是 openEuler 社区原生孵化的运维管理平台，采用插件式架构设计，功能模块轻量化组合、独立迭代演进，同时保证核心功能稳定；同时使用插件来增强平台功能、并打通不同运维组件之间的壁垒，实现了全局的状态感知及自动化流程。

#### 功能描述

PilotGo 核心功能模块包括：

* 用户管理：支持按照组织结构分组管理，支持导入已有平台账号，迁移方便；

* 权限管理：支持基于RBAC的权限管理，灵活可靠；
  
* 主机管理：状态前端可视化、直接执行软件包管理、服务管理、内核参数调优、简单易操作；
  
* 批次管理：支持运维操作并发执行，稳定高效；
 
* 日志审计：跟踪记录用户及插件的变更操作，方便问题回溯及安全审计；

* 告警管理：平台异常实时感知；

* 插件功能：支持扩展平台功能，插件联动，自动化能力倍增，减少人工干预。

![Alt text](./docs/images/functional%20modules.png)


当前OS发布版本还集成了以下插件：

* Prometheus：托管Prometheus监控组件，自动化下发及配置node-exporter监控数据采集，对接平台告警功能；![Alt text](./docs/images/prometheus%20plugin.png)

* Grafana：集成Grafana可视化平台，提供美观易用的指标监控面板功能。
![Alt text](./docs/images/grafana%20plugin.png)

#### 应用场景

PiotGo可用于典型的服务器集群管理场景，支持大批量的服务器集群基本管理及监控；通过集成对应的业务功能插件，还可实现业务集群的统一平台管理，例如Mysql数据库集群、redis数据缓存集群、nginx网关集群等。

#### 安装、启动教程

PilotGo可以单机部署也可以采用集群式部署。安装之前先关闭防火墙。
1.  安装mysql、redis，并设置密码；
2.  安装PilotGo-server，并修改配置文件:
   >dnf install -y PilotGo-server

   >vim /opt/PilotGo/server/config_server.yaml

   http_server：addr为安装PilotGo-server地址；

   socket_server：addr为安装PilotGo-server地址；

   mysql：host_name为安装mysql地址；user_name为DB的登录用户；password为DB访问密码；

   redis：redis_conn为安装redis服务地址；redis_pwd为redis密码；

   启动服务
   >systemctl start PilotGo-server

   停止服务
   >ystemctl stop PilotGo-server

   服务状态
   >systemctl status PilotGo-server
3.  安装PilotGo-agent：
   >dnf install -y PilotGo-agent
   
   >vim /opt/PilotGo/agent/config_agent.yaml
   
   server：addr为安装PilotGo-server地址；
   
   启动服务
   >systemctl start PilotGo-agent

   停止服务
   >systemctl stop PilotGo-agent

   服务状态
   >systemctl status PilotGo-agent
4.  插件安装：
   [PilotGo-plugin-grafana插件安装](https://gitee.com/src-openeuler/PilotGo-plugin-grafana)  
   [PilotGo-plugin-prometheus插件安装](https://gitee.com/src-openeuler/PilotGo-plugin-prometheus)

#### 补充链接

1.  [PilotGo使用手册](https://gitee.com/openeuler/docs/tree/master/docs/zh/docs/PilotGo/使用手册.md)
2.  PilotGo[软件包仓](https://gitee.com/src-openeuler/PilotGo)


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
