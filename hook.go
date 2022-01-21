package goLog

import "github.com/sirupsen/logrus"

type DateLogHook struct {
	skip int
}

func (h *DateLogHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}
}

func (h *DateLogHook) Fire(entry *logrus.Entry) error {
	entry.Caller = getCaller(h.skip)
	return nil
}
