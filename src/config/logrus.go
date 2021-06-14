package config

import "github.com/sirupsen/logrus"

var log = logrus.New()

// Log 日志格式
func Log() *logrus.Logger {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return log
}
