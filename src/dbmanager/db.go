package dbmanager

import (
	sconfig "gitee.com/openeuler/PilotGo/app/server/config"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/app/server/service/configfile"
	"gitee.com/openeuler/PilotGo/app/server/service/configmanage"
	"gitee.com/openeuler/PilotGo/app/server/service/cron"
	"gitee.com/openeuler/PilotGo/app/server/service/depart"
	"gitee.com/openeuler/PilotGo/app/server/service/machine"
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/app/server/service/role"
	"gitee.com/openeuler/PilotGo/app/server/service/script"
	"gitee.com/openeuler/PilotGo/app/server/service/user"
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/dbmanager/redismanager"
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

	mysqlmanager.MySQL().AutoMigrate(&cron.CrontabList{})
	mysqlmanager.MySQL().AutoMigrate(&machine.MachineNode{})
	mysqlmanager.MySQL().AutoMigrate(&batch.Batch{})
	mysqlmanager.MySQL().AutoMigrate(&auditlog.AuditLog{})
	mysqlmanager.MySQL().AutoMigrate(&configmanage.ConfigFiles{})
	mysqlmanager.MySQL().AutoMigrate(&configmanage.HistoryConfigFiles{})
	mysqlmanager.MySQL().AutoMigrate(&script.Script{})
	mysqlmanager.MySQL().AutoMigrate(&configfile.ConfigFile{})
	mysqlmanager.MySQL().AutoMigrate(&plugin.Plugin{})
	mysqlmanager.MySQL().AutoMigrate(&user.User{})
	mysqlmanager.MySQL().AutoMigrate(&role.UserRole{})

	// 创建超级管理员账户
	err = user.CreateAdministratorUser()
	if err != nil {
		return err
	}

	// 创建公司组织
	mysqlmanager.MySQL().AutoMigrate(&depart.DepartNode{})

	return depart.CreateOrganization()
}
