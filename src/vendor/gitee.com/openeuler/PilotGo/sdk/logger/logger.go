/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: yangzhao1
 * Date: 2022-03-01 09:59:30
 * LastEditTime: 2022-04-05 11:37:16
 * Description: provide agent log manager of pilotgo
 ******************************************************************************/
package logger

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type LogOpts struct {
	Level   string `yaml:"level"`
	Driver  string `yaml:"driver"`
	Path    string `yaml:"path"`
	MaxFile int    `yaml:"max_file"`
	MaxSize int    `yaml:"max_size"`
}

func setLogDriver(logopts *LogOpts) error {
	if logopts == nil {
		return errors.New("logopts is nil")
	}

	switch logopts.Driver {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "file":
		writer, err := rotatelogs.New(
			logopts.Path,
			rotatelogs.WithRotationCount(uint(logopts.MaxFile)),
			rotatelogs.WithRotationSize(int64(logopts.MaxSize)),
		)
		if err != nil {
			return err
		}
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:   true,
			ForceQuote:      false,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logrus.SetOutput(writer)
	default:
		logrus.SetOutput(os.Stdout)
		logrus.Warn("!!! invalid log output, use stdout !!!")
	}
	return nil
}

func setLogLevel(logopts *LogOpts) error {
	switch logopts.Level {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		return errors.New("invalid log level")
	}
	return nil
}
func Init(conf *LogOpts) error {
	setLogLevel(conf)
	err := setLogDriver(conf)
	if err != nil {
		return err
	}
	logrus.Debug("log init")

	return nil
}

func Trace(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

func Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		if reqUri != "/api/v1/pluginapi/heartbeat" {
			Debug("status_code:%d latency_time:%s client_ip:%s req_method:%s req_uri:%s",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
			)
		}
	}
}
