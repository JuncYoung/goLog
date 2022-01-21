package goLog

import (
	"sync"
)

const (
	qnLog     = "qnApi"
)

var (
	sets          settings
	fileLoggerMgr *FileLoggerMgr
	once          sync.Once
)

func NewLogFileMgr() *FileLoggerMgr {
	fm := make(map[string]Logger)
	return &FileLoggerMgr{
		fileWriters: fm,
	}
}

func GetFileLoggerMgr() *FileLoggerMgr {
	once.Do(func() {
		fileLoggerMgr = NewLogFileMgr()
	})

	return fileLoggerMgr
}

func (f *FileLoggerMgr) AddOne(key string, value Logger) {
	if f == nil {
		LogPrintf(ErrorLevel, "LogFileMgr empty\n")
		return
	}

	if f.fileWriters == nil {
		f.fileWriters = make(map[string]Logger)
	}

	f.fileWriters[key] = value
	return
}

func (f *FileLoggerMgr) getQnFileLogger() Logger {
	return f.fileWriters[qnLog]
}

func (f *FileLoggerMgr) getFileLogger(key string) Logger {
	return f.fileWriters[key]
}