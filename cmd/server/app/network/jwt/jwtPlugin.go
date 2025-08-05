package jwt

import (
	"errors"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type PluginServiceClaims struct {
	jwt.StandardClaims

	ServiceName string
}

func GeneratePluginServiceToken(serviceName string) (string, error) {
	claims := &PluginServiceClaims{
		ServiceName: serviceName,

		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			Issuer:   Issue,
			Subject:  "plugin token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.OptionsConfig.JWT.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParsePluginServiceClaims(c *gin.Context) (*PluginServiceClaims, error) {
	cookie, err := c.Request.Cookie("PluginToken")
	if err != nil {
		return nil, err
	}

	claims, err := parseClaims(cookie.Value, &PluginServiceClaims{})
	if err != nil {
		return nil, err
	}
	m, ok := claims.(*PluginServiceClaims)
	if !ok {
		return nil, errors.New("invalid plugin claims")
	}
	return m, nil
}
