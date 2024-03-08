/*
Navicat MySQL Data Transfer

Source Server         : 服务器
Source Server Version : 80028
Source Host           : 127.0.0.1:3306
Source Database       : PilotGo

Target Server Type    : MYSQL
Target Server Version : 80028
File Encoding         : 65001

Date: 2022-07-05 10:15:12
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for agent_log
-- ----------------------------
DROP TABLE IF EXISTS `agent_log`;
CREATE TABLE `agent_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `log_parent_id` int DEFAULT NULL,
  `ip` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status_code` int DEFAULT NULL,
  `operation_object` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `action` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `message` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_agent_log_log_parent_id` (`log_parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of agent_log
-- ----------------------------

-- ----------------------------
-- Table structure for agent_log_parent
-- ----------------------------
DROP TABLE IF EXISTS `agent_log_parent`;
CREATE TABLE `agent_log_parent` (
  `id` int NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `user_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `depart_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of agent_log_parent
-- ----------------------------

-- ----------------------------
-- Table structure for batch
-- ----------------------------
DROP TABLE IF EXISTS `batch`;
CREATE TABLE `batch` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `description` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `manager` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `machinelist` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `depart` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `depart_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_batch_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of batch
-- ----------------------------

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `p_type` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `v0` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `v1` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `v2` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `v3` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `v4` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `v5` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------

-- ----------------------------
-- Table structure for crontab_list
-- ----------------------------
DROP TABLE IF EXISTS `crontab_list`;
CREATE TABLE `crontab_list` (
  `id` int NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `machine_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `task_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `cron_spec` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `command` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of crontab_list
-- ----------------------------

-- ----------------------------
-- Table structure for depart_node
-- ----------------------------
DROP TABLE IF EXISTS `depart_node`;
CREATE TABLE `depart_node` (
  `id` int NOT NULL AUTO_INCREMENT,
  `p_id` int NOT NULL,
  `parent_depart` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `depart` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `node_locate` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of depart_node
-- ----------------------------
INSERT INTO `depart_node` VALUES ('1', '0', '', '组织名', '0');

-- ----------------------------
-- Table structure for files
-- ----------------------------
DROP TABLE IF EXISTS `files`;
CREATE TABLE `files` (
  `id` int NOT NULL AUTO_INCREMENT,
  `file_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `file_path` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_update` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_dept` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `controlled_batch` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `take_effect` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `file` text COLLATE utf8_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of files
-- ----------------------------

-- ----------------------------
-- Table structure for history_files
-- ----------------------------
DROP TABLE IF EXISTS `history_files`;
CREATE TABLE `history_files` (
  `id` int NOT NULL AUTO_INCREMENT,
  `file_id` int DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `user_update` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `user_dept` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `file_name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `file` text COLLATE utf8_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of history_files
-- ----------------------------

-- ----------------------------
-- Table structure for machine_node
-- ----------------------------
DROP TABLE IF EXISTS `machine_node`;
CREATE TABLE `machine_node` (
  `id` int NOT NULL AUTO_INCREMENT,
  `depart_id` int NOT NULL,
  `ip` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `machine_uuid` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `cpu` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `state` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `systeminfo` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of machine_node
-- ----------------------------

-- ----------------------------
-- Table structure for role_button
-- ----------------------------
DROP TABLE IF EXISTS `role_button`;
CREATE TABLE `role_button` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `button` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of role_button
-- ----------------------------
INSERT INTO `role_button` VALUES ('1', 'rpm_install');
INSERT INTO `role_button` VALUES ('2', 'rpm_uninstall');
INSERT INTO `role_button` VALUES ('3', 'batch_update');
INSERT INTO `role_button` VALUES ('4', 'batch_delete');
INSERT INTO `role_button` VALUES ('5', 'user_add');
INSERT INTO `role_button` VALUES ('6', 'user_import');
INSERT INTO `role_button` VALUES ('7', 'user_edit');
INSERT INTO `role_button` VALUES ('8', 'user_reset');
INSERT INTO `role_button` VALUES ('9', 'user_del');
INSERT INTO `role_button` VALUES ('10', 'role_add');
INSERT INTO `role_button` VALUES ('11', 'role_update');
INSERT INTO `role_button` VALUES ('12', 'role_delete');
INSERT INTO `role_button` VALUES ('13', 'role_modify');
INSERT INTO `role_button` VALUES ('14', 'config_install');
INSERT INTO `role_button` VALUES ('15', 'dept_change');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `depart_first` int DEFAULT NULL,
  `depart_second` int DEFAULT NULL,
  `depart_name` varchar(25) COLLATE utf8_unicode_ci DEFAULT NULL,
  `username` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `password` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `phone` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `email` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  `user_type` int DEFAULT NULL,
  `role_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', '2022-07-05 10:13:46', '0', '1', '超级用户', 'admin', '123456', '', 'admin@123.com', '0', '1');

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `id` int NOT NULL AUTO_INCREMENT,
  `role` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` int DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `menus` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `button_id` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of user_role
-- ----------------------------
INSERT INTO `user_role` VALUES ('1', '超级用户', '0', '超级管理员', 'overview,cluster,batch,usermanager,rolemanager,config,log,plug-in,plugin-web', '1,2,3,4,5,6,7,8,9,10,11,12,13,14,15');
