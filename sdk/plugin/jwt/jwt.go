/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Aug 07 16:18:53 2025 +0800
 */
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var Issue = "PilotGo"

const (
	TokenCookie = "PluginToken"
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
	tokenString, err := token.SignedString([]byte(serviceName))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParsePluginServiceClaims(c *gin.Context) (*PluginServiceClaims, error) {
	cookie, err := c.Request.Cookie(TokenCookie)
	if err != nil {
		return nil, err
	}

	jwtSecretKey, err := getJWTSecretKeyForService(cookie.Value)
	if err != nil {
		return nil, err
	}

	claims, err := parseClaims(cookie.Value, jwtSecretKey, &PluginServiceClaims{})
	if err != nil {
		return nil, err
	}
	m, ok := claims.(*PluginServiceClaims)
	if !ok {
		return nil, errors.New("invalid plugin claims")
	}
	return m, nil
}

func parseClaims(tokenString string, jwtSecretKey string, claims jwt.Claims) (jwt.Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if token != nil && !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token.Claims, nil
}

func getJWTSecretKeyForService(tokenString string) (string, error) {
	parser := new(jwt.Parser)
	tempClaims := &PluginServiceClaims{}
	_, _, err := parser.ParseUnverified(tokenString, tempClaims)
	if err != nil {
		return "", fmt.Errorf("failed to parse token without verifying: %s", err.Error())
	}
	return tempClaims.ServiceName, nil
}
