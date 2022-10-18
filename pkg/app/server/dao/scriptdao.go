package dao

import (
	"fmt"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
)

//添加脚本文件
func AddScript(s model.Script) error {
	fmt.Println("jjinru daoceng")
	version := s.Version
	if len(version) == 0 {
		return fmt.Errorf("版本号不能为空")
	}
	global.PILOTGO_DB.Save(&s)
	return nil
}

//根据脚本版本号查询文件是否存在
func IsVersionExist(scriptversion string) bool {
	var script model.Script
	global.PILOTGO_DB.Where("version=?", scriptversion).Find(&script)
	return script.Deleted == 0
}

//根据版本号删除文件（将标志位变为1）
func DeleteScript(scriptversion string) error {
	var script model.Script
	if IsVersionExist(scriptversion) {
		if err := global.PILOTGO_DB.Model(&script).Where("version=?", scriptversion).Update("deleted", 1).Error; err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("脚本不存在")
}
