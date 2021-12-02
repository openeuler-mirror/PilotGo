package os

import (
  "fmt"
  "openeluer.org/PilotGo/PilotGo/pkg/utils"
)

type ZonePort struct {
  Zone string
  Port int
}

func Restart() bool {
  tmp, _ := utils.RunCommand("service firewalld restart")
  if len(tmp) != 0 {
    fmt.Println("重启防火墙失败！")
    return false
  }
  return true
}

func Config() bool {
  tmp, _ := utils.RunCommand("firewall-cmd --list-all")
  if len(tmp) != 0 {
    fmt.Println("重启防火墙失败！")
    return false
  }
  return true
}

func Reload() bool {
  tmp, _ := utils.RunCommand("firewall-cmd --reload")
  if len(tmp) != 0 {
    fmt.Println("更新防火墙失败！")
    return false
  }
  return true
}

func Stop() bool {
  tmp, _ := utils.RunCommand("service firewalld stop")
  if len(tmp) != 0 {
    fmt.Println("关闭防火墙失败！")
    return false
  }
  return true
}

func DelZonePort(zp *ZonePort) string { //zone = block dmz drop external home internal public trusted work
  tmp, _ := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=public --remove-port=%v/tcp", zp.Port))
  return tmp
}

func AddZonePortPermanent(zp *ZonePort) string { //zone = block dmz drop external home internal public trusted work
  tmp, _ := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=public --add-port=%v/tcp", zp.Port))
  return tmp
}
