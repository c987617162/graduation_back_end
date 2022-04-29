package system

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/common/request"
	system "server/model/system"
	"server/utils"
)

type UserService struct{}

// Register 用户注册
func (service *UserService) Register(u system.SysUser) (err error, userInter system.SysUser) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名已注册"), userInter
	}

	u.Password = utils.MD5V([]byte(u.Password))
	u.UUID = uuid.NewV4()
	err = global.GVA_DB.Create(&u).Error
	return err, u
}

// Login 用户登录
func (service *UserService) Login(u *system.SysUser) (err error, userInter *system.SysUser) {
	if nil == global.GVA_DB {
		return fmt.Errorf("db not initialized"), nil
	}
	var user system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, u.Password).First(&user).Error
	return err, &user
}

// ChangePassword 修改密码
func (service *UserService) ChangePassword(u *system.SysUser, newPassword string) (err error, userInter *system.SysUser) {
	var user system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	newPasswordMd5 := utils.MD5V([]byte(newPassword))
	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, u.Password).First(&user).Update("password", newPasswordMd5).Error
	return err, u
}

// GetUserInfoList 分页查询用户
func (service *UserService) GetUserInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.Page
	offset := info.PageSize * (info.Page - 1)
	var userList []system.SysUser
	db := global.GVA_DB.Model(&system.SysUser{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").First(&userList).Error
	return err, userList, total
}

// SetUserAuthority  设置角色（权限）
func (service *UserService) SetUserAuthority(id uint, uuid uuid.UUID, authorityId string) (err error) {
	assignErr := global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ? ", id, authorityId).First(&system.SysUseAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = global.GVA_DB.Where("uuid = ?", uuid).First(&system.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

// DeleteUser 删除用户
func (service *UserService) DeleteUser(id int) (err error) {
	var user system.SysUser
	err = global.GVA_DB.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Delete([]system.SysUseAuthority{}, "sys_user_id", id).Error
	return err
}

// SetUserInfo 设置用户信息
func (service *UserService) SetUserInfo(req system.SysUser) error {
	return global.GVA_DB.Updates(&req).Error
}

// GetUserInfo 获取用户信息
func (service *UserService) GetUserInfo(uuid uuid.UUID) (err error, user system.SysUser) {
	var reqUser system.SysUser
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return err, reqUser
	}
	return err, reqUser
}

// FindById  通过id 查找用户
func (service *UserService) FindById(id int) (err error, user *system.SysUser) {
	var u system.SysUser
	err = global.GVA_DB.Where("id = ?", id).First(&u).Error
	return err, &u
}

// FindUserByUuid 通过uuid获取用户信息
func (service *UserService) FindUserByUuid(uuid string) (err error, user *system.SysUser) {
	var u system.SysUser
	if err = global.GVA_DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}

// ResetPassword 重置
func (service *UserService) ResetPassword(ID uint) (err error) {
	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.MD5V([]byte("123456"))).Error
	return err
}
