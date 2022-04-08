package system

import (
	"errors"
	"gorm.io/gorm"
	"server/global"
	"server/model/common/request"
	"server/model/system"
)

// AuthorityService 创建角色
type AuthorityService struct{}

var AuthorServiceApp = new(AuthorityService)

func (authorityService *AuthorityService) CreateAuthority(auth system.SysAuthority) (err error, authority system.SysAuthority) {
	var authorityBox system.SysAuthority
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色id"), auth
	}
	err = global.GVA_DB.Create(&auth).Error
	return err, auth
}

// UpdateAuthority 修改角色
func (authorityService *AuthorityService) UpdateAuthority(auth system.SysAuthority) (err error, authority system.SysAuthority) {
	err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Updates(&auth).Error
	return err, auth
}

// DeleteAuthority 删除角色
func (authorityService *AuthorityService) DeleteAuthority(auth system.SysAuthority) (err error) {
	if errors.Is(global.GVA_DB.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	if len(auth.Users) != 0 {
		return errors.New("该角色用户正在使用")
	}
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色用户正在使用")
	}
	db := global.GVA_DB.Preload("SysBaseMenus").Where("authority_id = ?", auth.AuthorityId).First(auth)
	err = db.Unscoped().Delete(auth).Error
	if err != nil {
		return
	}

	err = global.GVA_DB.Delete(&[]system.SysUseAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error
	err = global.GVA_DB.Delete(&[]system.SysAuthorityBtn{}, "authority_id = ?", auth.AuthorityId).Error
	return err
}

// GetAuthorityInfoList   分页查询角色
func (authorityService *AuthorityService) GetAuthorityInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysAuthority{})
	err = db.Select("count(*) as total").Pluck("total", &total).Error
	var authority []system.SysAuthority
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Find(&authority).Error

	return err, authority, total

}

// GetAuthorityInfo  获取所有角色信息
func (authorityService *AuthorityService) GetAuthorityInfo(auth system.SysAuthority) (err error, sa system.SysAuthority) {
	err = global.GVA_DB.Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(&sa).Error
	return err, sa
}
