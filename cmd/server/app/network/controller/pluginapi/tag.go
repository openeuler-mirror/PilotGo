package pluginapi

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/tag"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 插件返回tag数据
func GetTagHandler(c *gin.Context) {
	uuidTags := &struct {
		UUIDS []string `json:"uuids"`
	}{}
	if err := c.ShouldBindJSON(&uuidTags); err != nil {
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
		Action:     "publish tags",
	}
	auditlog.Add(log)
	//TODO:获取到了tag数据开始使用
	data, err := tag.RequestTag(uuidTags.UUIDS)
	if err != nil {
		auditlog.UpdateMessage(log, "agentuuid:"+err.Error())
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	response.Success(c, data, "get tag成功")
}
