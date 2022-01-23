package cmd

/**
 * @Author: zhang han
 * @Date: 2021/10/29 14:08
 * @Description: 启动程序、初始化、加载配置
 */

import (
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/router"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/net"
)

var (
	cfgFile string
	logger  = &logrus.Logger{}
	rootCmd = &cobra.Command{}
)

var sessionManage net.SessionManage
var sqlManager *mysqlmanager.MysqlManager

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
	sqlManager, err = mysqlmanager.Init(
		conf.Dbinfo.HostName,
		conf.Dbinfo.UserName,
		conf.Dbinfo.Password,
		conf.Dbinfo.DataBase,
		conf.Dbinfo.Port)
	if err != nil {
		return err
	}

	// sessionManage.Init(conf.MaxAge, conf.SessionCount)
	// go func() {
	// 	for true {
	// 		time.Sleep(time.Second * 10)
	// 		//每10秒读取一次数据库，并更改数据库状态
	// 		mi, err := mysqlmanager.GetMachInfo(sqlManager)
	// 		if err != nil {
	// 			continue
	// 		}

	// 		for _, m := range mi {
	// 			status := m.CheckStatus()
	// 			if m.SystemStatus != status {
	// 				m1 := mysqlmanager.MachInfo{
	// 					Id:           m.Id,
	// 					SystemStatus: status,
	// 				}
	// 				sqlManager.Update(&m1)
	// 			}
	// 		}
	// 	}
	// }()

	// pi, err := mysqlmanager.GetPluginInfo(sqlManager)
	// if err != nil {
	// 	return err
	// }

	// for _, value := range pi {
	// 	plugin.GetManager().Regist(&plugin.Plugin{
	// 		Name:        value.Name,
	// 		Version:     value.Version,
	// 		Description: value.Description,
	// 		Url:         value.Url,
	// 		Port:        value.Port,
	// 		Protocol:    value.Protocol,
	// 	})
	// }

	mysqlmanager.DB.AutoMigrate(&model.User{})
	mysqlmanager.DB.AutoMigrate(&model.DepartNode{})
	mysqlmanager.DB.AutoMigrate(&model.MachineNode{})
	mysqlmanager.DB.AutoMigrate(&model.Batch{})
	defer mysqlmanager.DB.Close()

	r := router.SetupRouter()
	server_url := ":" + strconv.Itoa(conf.S.ServerPort)
	r.Run(server_url)

	return http.ListenAndServe(server_url, nil) // listen and serve
}
