package logger

import (
	"errors"
	"os"
	"path"
	"runtime"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	conf "openeluer.org/PilotGo/PilotGo/pkg/config"
)

var logName string = "pilotgo"

func setLogDriver(logopts *conf.LogOpts) error {
	if logopts == nil {
		return errors.New("logopts == nil")
	}

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := path.Base(f.File)
			return f.Function, fileName
		},
	})

	switch logopts.LogDriver {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "file":
		writer, err := rotatelogs.New(
			logopts.LogPath+"/"+logName,
			rotatelogs.WithRotationCount(uint(logopts.MaxFile)),
			rotatelogs.WithRotationSize(int64(logopts.MaxSize)),
		)
		if err != nil {
			return err
		}
		logrus.SetOutput(writer)
	default:
		logrus.SetOutput(os.Stdout)
		logrus.Warn("!!! invalid log output, use stdout !!!")
	}
	return nil
}

func setLogLevel(logopts *conf.LogOpts) error {
	switch logopts.LogLevel {
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
		logrus.SetLevel(logrus.InfoLevel)
	}
	return nil
}
func Init(conf *conf.Configure) error {
	setLogLevel(&(conf.Logopts))
	err := setLogDriver(&(conf.Logopts))
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
