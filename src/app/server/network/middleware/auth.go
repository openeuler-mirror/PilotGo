package middleware

import (
	"net/http"

	"gitee.com/openeuler/PilotGo/app/server/service/auth"
	"gitee.com/openeuler/PilotGo/app/server/service/jwt"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwt.ParseMyClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "claims解析错误:" + err.Error()})
		c.Abort()
		return
	}
	logger.Debug("request from %d, %s", claims.UserId, claims.UserName)

	ok, err := auth.CheckAuth(claims.UserName, c.Request.RequestURI, "get")
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
		claims, err := jwt.ParseMyClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "claims解析错误:" + err.Error()})
			c.Abort()
			return
		}
		logger.Debug("request from %d, %s", claims.UserId, claims.UserName)

		ok, err := auth.CheckAuth(claims.UserName, resource, action)
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
	_, err := jwt.ParseMyClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "user token check error:" + err.Error()})
		c.Abort()
		return
	}

	c.Next()
}
