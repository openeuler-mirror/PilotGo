package logger

import (
	"github.com/sirupsen/logrus"
)

type LoggerConf struct {
	Level   int
	Path    string
	MaxSize int
	MaxFile int
	Driver  string
}

type Logger struct {
	config *LoggerConf
	inst   *logrus.Logger
}

func (l *Logger) Init() {

	l.inst = &logrus.Logger{}

	// TODO： 创建路径、初始化日志文件等
}

func (l *Logger) SetLevel() {

}

func (l *Logger) Trace() {}

func (l *Logger) Debug() {}

func (l *Logger) Info() {}

func (l *Logger) Warn() {}

func (l *Logger) Error() {}

func (l *Logger) Fatal() {}

func CreateLogger(conf *LoggerConf) *Logger {
	l := &Logger{
		config: conf,
	}
	l.Init()

	return l
}
