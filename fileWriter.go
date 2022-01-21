package goLog

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrLogFileWriterEmpty = errors.New("logFileWriter is nil")
	ErrFileNotOpen = errors.New("file not opened")
)

func (p *LogFileWriter) Write(data []byte) (int, error) {
	if p == nil {
		return 0, ErrLogFileWriterEmpty
	}
	if p.file == nil {
		return 0, ErrFileNotOpen
	}
	n, err := p.file.Write(data)
	p.size += int64(n)
	// 文件最大 64 K byte
	if p.size > p.maxSize {
		p.file.Close()
		sysLog.Errorf("log file full")
		count, err := countDirFileNum(p.dirPath)
		if err != nil {
			sysLog.Warnf("CountDirFileNum err: %s\n", err.Error())
			return n, err
		}
		fullPath := p.dirPath + "/" + p.fileName
		renamePath := fullPath + "." + strconv.Itoa(count)
		if err := os.Rename(fullPath, renamePath); err != nil {
			// 出错则不切割
			sysLog.Errorf("Rename file: %s, to %s, err: %s\n", fullPath, renamePath, err.Error())
			return n, err
		}
		p.file, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			sysLog.Errorf("OpenFile a new file err: %s\n", err.Error())
			return n, err
		}
		p.size = 0
	}
	return n, err
}
