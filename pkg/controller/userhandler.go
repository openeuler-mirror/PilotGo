package controller

/**
 * @Author: zhang han
 * @Date: 2021/10/28 14:58
 * @Description: 用户登录和注册
 */

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/common/dto"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/model"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

func Register(c *gin.Context) {

	username := c.PostForm("username") //Get the parameter
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	enable := c.PostForm("enable")

	if len(username) == 0 { //Data verification
		username = utils.RandomString(5)
	}
	if len(password) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"密码不能为空!")
		return
	}
	if len(email) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"邮箱不能为空!")
		return
	}
	if dao.IsEmailExist(email) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"Email exist!")
		return
	}
	hasedPassword, err := common.HashAndSalt(password)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Hased password error!")
		return
	}
	user := model.User{ //Create user
		Username: username,
		Password: string(hasedPassword),
		Phone:    phone,
		Email:    email,
		Enable:   enable,
	}
	mysqlmanager.DB.Create(&user)

	response.Success(c, nil, "注册成功!") //Return result
}

func Login(c *gin.Context) {

	email := c.PostForm("email") //get the argument
	password := c.PostForm("password")
	/*Data verification*/
	var user model.User //Data verification
	mysqlmanager.DB.Where("email = ?", email).Find(&user)

	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"用户不存在!")
		return
	}
	check := common.ComparePasswords(user.Password, password)
	if !check {
		response.Response(c, http.StatusBadRequest,
			400,
			nil,
			"密码错误!")
		return
	}

	token, err := common.ReleaseToken(user) // Issue token
	if err != nil {
		response.Response(c, http.StatusInternalServerError,
			500,
			nil,
			"服务器内部错误!")
		log.Printf("token生成错误:%v", err)
		return
	}
	response.Success(c, gin.H{"token": token}, "登陆成功!")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

// 查询所有用户
func UserAll(c *gin.Context) {
	users := model.User{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	list, total, err := users.All(query)
	if model.HandleError(c, err) {
		return
	}
	// 返回数据开始拼装分页的json
	model.JsonPagination(c, list, total, query)
}

// 刷新
func UserRefresh(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if model.HandleError(c, err) {
		return
	}
	err = user.Refresh()
	if model.HandleError(c, err) {
		return
	}
	response.Success(c, nil, "Updated User successfully!")
}

// 删除用户
func DeleteUser(c *gin.Context) {
	var user model.User
	userEmail := c.PostForm("email")
	if dao.IsEmailExist(userEmail) {
		mysqlmanager.DB.Where("email=?", userEmail).Delete(user)
		response.Response(c, http.StatusUnprocessableEntity,
			200,
			nil,
			"User deleted successfully!")
		return
	} else {
		response.Fail(c, nil, "No user found!")
	}
}

//修改用户信息
func UpdateUser(c *gin.Context) {
	var user model.User
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	if dao.IsEmailExist(email) {
		// 修改手机号
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("phone", phone)
		hasedPassword, err := common.HashAndSalt(password)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "Hased password error!")
			return
		}

		//修改密码
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("password", hasedPassword)
		response.Response(c, http.StatusUnprocessableEntity,
			200,
			gin.H{"data": user},
			"User update successfully!")
		return
	} else {
		response.Fail(c, nil, "No user found!")
	}
}

//一键导入用户数据
func ImportUser(c *gin.Context) {
	form, _ := c.MultipartForm()

	files := form.File["upload"]
	if files == nil {
		response.Fail(c, nil, "Please select a file first!")
		return
	}
	filePath := "static/"
	for _, file := range files {
		name := file.Filename
		filename := filePath + name

		// c.SaveUploadedFile(file, filename)
		xlFile, error := xlsx.OpenFile(filename)
		if error != nil {
			return
		}
		for _, sheet := range xlFile.Sheets {
			for rowIndex, row := range sheet.Rows {
				//跳过第一行表头信息
				if rowIndex == 0 {
					continue
				}
				user := model.User{}
				user.Username = row.Cells[0].Value
				// 设置默认密码为123456
				hasedPassword, err := common.HashAndSalt("123456")
				if err != nil {
					response.Response(c, http.StatusInternalServerError, 500, nil, "Hased password error!")
					return
				}
				user.Password = hasedPassword
				user.Phone = row.Cells[1].Value
				user.Email = row.Cells[2].Value
				mysqlmanager.DB.Create(&user)
			}
		}
	}
	response.Response(c, http.StatusUnprocessableEntity,
		200,
		nil,
		"import success")
}
