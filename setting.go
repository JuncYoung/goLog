package goLog

func InitConf(dir string, rotateMaxAge int, skip int, report bool) {
	sets.Com = settingDetail{
		LogDir: dir,
		RotateMaxAge:  rotateMaxAge,
		Skip:  skip,
		Report:  report,
	}
}
