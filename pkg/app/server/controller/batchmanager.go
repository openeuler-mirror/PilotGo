package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"

	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func CreateBatch(c *gin.Context) {
	name := c.Query("Name")
	descrip := c.PostForm("Description")
	manager := c.PostForm("Manager")
	depart := c.QueryArray("Depart")
	machine := c.QueryArray("Machine")
	if len(name) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"批次名未输入")
		return
	}
	if len(manager) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"创建人未输入")
		return
	}
	if len(depart) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"归属部门未输入")
		return
	}
	if len(machine) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"归属部门未输入")
		return
	}
	if dao.IsExistName(name) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"已存在该名称批次")
		return
	}
	Batch := model.Batch{
		Name:        name,
		Description: descrip,
		Manager:     manager,
		Depart:      strings.Join(depart, ","),
	}
	mysqlmanager.DB.Create(&Batch)
	response.Success(c, nil, "批次入库成功")
	for _, value := range machine {
		batchid := dao.GetBatchID(name)
		tmp := strconv.Itoa(int(batchid))
		x := dao.GetmachineBatch(value)
		if dao.GetmachineBatch(value) != "" {
			x += ","
		}
		dao.UpdatemachineBatch(value, x+tmp)
	}
	response.Success(c, nil, "机器已绑定批次")
}
