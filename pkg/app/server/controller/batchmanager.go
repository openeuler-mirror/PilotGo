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

func RemoveRepeatedElement(s []string) []string {
	result := make([]string, 0)
	m := make(map[string]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}
func CreateBatch(c *gin.Context) {
	// buf := make([]byte, 1024)
	// n, _ := c.Request.Body.Read(buf)
	// c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf[:n]))
	// j := buf[0:n]
	j, err := ioutil.ReadAll(c.Request.Body)
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
	batchinfo.DepartID = RemoveRepeatedElement(batchinfo.DepartID)
	batchinfo.DepartName = RemoveRepeatedElement(batchinfo.DepartName)
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

type BatchId struct {
	Page int `json:"page"`
	Size int `json:"size"`
	ID   int `json:"ID"`
}

func Batchmachineinfo(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var batchid BatchId
	err = json.Unmarshal(j, &batchid)
	logger.Info("%+v", batchid)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	// tmp, err := strconv.Atoi(batchid.ID)
	// if err != nil {
	// 	response.Response(c, http.StatusUnprocessableEntity,
	// 		422,
	// 		nil,
	// 		"批次ID有误")
	// 	return
	// }
	machinelist := dao.GetMachineID(batchid.ID)
	logger.Info("%+v", MachineInfo)
	MachineInfo := make([]model.MachineNode, 0)
	for _, value := range machinelist {
		// tmp1, err := strconv.Atoi(value)
		// if err != nil {
		// 	response.Response(c, http.StatusUnprocessableEntity,
		// 		422,
		// 		nil,
		// 		"批次ID有误")
		// 	return
		// }
		m := dao.MachineData(value)
		MachineInfo = append(MachineInfo, m)
	}

	len := len(MachineInfo)
	size := batchid.Size
	page := batchid.Page

	if len == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"data":  MachineInfo,
			"page":  page,
			"size":  size,
			"total": len,
		})
	}
	// page, _ := strconv.Atoi(batchid.Page)
	// size, _ := strconv.Atoi(batchid.Size)

	num := size * (page - 1)
	if num > len {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"页码超出")
	}

	if page*size >= len {
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"data":  MachineInfo[num:],
			"page":  page,
			"size":  size,
			"total": len,
		})
	} else {
		if page*size < num {
			response.Response(c, http.StatusUnprocessableEntity,
				422,
				nil,
				"读取错误")
		}

		if page*size == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":  200,
				"data":  MachineInfo,
				"page":  page,
				"size":  size,
				"total": len,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":  200,
				"data":  MachineInfo[num : page*size-1],
				"page":  page,
				"size":  size,
				"total": len,
			})
		}

	}
}
