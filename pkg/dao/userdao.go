package dao

/**
 * @Author: zhang han
 * @Date: 2021/10/30 15:22
 * @Description:
 */

import (
	"openeluer.org/PilotGo/PilotGo/pkg/model"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func IsEmailExist(email string) bool {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

//func GetUserList(c *gin.Context) {
//	var pageInfo model.PageInfo
//
//	pageStr := c.Query("page")
//	page, err := strconv.ParseInt(pageStr, 10, 64)
//	if err != nil {
//		response.Fail(c,
//			"查询失败",
//			nil)
//		return
//	}
//
//	pageSizeStr := c.Query("pageSize")
//	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
//	if err != nil {
//		response.Fail(c,
//			"查询失败",
//			nil)
//		return
//	}
//
//	pageInfo.Page = int(page)
//	pageInfo.PageSize = int(pageSize)
//
//	err, list, total := service.GetUserInfoList(pageInfo)
//	if err != nil {
//		response.Fail(c,
//			"查询失败",
//			nil)
//	} else {
//		response.Success(c, gin.H{}, model.PageResult{
//			List:     list,
//			Total:    total,
//			Page:     pageInfo.Page,
//			PageSize: pageInfo.PageSize,
//		})
//	}
//}
//
//func GetUser(c *gin.Context) {
//
//	idVal, ok := c.Params.Get("id")
//	if !ok {
//		response.FailWithMessage(fmt.Sprintf("id格式错误"), c)
//		return
//	}
//	id, err := strconv.ParseUint(idVal, 10, 64)
//	if err != nil {
//		response.FailWithMessage(fmt.Sprintf("id格式错误，%v", err), c)
//		return
//	}
//
//	err, user := service.GetUserById(uint(id))
//	if err != nil {
//		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
//	} else {
//		response.OkWithData(user, c)
//	}
//}
