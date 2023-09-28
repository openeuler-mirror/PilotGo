package dbmanager

import (
	sconfig "gitee.com/PilotGo/PilotGo/app/server/config"
	"gitee.com/PilotGo/PilotGo/app/server/dao"
	"gitee.com/PilotGo/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/PilotGo/PilotGo/dbmanager/redismanager"
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

	mysqlmanager.MySQL().AutoMigrate(&dao.CrontabList{})
	mysqlmanager.MySQL().AutoMigrate(&dao.MachineNode{})
	mysqlmanager.MySQL().AutoMigrate(&dao.Batch{})
	mysqlmanager.MySQL().AutoMigrate(&dao.AgentLogParent{})
	mysqlmanager.MySQL().AutoMigrate(&dao.AgentLog{})
	mysqlmanager.MySQL().AutoMigrate(&dao.AuditLog{})
	mysqlmanager.MySQL().AutoMigrate(&dao.ConfigFiles{})
	mysqlmanager.MySQL().AutoMigrate(&dao.HistoryConfigFiles{})
	mysqlmanager.MySQL().AutoMigrate(&dao.Script{})
	mysqlmanager.MySQL().AutoMigrate(&dao.ConfigFile{})
	mysqlmanager.MySQL().AutoMigrate(&dao.PluginModel{})
	mysqlmanager.MySQL().AutoMigrate(&dao.User{})
	mysqlmanager.MySQL().AutoMigrate(&dao.UserRole{})

	// 创建超级管理员账户
	err = dao.CreateAdministratorUser()
	if err != nil {
		return err
	}

	// 创建公司组织
	mysqlmanager.MySQL().AutoMigrate(&dao.DepartNode{})

	return dao.CreateOrganization()
}
