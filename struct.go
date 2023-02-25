package goLog

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type SettingDetail struct {
	LogDir       string `json:"logDir"`
	RotateMaxAge int    `json:"rotateMaxAge"`
	Skip         int    `json:"skip"`
	Report       bool   `json:"report"`
	Level        int    `json:"level"`
}

type settings struct {
	FileLog SettingDetail
	QnLog   SettingDetail
}

type FileLoggerMgr struct {
	fileWriters map[string]Logger // key: fileName value: LogFileWriter
}

type QnFormatter struct {
}

type standardLogFormat struct {
	Service string      `json:"service"`
	From    interface{} `json:"from,omitempty"`
	Out     interface{} `json:"out,omitempty"`
	Msg     string      `json:"msg"`
	Time    float64     `json:"time,omitempty"`
}

type Logger struct {
	*logrus.Logger
	ServiceName string
	Msg         string
}

type LogFileWriter struct {
	dirPath  string
	fileName string
	maxSize  int64
	file     *os.File
	size     int64
}

type LogFileWriterByDate struct {
	dirPath  string
	fileName string
	file     *os.File
	dateStr  string
	mutex    *sync.Mutex

	once sync.Once

	startMill sync.Once
	millCh    chan bool

	rotateMaxAge int
}
