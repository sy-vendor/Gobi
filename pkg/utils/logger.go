package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)

	// 支持日志写入文件和控制台
	file, err := os.OpenFile("gobi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(io.MultiWriter(os.Stdout, file))
	}
}
