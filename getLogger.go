package goLog

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyLog = errors.New("log pointer empty")
)

func GetFileLogger(fileName string) Logger {
	if fileName == "" {
		LogPrintf(ErrorLevel, "get empty file name")
		return Logger{}
	}

	flog := GetFileLoggerMgr().getFileLogger(fileName)
	if flog.Logger == nil {
		if err := SetupLoggerByDate(sets.Com.LogDir, fileName, sets.Com.RotateMaxAge, sets.Com.Skip, sets.Com.Report, DebugLevel); err != nil {
			LogPrintf(ErrorLevel, "setupLoggerByDate err: %s\n", err.Error())
		}
		flog = GetFileLoggerMgr().getFileLogger(fileName)
	}

	return flog
}

func GetQnFileLogger() Logger {
	return GetFileLoggerMgr().getQnFileLogger()
}

