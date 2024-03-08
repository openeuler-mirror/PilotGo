package middleware

import (
	"net/http"

	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	user, err := jwt.ParseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "claims解析错误:" + err.Error()})
		c.Abort()
		return
	}

	ok, err := auth.CheckAuth(user.Username, c.Request.RequestURI, "get", auth.DomainPilotGo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "check auth error:" + err.Error()})
		c.Abort()
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "no permission"})
		c.Abort()
		return
	}

	c.Next()
}

func NeedPermission(resource, action string) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := jwt.ParseUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "claims解析错误:" + err.Error()})
			c.Abort()
			return
		}

		ok, err := auth.CheckAuth(user.Username, resource, action, auth.DomainPilotGo)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "check auth error:" + err.Error()})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "no permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func TokenCheckMiddleware(c *gin.Context) {
	_, err := jwt.ParseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "user token check error:" + err.Error()})
		c.Abort()
		return
	}

	c.Next()
}
