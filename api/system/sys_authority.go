package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/common/response"
	M "server/model/system"
	systemRes "server/model/system/response"
	"server/service/system"
	"server/utils"
)

var authorityService = system.AuthorServiceApp

type AuthorityApi struct{}

func (a *AuthorityApi) CreateAuthority(c *gin.Context) {
	var authority M.SysAuthority
	_ = c.ShouldBind(&authority)
	if err := utils.Verify(authority, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	if err, authBack := authorityService.CreateAuthority(authority); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authBack}, "创建成功", c)
	}

}
