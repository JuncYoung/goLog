package goLog

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 10
)

var (
	callerInitOnce     sync.Once
	minimumCallerDepth int
	logrusPackage      string
)

func (p *LogFileWriterByDate) initLogFile(nowDateStr string) error {
	var err error
	logFullPath := p.dirPath + "/" + p.fileName + "-" + nowDateStr + ".log"
	p.file, err = os.OpenFile(logFullPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("openFile a new file err: %s\n", err.Error())
	}

	p.dateStr = nowDateStr
	return nil
}

func (p *LogFileWriterByDate) getOldLogFiles(timestampLess int64) ([]string, error) {
	result := make([]string, 0, 100)
	folderPath := p.dirPath
	infos, readDirError := ioutil.ReadDir(folderPath)
	if readDirError != nil {
		return nil, readDirError
	}

	for _, info := range infos {
		if !info.IsDir() && timestampLess > info.ModTime().Unix() {
			result = append(result, folderPath+"/"+info.Name())
		}
	}

	return result, nil
}

func (p *LogFileWriterByDate) Write(data []byte) (int, error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	nowDateStr := time.Now().Format(p.TimeFormat)
	p.mutex.Lock()
	if p.dateStr != nowDateStr {
		if err := p.file.Close(); err != nil {
			LogPrintf(ErrorLevel, "pfile close err: %s\n", err.Error())
		}
		if err := p.initLogFile(nowDateStr); err != nil {
			LogPrintf(ErrorLevel, "initLogFile err: %s\n", err.Error())
		}
	}
	p.mutex.Unlock()
	n, err := p.file.Write(data)
	return n, err
}

func (p *LogFileWriterByDate) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05.999")
	if !strings.HasSuffix(entry.Message, "\n") {
		entry.Message = entry.Message + "\n"
	}
	if entry.Caller == nil {
		msg := fmt.Sprintf("%s [%s]\t%s", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
		return []byte(msg), nil
	}

	file, line, fName := entry.Caller.File, entry.Caller.Line, entry.Caller.Function

	msg := fmt.Sprintf("%s [%s] [%s:%d %s]\t%s", timestamp, strings.ToUpper(entry.Level.String()), file, line, fName, entry.Message)
	return []byte(msg), nil
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller(skip int) *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)
		callDepth := -1
		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "LogxPrintf") && callDepth == -1 {
				callDepth = i
			}
			logrusPackage = getPackageName(funcName)
			if i == callDepth+skip { // 允许在找到顶层方法后再往栈顶追溯一定层级
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
