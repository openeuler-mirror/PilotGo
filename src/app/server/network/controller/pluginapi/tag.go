package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/tag"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 插件返回tag数据
func GetTagHandler(c *gin.Context) {
	msg := &tag.TageMessage{}
	if err := c.ShouldBind(msg); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "get tag",
	}
	auditlog.Add(log)
	//TODO:获取到了tag数据开始使用
	response.Success(c, gin.H{"data": msg}, "get tag成功")
}
