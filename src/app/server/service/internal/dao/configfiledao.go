package dao

import (
	"fmt"
	"time"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
)

type ConfigFile struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	MachineUUID string `json:"uuid"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	UpdatedAt   time.Time
}

func AddConfigFile(cf ConfigFile) error {
	node, err := MachineInfoByUUID(cf.MachineUUID)
	if err != nil {
		return err
	}
	if node.ID != 0 {
		return mysqlmanager.MySQL().Save(&cf).Error
	}
	return fmt.Errorf("机器不存在")
}
