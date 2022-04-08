package main

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/init"
	"time"
)

func main() {
	global.GVA_VP = init.Viper()
	global.GVA_LOG = init.Zap()
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = init.GormMysql()
	init.Timer()
	if global.GVA_DB != nil {
		init.RegisterTables(global.GVA_DB)
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
	Router := init.Routers()
	port := "8080"
	s := init.InitServer(port, Router)
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", port))
	fmt.Printf(`
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
`, port)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
