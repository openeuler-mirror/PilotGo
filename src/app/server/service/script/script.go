package script

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
)

type Script = dao.Script

// 存储脚本文件
func AddScript(script *dao.Script) error {
	if len(script.Name) == 0 {
		return errors.New("请输入脚本文件名字")
	}
	if len(script.Content) == 0 {
		return errors.New("请输入脚本内容")
	}
	if len(script.Description) == 0 {
		return errors.New("请输入脚本描述")
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	version := time.Now().Format("2006-01-02 15:04:05") + "-" + vcode
	sc := dao.Script{
		Name:        script.Name,
		Content:     script.Content,
		Description: script.Description,
		UpdatedAt:   time.Time{},
		Version:     version,
		Deleted:     0,
	}
	err := dao.AddScript(sc)
	if err != nil {
		return errors.New("脚本文件添加失败")
	}
	return nil
}
