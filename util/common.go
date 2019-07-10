package util

import "github.com/astaxie/beego/logs"

func GetLogs()*logs.BeeLogger{
	return logs.GetBeeLogger()
}
