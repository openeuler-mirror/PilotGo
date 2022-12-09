package controller

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func AlertMessageConfigHandler(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	logger.Info(string(j))
	if err != nil {
		logger.Error("%s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	err = service.AlertMessageConfig(j)
	if err != nil {
		logger.Error("%s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "success")
}
