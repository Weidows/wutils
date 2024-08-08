package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger = &logrus.Logger{
		// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
		// 日志消息输出可以是任意的io.writer类型
		Out: os.Stdout,
		//Formatter: new(logrus.TextFormatter),
		Formatter: &logrus.JSONFormatter{
			TimestampFormat:   "15:03:04",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			FieldMap:          nil,
			CallerPrettyfier:  nil,
			PrettyPrint:       true,
		},
		Hooks: make(logrus.LevelHooks),
		// 设置日志级别为warn以上
		//Level: logrus.WarnLevel, 这个会导致错误
		Level: logrus.DebugLevel,
	}
)

func GetLogger() *logrus.Logger {
	logger.SetReportCaller(true)
	return logger
}
