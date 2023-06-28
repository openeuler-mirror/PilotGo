/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-06-27 09:59:30
 * LastEditTime: 2023-06-28 16:10:54
 * Description: provide global const of pilotgo
 ******************************************************************************/
package global

import (
	"github.com/casbin/casbin/v2"
)

// 实例变量
var (
	// PILOTGO_DB    *gorm.DB
	// PILOTGO_REDIS *redis.Client
	PILOTGO_E *casbin.Enforcer
)

// 权限菜单、按钮
const (
	PILOTGO_MENUS    = "overview,cluster,batch,usermanager,rolemanager,config,log,plugin"
	PILOTGO_BUTTONID = "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15"
)

// 用户、角色
const (
	// 超级管理员
	AdminUserType = 0
	// 部门管理员
	DepartManagerType = 1
	// 普通用户
	OrdinaryUserType = 2
	// 其他用户，如实习生
	OtherUserType = 3
	//普通用户角色id
	OrdinaryUserRoleId = 3
	// 默认用户密码
	DefaultUserPassword = "123456"
)

// 单机操作成功状态:是否成功，机器数量，比率
const (
	ActionOK    = "1,1,1.00"
	ActionFalse = "0,1,0.00"
)

// 配置文件类型
const (
	ConfigRepo = "repo配置"
)

// 配置文件源路径
const (
	// repo路径
	RepoPath = "/etc/yum.repos.d"
	// 网络配置
	NetWorkPath = "/etc/sysconfig/network-scripts"
)

// 机器运行状态
const (
	// 机器运行
	Normal = 1
	// 脱机
	OffLine = 2
	// 空闲
	Free = 3
	// 新注册机器添加到部门根节点
	UncateloguedDepartId = 1
	// 是否为部门根节点
	Departroot   = 0
	DepartUnroot = 1
)
