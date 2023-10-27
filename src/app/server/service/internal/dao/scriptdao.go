package dao

import (
	"fmt"
	"time"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
)

type Script struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	Description string `json:"description"`
	UpdatedAt   time.Time
	Version     string `gorm:"unique" json:"version"`
	Deleted     int    `json:"deleted"` //deleted为1的时候表示删除，一般表示为0
}

// 添加脚本文件
func AddScript(s Script) error {
	fmt.Println("jjinru daoceng")
	version := s.Version
	if len(version) == 0 {
		return fmt.Errorf("版本号不能为空")
	}
	return mysqlmanager.MySQL().Save(&s).Error
}

// 根据脚本版本号查询文件是否存在
func IsVersionExist(scriptversion string) (bool, error) {
	var script Script
	err := mysqlmanager.MySQL().Where("version=?", scriptversion).Find(&script).Error
	return script.Deleted == 0, err
}

// 根据版本号删除文件（将标志位变为1）
func DeleteScript(scriptversion string) error {
	var script Script
	VersionExistBool, err := IsVersionExist(scriptversion)
	if err != nil {
		return err
	}
	if VersionExistBool {
		if err := mysqlmanager.MySQL().Model(&script).Where("version=?", scriptversion).Update("deleted", 1).Error; err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("脚本不存在")
}

// 根据版本号查询脚本文件内容
func ShowScript(scriptversion string) (string, error) {
	var script Script
	err := mysqlmanager.MySQL().Where("version=?", scriptversion).Find(&script).Error
	return script.Content, err
}
