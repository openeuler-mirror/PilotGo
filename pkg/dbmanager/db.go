package dbmanager

import (
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
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
	global.PILOTGO_DB.AutoMigrate(&dao.MachineNode{})
	global.PILOTGO_DB.AutoMigrate(&dao.RoleButton{})
	global.PILOTGO_DB.AutoMigrate(&dao.Batch{})
	global.PILOTGO_DB.AutoMigrate(&dao.AgentLogParent{})
	global.PILOTGO_DB.AutoMigrate(&dao.AgentLog{})
	global.PILOTGO_DB.AutoMigrate(&dao.AuditLog{})
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
	global.PILOTGO_DB.AutoMigrate(&dao.DepartNode{})

	return dao.CreateOrganization()
}
