package system

import (
	"github.com/gin-gonic/gin"
	SysApi "server/api/system"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	userRouter := Router.Group("user")
	var baseApi SysApi.BaseApi
	{
		userRouter.POST("login", baseApi.Login)                   //登录
		userRouter.POST("admin_register", baseApi.Register)       // 管理员注册账号
		userRouter.POST("changePassword", baseApi.ChangePassword) // 用户修改密码
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)       // 删除用户
		userRouter.PUT("setUserInfo", baseApi.SetUserInfo)        // 设置用户信息
		userRouter.PUT("setSelfInfo", baseApi.SetSelfInfo)        // 设置自身信息
		userRouter.POST("resetPassword", baseApi.ResetPassword)   // 设置用户权限组
		userRouter.POST("getUserInfo", baseApi.GetUserInfo)       // 分页获取用户列表
		userRouter.GET("getUserInfo", baseApi.GetUserInfo)        // 获取自身信息
	}
	{
		baseRouter.GET("captcha", baseApi.Captcha)
	}

}
