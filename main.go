package main

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/initialize"
	"time"
)

func main() {
	global.GVA_VP = initialize.Viper()
	global.GVA_LOG = initialize.Zap()
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.GormMysql()
	//initialize.Timer()
	if global.GVA_DB != nil {
		initialize.RegisterTables(global.GVA_DB)
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	// 预览Redis 相关配置

	// 读取JWT
	//if global.GVA_DB != nil {
	//	system.LoadAll()
	//}

	// 初始化路由

	// 启动服务
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initialize.InitServer(address, Router)
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))
	fmt.Printf(`
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:%s
`, address, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
