/*
PilotGo项目数据库中权限按钮菜单，运行此sql文件可鉴权登陆
Source Host           : 127.0.1.1:3306
Source Database       : PilotGo
Target Server Type    : MYSQL
*/

SET FOREIGN_KEY_CHECKS=0;

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
