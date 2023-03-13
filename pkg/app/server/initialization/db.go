/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-07-05 13:03:16
 * LastEditTime: 2022-07-05 14:10:23
 * Description: db and redis init
 ******************************************************************************/
package initialization

import (
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/dbmanager/redismanager"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

func RedisdbInit(conf *sconfig.RedisDBInfo) error {
	err := redismanager.RedisInit(
		conf.RedisConn,
		conf.RedisPwd,
		conf.DefaultDB,
		conf.DialTimeout,
		conf.EnableRedis)
	if err != nil {
		return err
	}
	return nil
}

func MysqldbInit(conf *sconfig.MysqlDBInfo) error {
	_, err := mysqlmanager.MysqlInit(
		conf.HostName,
		conf.UserName,
		conf.Password,
		conf.DataBase,
		conf.Port)
	if err != nil {
		return err
	}

	global.PILOTGO_DB.AutoMigrate(&dao.CrontabList{})
	global.PILOTGO_DB.AutoMigrate(&model.MachineNode{})
	global.PILOTGO_DB.AutoMigrate(&dao.RoleButton{})
	global.PILOTGO_DB.AutoMigrate(&model.Batch{})
	global.PILOTGO_DB.AutoMigrate(&dao.AgentLogParent{})
	global.PILOTGO_DB.AutoMigrate(&dao.AgentLog{})
	global.PILOTGO_DB.AutoMigrate(&dao.Files{})
	global.PILOTGO_DB.AutoMigrate(&dao.HistoryFiles{})
	global.PILOTGO_DB.AutoMigrate(&dao.Script{})
	global.PILOTGO_DB.AutoMigrate(&dao.ConfigFile{})
	global.PILOTGO_DB.AutoMigrate(&dao.PluginModel{})

	// 创建超级管理员账户
	global.PILOTGO_DB.AutoMigrate(&dao.User{})
	global.PILOTGO_DB.AutoMigrate(&dao.UserRole{})
	err = dao.CreateAdministratorUser()
	if err != nil {
		return err
	}

	// 创建公司组织
	global.PILOTGO_DB.AutoMigrate(&model.DepartNode{})

	return dao.CreateOrganization()
}
