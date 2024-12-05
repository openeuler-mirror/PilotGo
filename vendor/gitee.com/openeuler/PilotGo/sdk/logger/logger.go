/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
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
	Driver  string `yaml:"driver" comment:"可选stdout和file.stdout:输出到终端控制台;file:输出到path下的指定文件。"`
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

func ErrorStack(msg string, err error) {
	logrus.Errorf(msg+"\n%+v", err)
}

func Fatal(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func RequestLogger(_skipPaths []string) gin.HandlerFunc {
	var skip map[string]struct{}
	if len(_skipPaths) > 0 {
		skip = make(map[string]struct{}, len(_skipPaths))

		for _, path := range _skipPaths {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if _, ok := skip[path]; !ok {
			endTime := time.Now()
			latency := endTime.Sub(start)
			method := c.Request.Method
			statusCode := c.Writer.Status()
			clientIP := c.ClientIP()
			errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}
			if latency > time.Minute {
				latency = latency.Truncate(time.Second)
			}

			Debug("%s  status_code:%3d  latency_time:%-13v client_ip:%-15s  method:%-7s req_uri:%#v\n%s",
				start.Format("2006-01-02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
				errorMessage)
		}
	}
}
