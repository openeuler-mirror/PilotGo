package controller

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 存储脚本文件
func AddScript(c *gin.Context) {
	var script model.Script
	c.Bind(&script)
	if len(script.Name) == 0 {
		response.Fail(c, nil, "请输入脚本文件名字")
		return
	}
	if len(script.Content) == 0 {
		response.Fail(c, nil, "请输入脚本内容")
		return
	}
	if len(script.Description) == 0 {
		response.Fail(c, nil, "请输入脚本描述")
		return
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	version := time.Now().Format("2006-01-02 15:04:05") + "-" + vcode
	fmt.Println(vcode)
	fmt.Println(version)
	sc := model.Script{
		Name:        script.Name,
		Content:     script.Content,
		Description: script.Description,
		UpdatedAt:   time.Time{},
		Version:     version,
		Deleted:     0,
	}
	err := dao.AddScript(sc)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "脚本文件添加失败")
	}
	response.Success(c, nil, "脚本文件添加成功")
}
