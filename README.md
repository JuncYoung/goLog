#goLog

goLog.InitConf("E:\\qnzs\\aiqc\\src\\goLog\\logeeee", 0, 0, false)
goLog.SetupQnFormatByDate("E:\\qnzs\\aiqc\\src\\goLog\\logsss", "qnApi.log", "demoService", 0, 0, true, goLog.DebugLevel)

goLog.SetSysLevel(goLog.DebugLevel)
goLog.LogPrintfWithID(goLog.DebugLevel, "xxx", "xxx %s", "ddd")
goLog.LogPrintf(goLog.DebugLevel,"xxx %s", "ddd")

goLog.GetQnFileLogger().ParseQnApiLogFormat(true, 0, []string{"12345", "aaaaa", "-----"})
goLog.GetQnFileLogger().ParseQnApiLogFormat(false, 100, "success")

goLog.GetFileLogger("demo.log").Errorf("this is %s", "demo log")
goLog.GetFileLogger("success.log").Errorf("this is %s", "success log")
