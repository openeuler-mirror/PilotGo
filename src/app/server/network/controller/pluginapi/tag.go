package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/service/tag"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
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

	//TODO:获取到了tag数据开始使用
	data, err := tag.RequestTag(uuidTags.UUIDS)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	response.Success(c, data, "get tag成功")
}
