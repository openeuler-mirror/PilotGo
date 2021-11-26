package controller

/**
 * @Author: zhang han
 * @Date: 2021/10/28 14:58
 * @Description: 用户登录和注册
 */

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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

	user := model.User{ //Create user
		Username: username,
		Password: password,
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
	if password != user.Password {
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
