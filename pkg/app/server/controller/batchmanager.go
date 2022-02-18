package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"

	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

type BatchInfo struct {
	Name       string
	Descrip    string
	Manager    string
	DepartName []string
	DepartID   []string
	Machine    []string
}

func CreateBatch(c *gin.Context) {
	// buf := make([]byte, 1024)
	// n, _ := c.Request.Body.Read(buf)
	// c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf[:n]))
	// j := buf[0:n]
	// fmt.Println("body:", string(j)) //获取到post传递过来的数据
	j, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println("body:", string(j))
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var batchinfo BatchInfo
	err = json.Unmarshal(j, &batchinfo)
	logger.Info("%+v", batchinfo)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	fmt.Println("====>" + batchinfo.Name)
	if len(batchinfo.Name) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"批次名未输入")
		return
	}
	if len(batchinfo.Manager) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"创建人未输入")
		return
	}
	if len(batchinfo.DepartID) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"归属部门未输入")
		return
	}
	if len(batchinfo.Machine) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"归属部门未输入")
		return
	}
	if dao.IsExistName(batchinfo.Name) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"已存在该名称批次")
		return
	}
	for _, value := range batchinfo.DepartID {
		tmp, err := strconv.Atoi(value)
		if err != nil {
			response.Response(c, http.StatusUnprocessableEntity,
				422,
				nil,
				"部门ID有误")
			return
		}
		if !dao.IsDepartIDExist(tmp) {
			response.Response(c, http.StatusUnprocessableEntity,
				422,
				nil,
				"不存在该部门")
			return
		}
	}

	Batch := model.Batch{
		Name:        batchinfo.Name,
		Description: batchinfo.Descrip,
		Manager:     batchinfo.Manager,
		Depart:      strings.Join(batchinfo.DepartID, ","),
		DepartName:  strings.Join(batchinfo.DepartName, ","),
		Machinelist: strings.Join(batchinfo.Machine, ","),
	}
	mysqlmanager.DB.Create(&Batch)
	logger.Info("%s", batchinfo.Machine)
	response.Success(c, nil, "批次入库成功")
}

func BatchInform(c *gin.Context) {
	batch := model.Batch{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	list, total, err := batch.ReturnBatch(query)
	logger.Info("%+v", list)
	if model.HandleError(c, err) {
		return
	}
	// 返回数据开始拼装分页的json
	model.JsonPagination(c, list, total, query)
}

type Batchdel struct {
	BatchID []string
}

func DeleteBatch(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println("body:", string(j))
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var batchdel Batchdel
	err = json.Unmarshal(j, &batchdel)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	if len(batchdel.BatchID) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"请输入删除批次ID")
		return
	}
	for _, value := range batchdel.BatchID {
		tmp, err := strconv.Atoi(value)
		if err != nil {
			response.Response(c, http.StatusUnprocessableEntity,
				422,
				nil,
				err.Error())
			return
		}
		dao.DeleteBatch(tmp)
	}
	response.Success(c, nil, "批次删除成功")
}

type Batchupdate struct {
	BatchID   string
	BatchName string
	Descrip   string
}

func UpdateBatch(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println("body:", string(j))
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var batchinfo Batchupdate
	err = json.Unmarshal(j, &batchinfo)
	logger.Info("%+v", batchinfo)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	if len(batchinfo.BatchID) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"请输入修改批次ID")
		return
	}
	tmp, err := strconv.Atoi(batchinfo.BatchID)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID有误")
		return
	}
	if !dao.IsExistID(tmp) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"不存在该批次")
		return
	}
	dao.UpdateBatch(tmp, batchinfo.BatchName, batchinfo.Descrip)
	response.Success(c, nil, "批次修改成功")
}
