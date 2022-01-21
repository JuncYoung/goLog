package goLog

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	logInfoPrefix = "UNIQUE-DialogSessionID||"

	DebugLevel = 0
	WarnLevel  = 1
	ErrorLevel = 2
)

var (
	sysLog     = &logrus.Logger{}
	onceSysLog = sync.Once{}
)

func GetSysLogger() *logrus.Logger {
	onceSysLog.Do(func() {
		sysLog = logrus.New()
		// 设置日志格式为json格式
		sysLog.SetFormatter(&logrus.JSONFormatter{})
		sysLog.SetOutput(os.Stdout)
		// 设置日志级别为info以上
		sysLog.SetReportCaller(false)
		sysLog.SetLevel(logrus.DebugLevel)
	})

	return sysLog
}

func LogPrintf(level int, format string, parameters ...interface{}) {
	format = fmt.Sprintf("%s  %s", logInfoPrefix, format)
	switch level {
	case 0:
		LogDebug(format, parameters...)
	case 1:
		LogWarnf(format, parameters...)
	case 2:
		LogErrorf(format, parameters...)
	default:
		LogDebug(format, parameters...)
	}
}

func LogPrintfWithID(level int, uid string, format string, parameters ...interface{}) {
	format = fmt.Sprintf("%s  %s", logInfoPrefix+uid, format)
	switch level {
	case 0:
		LogDebug(format, parameters...)
	case 1:
		LogWarnf(format, parameters...)
	case 2:
		LogErrorf(format, parameters...)
	default:
		LogDebug(format, parameters...)
	}
}

func SetSysLevel(level int) {
	GetSysLogger().SetLevel(parseLogLevel(level))
}

func LogDebug(format string, parameters ...interface{}) {
	GetSysLogger().Debugf(format, parameters...)
}

func LogWarnf(format string, parameters ...interface{}) {
	GetSysLogger().Warnf(format, parameters...)
}

func LogErrorf(format string, parameters ...interface{}) {
	GetSysLogger().Errorf(format, parameters...)
}

func parseLogLevel(level int) logrus.Level {
	switch level {
	case 0:
		return logrus.DebugLevel
	case 1:
		return logrus.WarnLevel
	case 2:
		return logrus.ErrorLevel
	}
	return logrus.DebugLevel
}
