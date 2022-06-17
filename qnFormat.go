package goLog

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *QnFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05.999")
	if !strings.HasSuffix(entry.Message, "\n") {
		entry.Message = entry.Message + "\n"
	}
	if entry.Caller == nil {
		msg := fmt.Sprintf("%s [%s]\t%s", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
		return []byte(msg), nil
	}

	file, line, fName := path.Base(entry.Caller.File), entry.Caller.Line, entry.Caller.Function
	msg := fmt.Sprintf("%s [%s] [%s:%d %s]\t%s", timestamp, strings.ToUpper(entry.Level.String()), file, line, fName, entry.Message)
	return []byte(msg), nil
}
