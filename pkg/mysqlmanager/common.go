package mysqlmanager

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MachInfo struct {
	Id               int       `json:"id"`
	Ip               string    `json:"ip"`
	SystemStatus     int       `json:"system_status"`
	SystemInfo       string    `json:"system_info"`
	SystemVersion    string    `json:"system_version"`
	Arch             string    `json:"arch"`
	InstallationTime time.Time `json:"installation_time"`
	MachineType      int       `json:"machine_type"`
}

type PluginInfo struct {
	Id          int
	Name        string
	Description string
	Url         string
	Port        string
	Protocol    string
	Version     string
}

type MysqlManager struct {
	ip       string
	port     int
	userName string
	passWord string
	dbName   string
	db       *gorm.DB
}

func (m *MachInfo) CheckStatus() int {
	return checkMachineOnline(m.Ip)
}

func checkMachineOnline(remote string) int {
	return 1
}
