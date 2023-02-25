package goLog

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyLog = errors.New("log pointer empty")
)

func GetFileLogger(fileName string) Logger {
	flSets := sets.FileLog
	if fileName == "" {
		LogPrintf(ErrorLevel, "get empty file name")
		return Logger{}
	}

	flog := GetFileLoggerMgr().getFileLogger(fileName)
	if flog.Logger == nil {
		if err := SetupLoggerByDate(flSets.LogDir, fileName, flSets.RotateMaxAge, flSets.Skip, flSets.Report, DebugLevel, flSets.TimeFormat); err != nil {
			LogPrintf(ErrorLevel, "setupLoggerByDate err: %s\n", err.Error())
		}
		flog = GetFileLoggerMgr().getFileLogger(fileName)
	}

	return flog
}

func GetQnFileLogger() Logger {
	return GetFileLoggerMgr().getQnFileLogger()
}
