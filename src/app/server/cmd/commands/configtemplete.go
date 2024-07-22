package commands

import (
	"os"

	"gitee.com/openeuler/PilotGo/app/server/cmd/options"
	"gitee.com/openeuler/PilotGo/utils"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewTempleteCommand() *cobra.Command {

	command := &cobra.Command{
		Use:   "templete",
		Short: "create the templete of pilotgo server start",
		RunE: func(cmd *cobra.Command, args []string) error {
			return templeteRun()
		},
		Example: `
		# create the templete for api-server start
		pilotgo-server templete
		`,
	}
	return command
}
func templeteRun() error {
	config := options.ServerConfig{
		HttpServer:   options.NewHttpServerOptions(),
		SocketServer: options.NewSocketServerOptions(),
		JWT:          options.NewJWTConfigOptions(),
		Logopts:      options.NewLogOptsOptions(),
		MysqlDBinfo:  options.NewMysqlDBInfoOpts(),
		RedisDBinfo:  options.NewRedisDBInfoOpts(),
		Storage:      options.NewStorageOptions(),
	}
	operator := utils.NewYamlOpeartor(config,
		utils.WithCommentsTagFlag(utils.PilotGoHeadComment),
		utils.WithDefaultTagName(utils.DefaultTagName))
	yamlContent, err := operator.Marshal()
	if err != nil {
		klog.Infof("Marshal error: %v", err)
		return err
	}
	err = os.WriteFile("./config_server.yaml.templete", yamlContent, os.ModePerm)
	if err != nil {
		klog.Infof("create templete error: %v", err)
		return err
	}
	klog.Info("create config_server.yaml.templete file success,please search it in current directory.")
	return nil
}
