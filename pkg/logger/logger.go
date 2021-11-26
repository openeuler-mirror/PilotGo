package logger

import (
	"errors"
	"os"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	conf "openeluer.org/PilotGo/PilotGo/pkg/config"
)

var logName string = "pilotgo"

func setLogDriver(logopts *conf.LogOpts) error {
	if logopts == nil {
		return errors.New("logopts == nil")
	}

	switch logopts.LogDriver {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		writer, err := rotatelogs.New(
			logopts.LogPath+"/"+logName,
			rotatelogs.WithRotationCount(uint(logopts.MaxFile)),
			rotatelogs.WithRotationSize(int64(logopts.MaxSize)),
		)
		if err != nil {
			return err
		}
		log.SetOutput(writer)
	}
	return nil
}

func setLogLevel(logopts *conf.LogOpts) error {
	switch logopts.LogLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return nil
}
func Init(conf *conf.Configure) error {
	setLogLevel(&(conf.Logopts))
	err := setLogDriver(&(conf.Logopts))
	if err != nil {
		return err
	}
	log.Debug("log init")

	return nil
}

func Trace(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
