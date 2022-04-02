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
 * LastEditTime: 2022-04-02 11:28:13
 * Description: server main
 ******************************************************************************/
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/router"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/net"
)

func main() {

	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	if err := logger.Init(conf); err != nil {
		fmt.Println("logger init failed, please check the config file")
		os.Exit(-1)
	}
	logger.Info("Thanks to choose PilotGo!")

	// db初始化
	if err := dbInit(conf); err != nil {
		logger.Error("logger init failed, please check the config file")
		os.Exit(-1)
	}

	// 监控初始化
	if err := monitorInit(); err != nil {
		logger.Error("monitor init failed")
		os.Exit(-1)
	}

	// 启动agent socket server
	if err := sockerServerInit(conf); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		os.Exit(-1)
	}

	//此处启动前端及REST http server
	if err := httpServerInit(conf); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		os.Exit(-1)
	}

	logger.Info("start to serve.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			// TODO: DO EXIT

			mysqlmanager.DB.Close()

			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
}

func sessionManagerInit(conf *config.Configure) error {
	var sessionManage net.SessionManage
	sessionManage.Init(conf.MaxAge, conf.SessionCount)
	go func() {
		for {
			controller.AddAgents()
			time.Sleep(time.Second * 30)
		}
	}()
	return nil
}

func dbInit(conf *config.Configure) error {
	var menus string = "cluster,batch,usermanager,rolemanager"

	_, err := mysqlmanager.Init(
		conf.Dbinfo.HostName,
		conf.Dbinfo.UserName,
		conf.Dbinfo.Password,
		conf.Dbinfo.DataBase,
		conf.Dbinfo.Port)
	if err != nil {
		return err
	}

	// mysqlmanager.DB.Delete(&model.MachineNode{})
	mysqlmanager.DB.AutoMigrate(&model.User{})
	mysqlmanager.DB.AutoMigrate(&model.UserRole{})

	var user model.User
	var role model.UserRole
	pid := 0
	mysqlmanager.DB.Where("depart_first=?", pid).Find(&user)
	if user.ID == 0 {
		user = model.User{
			CreatedAt:    time.Time{},
			DepartFirst:  0,
			DepartSecond: 1,
			DepartName:   "超级用户",
			Username:     "admin",
			Password:     "1234",
			Email:        "admin@123.com",
			UserType:     0,
			RoleID:       "1",
		}
		mysqlmanager.DB.Create(&user)
		role = model.UserRole{
			Role:  "超级管理员",
			Type:  0,
			Menus: menus,
		}
		mysqlmanager.DB.Create(&role)
	}

	mysqlmanager.DB.AutoMigrate(&model.DepartNode{})
	var Depart model.DepartNode
	mysqlmanager.DB.Where("p_id=?", pid).Find(&Depart)
	if Depart.ID == 0 {
		Depart = model.DepartNode{
			PID:          0,
			ParentDepart: "",
			Depart:       "组织名",
			NodeLocate:   0,
		}
		mysqlmanager.DB.Save(&Depart)
	}
	mysqlmanager.DB.AutoMigrate(&model.MachineNode{})
	mysqlmanager.DB.AutoMigrate(&model.RoleButton{})
	mysqlmanager.DB.AutoMigrate(&model.Batch{})
	mysqlmanager.DB.AutoMigrate(&model.AgentLogParent{})
	mysqlmanager.DB.AutoMigrate(&model.AgentLog{})
	// defer mysqlmanager.DB.Close()

	return nil
}

func sockerServerInit(conf *config.Configure) error {
	server := &network.SocketServer{
		// MessageProcesser: protocol.NewMessageProcesser(),
		OnAccept: agentmanager.AddandRunAgent,
		OnStop:   agentmanager.StopAgentManager,
	}
	url := conf.S.ServerIP + ":" + strconv.Itoa(conf.SocketPort)
	go func() {
		server.Run(url)
	}()
	return nil
}

func httpServerInit(conf *config.Configure) error {
	if err := sessionManagerInit(conf); err != nil {
		return err
	}

	go func() {
		r := router.SetupRouter()
		server_url := ":" + strconv.Itoa(conf.S.ServerPort)
		r.Run(server_url)

		err := http.ListenAndServe(server_url, nil) // listen and serve
		if err != nil {
			logger.Error("failed to start http server, error:%v", err)
		}
	}()

	return nil
}

func monitorInit() error {
	go func() {
		logger.Info("start monitor")
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
			err := controller.WritePrometheusYml(a)
			if err != nil {
				logger.Error("写入promethues配置文件失败")
			}
			conf, err := config.Load()
			if err != nil {
				logger.Error("%s", "failed to load configure, exit.."+err.Error())

			}
			err = controller.PrometheusConfigReload(conf.S.ServerIP)
			if err != nil {
				logger.Error("%s", err.Error())
			}
			time.Sleep(100 * time.Second)
		}

	}()

	return nil
}
