package main

import "github.com/JuncYoung/goLog"

// example
func main() {
	goLog.SetSysLevel(0)
	goLog.InitConf("E:\\qnzs\\goPro\\goLog\\logeeee", 7, 0, false)
	goLog.SetupQnFormatByDate("E:\\qnzs\\goPro\\goLog\\logsss", "qnApi.log", "demoService", 7, 0, true, goLog.DebugLevel)


	goLog.LogPrintfWithID(goLog.DebugLevel, "xxx", "ssssssssssssss %s", "ddd")
	goLog.LogPrintf(goLog.DebugLevel,"ssssssssssssss %s", "ddd")

	goLog.GetQnFileLogger().ParseQnApiLogFormat(true, 0, []string{"12345", "aaaaa", "-----"})
	goLog.GetQnFileLogger().ParseQnApiLogFormat(false, 100, "success")

	goLog.GetFileLogger("demo.log").Errorf("fffffff %s", "good")
}
