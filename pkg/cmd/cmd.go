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
 * Date: 2022-02-23 17:11:01
 * LastEditTime: 2022-03-16 13:03:03
 * Description: 启动程序、初始化、加载配置
 ******************************************************************************/
package cmd

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/router"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/net"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{}
)

var sessionManage net.SessionManage
var SqlManager *mysqlmanager.MysqlManager
var menus string = "cluster,batch,usermanager,rolemanager"

func initConfig() {
	config.MustInit(os.Stdout, cfgFile) // 配置初始化
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "pkg/config", "pkg/config/dev.yaml", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("debug", true, "开启debug")
	viper.SetDefault("gin.mode", rootCmd.PersistentFlags().Lookup("debug"))
}

func Start(conf *config.Configure) (err error) {
	SqlManager, err = mysqlmanager.Init(
		conf.Dbinfo.HostName,
		conf.Dbinfo.UserName,
		conf.Dbinfo.Password,
		conf.Dbinfo.DataBase,
		conf.Dbinfo.Port)
	if err != nil {
		return err
	}

	sessionManage.Init(conf.MaxAge, conf.SessionCount)
	go func() {
		for {
			controller.AddAgents()
			time.Sleep(time.Second * 30)
		}
	}()

	mysqlmanager.DB.Delete(&model.MachineNode{})

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
	defer mysqlmanager.DB.Close()

	r := router.SetupRouter()
	server_url := ":" + strconv.Itoa(conf.S.ServerPort)
	r.Run(server_url)

	return http.ListenAndServe(server_url, nil) // listen and serve
}
