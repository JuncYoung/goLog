package goLog

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func (l Logger) ParseQnApiLogFormat(isInput bool, t int64, data interface{}) {
	s := standardLogFormat{
		Service: l.ServiceName,
	}
	if isInput {
		s.From = data
	} else {
		s.Out = data
		s.Time = float32(t) / float32(time.Second)
	}
	bs, err := json.Marshal(s)
	if err != nil {
		LogPrintf(ErrorLevel, "ParseQnApiLogFormat Marshal err:%s", err.Error())
	}
	str := string(bs)
	l.Info(str)
}

func SetupLoggerByDate(logDir, logName string, rotateMaxAge, skip int, report bool, level int) error {
	LogPrintf(DebugLevel, "SetupLoggerCommonByDate using log dir : [%s]\n", logDir)

	log := logrus.New()
	return setupLoggerByDate(log, parseLogLevel(level), logName, logDir, rotateMaxAge, skip, report)
}

func setupLoggerByDate(logs *logrus.Logger, level logrus.Level, fileName, logDir string, rotateMaxAge, skip int, report bool) error {
	if logs == nil {
		return ErrEmptyLog
	}
	isExist, err := pathExists(logDir)
	if err != nil {
		return fmt.Errorf("failed to create trace output file: %v", err)
	}
	if !isExist {
		err = os.MkdirAll(logDir, 0755)

		if err != nil {
			return fmt.Errorf("failed to create trace output dir: %v", err)
		}
	}

	dateStr := time.Now().Format("20060102")
	logFullPath := logDir + "/" + fileName + "." + dateStr
	// os.O_O_CREATE auto create
	logFile, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("open logger file fail: %s", err.Error())
	}
	fileWriter := LogFileWriterByDate{
		dirPath:      logDir,
		fileName:     fileName,
		file:         logFile,
		dateStr:      dateStr,
		mutex:        &sync.Mutex{},
		rotateMaxAge: rotateMaxAge,
	}
	fileWriter.mustDoMillFirst()
	logs.SetFormatter(&logrus.JSONFormatter{})
	logs.SetOutput(&fileWriter)
	logs.SetReportCaller(false)
	if report {
		logs.AddHook(&DateLogHook{skip: skip})
	}

	logs.SetLevel(level)

	GetFileLoggerMgr().AddOne(fileName, Logger{Logger: logs})
	return nil
}

func SetupQnFormatByDate(logDir, logName, serviceName string, rotateMaxAge, skip int, report bool, level int) error {
	LogPrintf(DebugLevel, "SetupQnFormatByDate using log dir : [%s]\n", logDir)

	log := logrus.New()
	return setupDIYFormatByDate(log, parseLogLevel(level), logName, logDir, serviceName, rotateMaxAge, skip, report)
}

func setupDIYFormatByDate(logs *logrus.Logger, level logrus.Level, fileName, logDir, serviceName string, rotateMaxAge, skip int, report bool) error {
	if logs == nil {
		return ErrEmptyLog
	}
	isExist, err := pathExists(logDir)
	if err != nil {
		return fmt.Errorf("failed to create trace output file: %v", err)
	}
	if !isExist {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create trace output dir: %v", err)
		}
	}

	dateStr := time.Now().Format("20060102")
	logFullPath := logDir + "/" + fileName + "." + dateStr
	// os.O_O_CREATE auto create
	logFile, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("open logger file fail: %s", err.Error())
	}

	fileWriter := LogFileWriterByDate{
		dirPath:      logDir,
		fileName:     fileName,
		file:         logFile,
		dateStr:      dateStr,
		mutex:        &sync.Mutex{},
		rotateMaxAge: rotateMaxAge,
	}
	fileWriter.mustDoMillFirst()
	// 设置日志格式为json格式
	logs.SetFormatter(&QnFormatter{})
	// 开启调用上下文记录
	logs.SetReportCaller(false)
	// 日志消息输出可以是任意的io.writer类型
	logs.SetOutput(&fileWriter)
	if report {
		logs.AddHook(&DateLogHook{skip: skip})
	}

	logs.SetLevel(level)

	GetFileLoggerMgr().AddOne(qnLog, Logger{
		Logger:      logs,
		ServiceName: serviceName,
	})

	return nil
}
