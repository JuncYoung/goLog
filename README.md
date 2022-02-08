#goLog

##useage:
    nowTime := time.Now()
    goLog.SetSysLevel(goLog.DebugLevel)
    goLog.InitConf(goLog.SettingDetail{
        LogDir:       "/tmp/myLog",
        RotateMaxAge: 0,
        Skip:         0,
        Report:       false,
        Level:        goLog.DebugLevel,
    }, goLog.SettingDetail{
        LogDir:       "/tmp/qnLog",
        RotateMaxAge: 0,
        Skip:         0,
        Report:       true,
        Level:        goLog.DebugLevel,
    })
    goLog.SetupQnFormatByDate("qnApi.log", "demoService")
    
    goLog.LogPrintfWithID(goLog.DebugLevel, "xxx", "xxx %s", "ddd")
    goLog.LogPrintf(goLog.DebugLevel, "xxx %s", "ddd")
    
    goLog.GetQnFileLogger().ParseQnApiLogInput([]string{"12345", "aaaaa", "-----"})
    goLog.GetQnFileLogger().ParseQnApiLogOutput(nowTime, "success")
    
    goLog.GetFileLogger("demo.log").Warnf("this is %s", "demo warning log")
    goLog.GetFileLogger("demo.log").Errorf("this is %s", "demo err log")
    goLog.GetFileLogger("success.log").Infof("this is %s", "success log")
