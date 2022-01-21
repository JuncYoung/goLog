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
	logFullPath := p.dirPath + "/" + p.fileName + "." + nowDateStr
	p.file, err = os.OpenFile(logFullPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		return fmt.Errorf("openFile a new file err: %s\n", err.Error())
	}

	p.dateStr = nowDateStr
	return nil
}

func (p *LogFileWriterByDate) doMill() {
	p.startMill.Do(func() {
		p.millCh = make(chan bool, 1)
		go p.runMill()
	})
	select {
	case p.millCh <- true:
	case <-time.After(1 * time.Second):
		LogPrintf(ErrorLevel, "doMill stuck\n")
	default:
	}
}

func (p *LogFileWriterByDate) runMill() {
	for range p.millCh {
		err := p.millRunOnce()
		if err != nil {
			LogPrintf(ErrorLevel, "millRunOnce err: %s\n", err.Error())
		}
	}
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

func (p *LogFileWriterByDate) millRunOnce() error {
	if p.rotateMaxAge <= 0 {
		return nil // no need to clean old files
	}

	diff := time.Duration(int64(24*time.Hour) * int64(p.rotateMaxAge))
	cutoff := time.Now().Add(-1 * diff)
	oldFiles, err := p.getOldLogFiles(cutoff.Unix())
	if err != nil {
		return err
	}

	for _, file := range oldFiles {
		if err := os.Remove(file); err != nil {
			LogPrintf(WarnLevel, "remove old log file: %s err: %s\n", file, err.Error())
		}
	}

	LogPrintf(DebugLevel, "success do millRunOnce, files: %+v, dirPath: %s\n", oldFiles, p.dirPath)
	return nil
}

func (p *LogFileWriterByDate) Write(data []byte) (int, error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	nowDateStr := time.Now().Format("20060102")
	p.mutex.Lock()
	if p.dateStr != nowDateStr {
		if err := p.initLogFile(nowDateStr); err != nil {
			LogPrintf(ErrorLevel, "initLogFile err: %s\n", err.Error())
		}
		p.doMill()
	}
	p.mutex.Unlock()
	n, err := p.file.Write(data)
	return n, err
}

func (p *LogFileWriterByDate) mustDoMillFirst() {
	p.once.Do(func() {
		p.doMill()
	})
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

