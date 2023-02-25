package main

import (
	"time"

	"github.com/JuncYoung/goLog"
)

// example
func main() {
	nowTime := time.Now()
	goLog.SetSysLevel(goLog.DebugLevel)
	goLog.InitConf(goLog.SettingDetail{
		LogDir:       "./myLog", // absolute path is suggested
		RotateMaxAge: 0,
		Skip:         0,
		Report:       false,
		Level:        goLog.DebugLevel,
		TimeFormat:   "2006-01-02-15-04-05",
	}, goLog.SettingDetail{
		LogDir:       "./qnLog", // absolute path is suggested
		RotateMaxAge: 0,
		Skip:         0,
		Report:       true,
		Level:        goLog.DebugLevel,
		TimeFormat:   "2006-01-02-15-04-05",
	})
	goLog.SetupQnFormatByDate("qnApi", "demoService")

	goLog.LogPrintfWithID(goLog.DebugLevel, "xxx", "xxx %s", "ddd")
	goLog.LogPrintf(goLog.DebugLevel, "xxx %s", "ddd")

	goLog.GetQnFileLogger().ParseQnApiLogInput([]string{"12345", "aaaaa", "-----"})
	goLog.GetQnFileLogger().ParseQnApiLogOutput(nowTime, "success")
	time.Sleep(10 * time.Second)
	goLog.GetQnFileLogger().WithMethod("mainFunc").WithMsg("test output").ParseQnApiLogOutput(nowTime, "mainFunc success")
	goLog.GetQnFileLogger().QnInternalError("error occur")

	goLog.GetFileLogger("demo").Warnf("this is %s", "demo warning log")
	goLog.GetFileLogger("demo").Errorf("this is %s", "demo err log")
	goLog.GetFileLogger("success").Infof("this is %s", "success log")
	goLog.GetFileLogger("other").WithFields(map[string]interface{}{"nice": "good", "number": 233.333}).Debug("fields print")
}
