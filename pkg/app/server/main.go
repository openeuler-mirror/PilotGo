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
 * Date: 2021-11-18 10:25:52
 * LastEditTime: 2022-04-18 15:16:10
 * Description: server main
 ******************************************************************************/
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	sconfig "openeluer.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/router"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/redismanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func main() {
	err := sconfig.Init()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	if err := logger.Init(&sconfig.Config().Logopts); err != nil {
		fmt.Println("logger init failed, please check the config file")
		os.Exit(-1)
	}
	logger.Info("Thanks to choose PilotGo!")

	// redis db初始化
	if err := redisdbInit(&sconfig.Config().RedisDBinfo); err != nil {
		logger.Error("redis db init failed, please check the config file")
		os.Exit(-1)
	}

	// mysql db初始化
	if err := mysqldbInit(&sconfig.Config().MysqlDBinfo); err != nil {
		logger.Error("mysql db init failed, please check the config file")
		os.Exit(-1)
	}

	// 监控初始化
	if err := monitorInit(&sconfig.Config().Monitor); err != nil {
		logger.Error("monitor init failed")
		os.Exit(-1)
	}

	// 启动agent socket server
	if err := socketServerInit(&sconfig.Config().SocketServer); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		os.Exit(-1)
	}

	//此处启动前端及REST http server
	if err := httpServerInit(&sconfig.Config().HttpServer); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		os.Exit(-1)
	}

	logger.Info("start to serve.")

	// 信号监听
	agentmanager.WARN_MSG = make(chan interface{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			// TODO: DO EXIT

			mysqlmanager.DB.Close()
			redismanager.Redis.Close()

			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
}

func sessionManagerInit(conf *sconfig.HttpServer) error {
	var sessionManage service.SessionManage
	sessionManage.Init(conf.SessionMaxAge, conf.SessionCount)
	return nil
}

func redisdbInit(conf *sconfig.RedisDBInfo) error {
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

func mysqldbInit(conf *sconfig.MysqlDBInfo) error {
	_, err := mysqlmanager.MysqlInit(
		conf.HostName,
		conf.UserName,
		conf.Password,
		conf.DataBase,
		conf.Port)
	if err != nil {
		return err
	}

	// 创建超级管理员账户
	mysqlmanager.DB.AutoMigrate(&model.User{})
	mysqlmanager.DB.AutoMigrate(&model.UserRole{})
	dao.CreateSuperAdministratorUser()

	// 创建公司组织
	mysqlmanager.DB.AutoMigrate(&model.DepartNode{})
	dao.CreateOrganization()

	mysqlmanager.DB.AutoMigrate(&model.CrontabList{})
	mysqlmanager.DB.AutoMigrate(&model.MachineNode{})
	mysqlmanager.DB.AutoMigrate(&model.RoleButton{})
	mysqlmanager.DB.AutoMigrate(&model.Batch{})
	mysqlmanager.DB.AutoMigrate(&model.AgentLogParent{})
	mysqlmanager.DB.AutoMigrate(&model.AgentLog{})
	mysqlmanager.DB.AutoMigrate(&model.Files{})
	mysqlmanager.DB.AutoMigrate(&model.HistoryFiles{})

	return nil
}

func socketServerInit(conf *sconfig.SocketServer) error {
	server := &network.SocketServer{
		// MessageProcesser: protocol.NewMessageProcesser(),
		OnAccept: agentmanager.AddandRunAgent,
		OnStop:   agentmanager.StopAgentManager,
	}

	go func() {
		server.Run(conf.Addr)
	}()
	return nil
}

func httpServerInit(conf *sconfig.HttpServer) error {
	if err := sessionManagerInit(conf); err != nil {
		return err
	}

	go func() {
		r := router.SetupRouter()
		r.Run(conf.Addr)

		err := http.ListenAndServe(conf.Addr, nil) // listen and serve
		if err != nil {
			logger.Error("failed to start http server, error:%v", err)
		}
	}()

	return nil
}

func monitorInit(conf *sconfig.Monitor) error {
	go func() {
		logger.Info("start monitor")
		err := controller.InitPromeYml()
		if err != nil {
			logger.Error("初始化promethues配置文件失败")
		}
		for {
			// TODO: 重构为事件触发机制
			a := make([]map[string]string, 0)
			var m []model.MachineNode
			mysqlmanager.DB.Find(&m)
			for _, value := range m {
				r := map[string]string{}
				r[value.MachineUUID] = value.IP
				a = append(a, r)
			}
			err := controller.WriteYml(a)
			if err != nil {
				logger.Error("%s", err.Error())
			}
			time.Sleep(100 * time.Second)
		}

	}()

	return nil
}
