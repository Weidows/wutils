package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	logger = &logrus.Logger{
		// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
		// 日志消息输出可以是任意的io.writer类型
		Out: os.Stdout,
		//Formatter: new(logrus.TextFormatter),
		Formatter: &logrus.JSONFormatter{
			//TimestampFormat:   "yyyy-MM-dd HH:mm:ss.ffffff",
			//TimestampFormat:   "HH.mm.ss",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			FieldMap:          nil,
			CallerPrettyfier:  nil,
			PrettyPrint:       true,
		},
		Hooks: make(logrus.LevelHooks),
		// 设置日志级别为warn以上
		//Level: logrus.DebugLevel,
		Level: logrus.WarnLevel,
	}
)

func GetLogger() *logrus.Logger {
	return logger
}
