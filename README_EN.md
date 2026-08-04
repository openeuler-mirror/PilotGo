# PilotGo

## Overview

PilotGo is an operations and maintenance (O\&M) management platform natively incubated by the openEuler community. It adopts a plugin-based architecture with modular, lightweight, and independently evolving functional components, while ensuring the stability of core features. Plugins enhance the platform's capabilities and break down barriers between different O\&M components, achieving global state awareness and automated workflows.

## Features

The core modules of PilotGo include:

* User Management: Supports group management based on organizational structure and importing existing platform accounts for easy migration.

* Access Control: Role-Based Access Control (RBAC) for flexible and reliable permission management.

* Host Management: Visualized host state, direct management of software packages, services, and kernel parameters. Simple and user-friendly.

* Batch Execution: Enables concurrent execution of O\&M tasks, ensuring stability and efficiency.

* Audit Logging: Tracks user and plugin operations, facilitating issue tracing and security auditing.

* Alert Management: Real-time awareness of platform abnormalities.

* Plugin System: Extends platform functionalities. Plugins can interact with each other, boosting automation and reducing manual effort.

![Alt text](./docs/other/images/functional%20modules_en.png)

Integrated plugins in the current OS release:

* Prometheus: Hosts the Prometheus monitoring component, automates the deployment and configuration of `node-exporter` for data collection, and integrates with the platform's alerting system.
  ![Alt text](./docs/other/images/prometheus%20plugin.png)

* Grafana: Integrates Grafana visualization, offering a user-friendly metrics dashboard.
  ![Alt text](./docs/other/images/grafana%20plugin.png)

## Application Scenarios

PilotGo is ideal for managing large-scale server clusters. It provides basic management and monitoring functionalities for mass server clusters. By integrating business-specific plugins, PilotGo also enables unified platform management for business clusters such as MySQL databases, Redis caches, and Nginx gateways.

## Installation and Startup Guide

PilotGo supports both **standalone** and **clustered** deployment. Be sure to disable the firewall before installation.

1. Install MySQL and Redis, and configure their passwords.

2. Install PilotGo-server and modify the configuration file:

   ```bash
   dnf install -y PilotGo-server
   vim /opt/PilotGo/server/config_server.yaml
   ```
   
   `http_server`: *addr* is the address where PilotGo-server is installed.

   `socket_server`: *addr* is the address where PilotGo-server is installed.

   `mysql`: *host_name* is the address where MySQL is installed; *user_name* is the database login user; *password* is the database access password.

   `redis`: *redis_conn* is the address where the Redis service is installed; *redis_pwd* is the Redis password.

   `storage`: *path* is the storage path for the file service.

   For the first local startup, please fill in the MySQL password, Redis password, and file service storage path.

   Start the service:

     ```bash
     systemctl start PilotGo-server
     ```

   Stop the service:

     ```bash
     systemctl stop PilotGo-server
     ```
   
   Check service status:

     ```bash
     systemctl status PilotGo-server
     ```

3. Install PilotGo-agent:

   ```bash
   dnf install -y PilotGo-agent
   vim /opt/PilotGo/agent/config_agent.yaml
   ```

   `server`: *addr* is the address where PilotGo-server is installed.

   Start the service:

     ```bash
     systemctl start PilotGo-agent
     ```

   Stop the service:

     ```bash
     systemctl stop PilotGo-agent
     ```

   Check service status:

     ```bash
     systemctl status PilotGo-agent
     ```

4. Install plugins.

   * [PilotGo-plugin-grafana](https://atomgit.com/src-openeuler/PilotGo-plugin-grafana)
   * [PilotGo-plugin-prometheus](https://atomgit.com/src-openeuler/PilotGo-plugin-prometheus)
   * [PilotGo-plugin-a-tune](https://atomgit.com/openeuler/PilotGo-plugin-a-tune)
   * [PilotGo-plugin-topology](https://atomgit.com/openeuler/PilotGo-plugin-topology)

## Additional Links

[PilotGo Software Repository](https://atomgit.com/src-openeuler/PilotGo)

## Contribution

1. Fork this repository.
2. Create a Feat_*xxx* branch.
3. Commit code.
4. Create a pull request (PR).

## Notes

Use the file naming pattern `README_xx.md` to indicate a supported language (for example, `README_EN.md`).
