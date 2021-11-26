package controller

/**
 * @Author: zhang han
 * @Date: 2021/11/15 10:13
 * @Description: 防火墙命令及配置
 */

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func Config(c *gin.Context) {
  ip := c.PostForm("ip")
  host_user := c.PostForm("host_user")
  host_password := c.PostForm("host_password")
  cli := common.NewSsh(ip, host_user, host_password, 22)
  tmp, err := cli.Run("firewall-cmd --list-all")
  if err != nil {
    response.Response(c, http.StatusUnprocessableEntity,
      422,
      nil,
      "获取防火墙配置失败")
    return
  }
  response.Success(c, gin.H{"tmp": tmp}, "获取防火墙配置成功")
}

func Stop(c *gin.Context) {
	ip := c.PostForm("ip")
	host_user := c.PostForm("host_user")
	host_password := c.PostForm("host_password")
	cli := common.NewSsh(ip, host_user, host_password, 22)
	tmp, err := cli.Run("service firewalld stop")
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"关闭防火墙失败")
		return
	}
	response.Success(c, gin.H{"tmp": tmp}, "关闭防火墙成功")
}

func Restart(c *gin.Context) {
	ip := c.PostForm("ip")
	host_user := c.PostForm("host_user")
	host_password := c.PostForm("host_password")
	cli := common.NewSsh(ip, host_user, host_password, 22)
	tmp, err := cli.Run("service firewalld restart")
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"重启防火墙失败")
		return
	}
	response.Success(c, gin.H{"tmp": tmp}, "重启防火墙成功")
}

func Reload(c *gin.Context) {
	ip := c.PostForm("ip")
	host_user := c.PostForm("host_user")
	host_password := c.PostForm("host_password")
	cli := common.NewSsh(ip, host_user, host_password, 22)
	tmp, err := cli.Run("firewall-cmd --reload")
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"更新防火墙失败")
		return
	}
	response.Success(c, gin.H{"tmp": tmp}, "更新防火墙成功")
}

func AddZonePort(c *gin.Context) { //zone = block dmz drop external home internal public trusted work
	ip := c.PostForm("ip")
	host_user := c.PostForm("host_user")
	host_password := c.PostForm("host_password")
	zone := c.PostForm("zone")
	port := c.PostForm("port")
	cli := common.NewSsh(ip, host_user, host_password, 22)
	tmp, err := cli.Run(fmt.Sprintf("firewall-cmd --zone=%v --add-port=%v/tcp", zone, port))
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"指定区域开放端口失败")
		return
	}
	response.Success(c, gin.H{"tmp": tmp}, "指定区域开放端口成功")
}

func DelZonePort(c *gin.Context) { //zone = block dmz drop external home internal public trusted work
	ip := c.PostForm("ip")
	host_user := c.PostForm("host_user")
	host_password := c.PostForm("host_password")
	zone := c.PostForm("zone")
	port := c.PostForm("port")
	cli := common.NewSsh(ip, host_user, host_password, 22)
	tmp, err := cli.Run(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-port=%v/tcp", zone, port))
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"指定区域移除端口失败")
		return
	}
	response.Success(c, gin.H{"tmp": tmp}, "指定区域移除端口成功")
}

func AddZonePortPermanent(c *gin.Context) {
	ip := c.PostForm("ip")
	host_user := c.PostForm("host_user")
	host_password := c.PostForm("host_password")
	zone := c.PostForm("zone")
	port := c.PostForm("port")
	cli := common.NewSsh(ip, host_user, host_password, 22)
	tmp, err := cli.Run(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-port=%v/tcp", zone, port))
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"永久在该区域开放端口失败")
		return
	}
	response.Success(c, gin.H{"tmp": tmp}, "永久在该区域开放端口成功")
}
