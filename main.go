package main

import (
	"go.uber.org/zap"
	"server/global"
	"server/init"
)

func main() {
	global.GVA_VP = init.Viper()
	global.GVA_LOG = init.Zap()
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = init.GormMysql()
	init.Timer()

}
