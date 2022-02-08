package goLog

func InitConf(fileLog SettingDetail, qnLog SettingDetail) {
	sets = settings{
		FileLog: fileLog,
		QnLog:   qnLog,
	}
}
