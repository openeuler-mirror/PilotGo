# PilotGo

#### Introduction

PilotGo is an operations and maintenance (O\&M) management platform natively incubated by the openEuler community. It adopts a plugin-based architecture with modular, lightweight, and independently evolving functional components, while ensuring the stability of core features. Plugins enhance the platform's capabilities and break down barriers between different O\&M components, achieving global state awareness and automated workflows.

#### Feature Overview

The core modules of PilotGo include:

* User Management: Supports group management based on organizational structure and importing existing platform accounts for easy migration.

* Access Control: Role-Based Access Control (RBAC) for flexible and reliable permission management.

* Host Management: Visualized host state, direct management of software packages, services, and kernel parameters. Simple and user-friendly.

* Batch Execution: Enables concurrent execution of O\&M tasks, ensuring stability and efficiency.

* Audit Logging: Tracks user and plugin operations, facilitating issue tracing and security auditing.

* Alert Management: Real-time awareness of platform abnormalities.

* Plugin System: Extends platform functionalities. Plugins can interact with each other, boosting automation and reducing manual effort.

![Alt text](./docs/other/images/functional%20modules.png)


Integrated Plugins in the Current OS Release

* **Prometheus Plugin**: Hosts the Prometheus monitoring component, automates the deployment and configuration of `node-exporter` for data collection, and integrates with the platform’s alerting system.
  ![Alt text](./docs/other/images/prometheus%20plugin.png)

* **Grafana Plugin**: Integrates Grafana visualization, offering a beautiful and user-friendly metrics dashboard.
  ![Alt text](./docs/other/images/grafana%20plugin.png)

#### Application Scenarios

PilotGo is ideal for managing large-scale server clusters. It provides basic management and monitoring functionalities for mass server clusters. By integrating business-specific plugins, PilotGo also enables unified platform management for business clusters such as MySQL databases, Redis caches, and Nginx gateways.

#### Installation and Startup Guide

PilotGo supports both **standalone** and **clustered** deployment. Be sure to **disable the firewall** before installation.

1. **Install MySQL and Redis**, and configure their passwords.

2. **Install PilotGo Server** and modify the configuration file:

   ```bash
   dnf install -y PilotGo-server
   vim /opt/PilotGo/server/config_server.yaml
   ```
   
   **http_server**: addr is the address where PilotGo-server is installed;

   **socket_server**: addr is the address where PilotGo-server is installed;

   **mysql**: host_name is the address where MySQL is installed; user_name is the database login user; password is the database access password;

   **redis**: redis_conn is the address where the Redis service is installed; redis_pwd is the Redis password;

   **storage**: path is the storage path for the file service.

   For the first local startup, please fill in the MySQL password, Redis password, and file service storage path.

   **Start the service**:

     ```bash
     systemctl start PilotGo-server
     ```
   **Stop the service**:

     ```bash
     systemctl stop PilotGo-server
     ```
   **Check service status**:

     ```bash
     systemctl status PilotGo-server
     ```

3. **Install PilotGo Agent**:

   ```bash
   dnf install -y PilotGo-agent
   vim /opt/PilotGo/agent/config_agent.yaml
   ```

   `server.addr`: Address of the installed PilotGo-server.

   **Start the agent**:

     ```bash
     systemctl start PilotGo-agent
     ```

   **Stop the agent**:

     ```bash
     systemctl stop PilotGo-agent
     ```

   **Check agent status**:

     ```bash
     systemctl status PilotGo-agent
     ```

4. **Plugin Installation**:

   * [PilotGo-plugin-grafana Installation](https://gitee.com/src-openeuler/PilotGo-plugin-grafana)
   * [PilotGo-plugin-prometheus Installation](https://gitee.com/src-openeuler/PilotGo-plugin-prometheus)
   * [PilotGo-plugin-a-tune Installation](https://gitee.com/openeuler/PilotGo-plugin-a-tune)
   * [PilotGo-plugin-topology Installation](https://gitee.com/openeuler/PilotGo-plugin-topology)


#### Additional Links

1. [PilotGo User Manual](https://gitee.com/openeuler/docs/tree/master/docs/zh/docs/PilotGo/使用手册.md)
2. [PilotGo Software Repository](https://gitee.com/src-openeuler/PilotGo)
3. PilotGo Community Developer WeChat Group
![Alt text](./docs/other/images/PilotGo社区开发群.jpg)

#### Contribution

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request


#### Gitee Feature

1.  You can use Readme\_XXX.md to support different languages, such as Readme\_en.md, Readme\_zh.md
2.  Gitee blog [blog.gitee.com](https://blog.gitee.com)
3.  Explore open source project [https://gitee.com/explore](https://gitee.com/explore)
4.  The most valuable open source project [GVP](https://gitee.com/gvp)
5.  The manual of Gitee [https://gitee.com/help](https://gitee.com/help)
6.  The most popular members  [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)


