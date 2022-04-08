package init

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"server/global"
	SysR "server/router/system"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	systemRouter := SysR.UserRouter{}

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Info("register swagger handler")

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	{
		systemRouter.InitUserRouter(PublicGroup)
	}

	global.GVA_LOG.Info("router register success")
	return Router
}
