/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"fmt"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
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
