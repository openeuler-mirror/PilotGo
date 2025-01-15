/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Jun 9 16:03:53 2023 +0800
 */
package dbmanager

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/configfile"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/configmanage"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/cron"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/depart"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/role"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/script"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/user"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/redismanager"
)

func RedisdbInit(conf *options.RedisDBInfo, stopCh <-chan struct{}) error {
	err := redismanager.RedisInit(
		conf.RedisConn,
		conf.RedisPwd,
		conf.DefaultDB,
		conf.DialTimeout,
		conf.EnableRedis,
		stopCh,
		conf.UseTLS)
	if err != nil {
		return err
	}
	return nil
}

func MysqldbInit(conf *options.MysqlDBInfo) error {
	_, err := mysqlmanager.MysqlInit(
		conf.HostName,
		conf.UserName,
		conf.Password,
		conf.DataBase,
		conf.Port)
	if err != nil {
		return err
	}

	mysqlmanager.MySQL().AutoMigrate(&cron.CrontabList{})
	mysqlmanager.MySQL().AutoMigrate(&machine.MachineNode{})
	mysqlmanager.MySQL().AutoMigrate(&batch.Batch{})
	mysqlmanager.MySQL().AutoMigrate(&batch.BatchMachines{})
	mysqlmanager.MySQL().AutoMigrate(&auditlog.AuditLog{})
	mysqlmanager.MySQL().AutoMigrate(&configmanage.ConfigFiles{})
	mysqlmanager.MySQL().AutoMigrate(&configmanage.HistoryConfigFiles{})
	mysqlmanager.MySQL().AutoMigrate(&script.Script{}, &script.HistoryVersion{}, &script.DangerousCommands{})
	mysqlmanager.MySQL().AutoMigrate(&configfile.ConfigFile{})
	mysqlmanager.MySQL().AutoMigrate(&user.User{})
	mysqlmanager.MySQL().AutoMigrate(&role.Role{})
	mysqlmanager.MySQL().AutoMigrate(&role.UserRole{})

	// 创建超级管理员账户
	err = user.CreateAdministratorUser()
	if err != nil {
		return err
	}

	// 创建公司组织
	mysqlmanager.MySQL().AutoMigrate(&depart.DepartNode{})

	// 创建高危命令黑名单
	if err := script.CreateDangerousCommands(); err != nil {
		return err
	}

	return depart.CreateOrganization()
}
