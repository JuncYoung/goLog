package goLog

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultDir     = "./logs"
	defaultLogName = "qnLog"
)

func (l Logger) ParseQnApiLogInput(data interface{}) {
	s := standardLogFormat{
		Service: l.ServiceName,
		From:    data,
	}

	bs, err := json.Marshal(s)
	if err != nil {
		LogPrintf(ErrorLevel, "qnApiLog Marshal err:%s", err.Error())
	}
	str := string(bs)
	l.Info(str)
}

func (l Logger) ParseQnApiLogOutput(t time.Time, data interface{}) {
	s := standardLogFormat{
		Service: l.ServiceName,
		Out:     data,
		Time:    time.Now().Sub(t).Seconds(),
	}

	bs, err := json.Marshal(s)
	if err != nil {
		LogPrintf(ErrorLevel, "qnApiLog Marshal err:%s", err.Error())
	}
	str := string(bs)
	l.Info(str)
}

func (l Logger) QnInternalError(data interface{}) {
	s := standardLogFormat{
		Service: l.ServiceName,
		From:    data,
	}

	bs, err := json.Marshal(s)
	if err != nil {
		LogPrintf(ErrorLevel, "err log Marshal err:%s", err.Error())
	}
	str := string(bs)
	l.Error(str)
}

func (l Logger) Method(name string) Logger {
	l.ServiceName = name
	return l
}

func SetupLoggerByDate(logDir, logName string, rotateMaxAge, skip int, report bool, level int) error {
	if logDir == "" {
		logDir = defaultDir
	}
	LogPrintf(DebugLevel, "SetupLoggerByDate using log dir : [%s]\n", logDir)

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

func SetupQnFormatByDate(logName, serviceName string) error {
	qnSets := sets.QnLog
	logDir := qnSets.LogDir
	if logDir == "" {
		logDir = defaultDir
	}
	if logName == "" {
		logName = defaultLogName
	}
	LogPrintf(DebugLevel, "SetupQnFormatByDate using log dir : [%s]\n", logDir)
	log := logrus.New()
	return setupDIYFormatByDate(log, parseLogLevel(qnSets.Level), logName, logDir, serviceName, qnSets.RotateMaxAge, qnSets.Skip, qnSets.Report)
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
	// ?????????????????????json??????
	logs.SetFormatter(&QnFormatter{})
	// ???????????????????????????
	logs.SetReportCaller(false)
	// ????????????????????????????????????io.writer??????
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
