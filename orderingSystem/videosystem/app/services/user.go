package services

import (
	"errors"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"orderingsystem/utils"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

type userServices struct {
}

var UserServices = new(userServices)

//  用户注册逻辑，涉及 DB 相关操作

func (userServices userServices) Test(params *request.CreateRole) (err error, role models.Role) {
	err = global.App.DB.First(&models.Role{}, "name = ?", params.Name).Error
	if err != nil {
		role = models.Role{Name: params.Name}
		err = global.App.DB.Create(&role).Error
	} else {
		err = errors.New("角色名已存在")
	}
	return
}

// 登录
func (userServices userServices) Login(params *request.Login) (err error, user models.User) {
	err = global.App.DB.Model(&user).Where("telnumber = ?", params.Mobile).Preload("Roles").First(&user).Error
	if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
		err = errors.New("用户不存在/密码错误")
	}
	return
}

// 获取用户信息
// 用户名， 角色
func (userServices userServices) GetUserInfo(params *request.GetUserInfo) (err error, user models.User) {
	tokenStr := params.Jwt
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.App.Config.Jwt.Secret), nil
	})
	if err != nil {
		err = errors.New("授权失败")
		return
	}
	claims := token.Claims.(*CustomClaims)
	userId := claims.Id
	if err = global.App.DB.Preload("Roles").First(&user, userId).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	return
}

func (userServices userServices) GetUserInfoID(userId string) (err error, users []models.User) {
	if userId == "0" {
		if err = global.App.DB.Preload("UserList").Preload("Roles").Preload("Permissions").Find(&users).Error; err != nil {
			err = errors.New("查询失败")
		}
	} else {
		if userintId, err := strconv.Atoi(userId); err != nil {
			err = errors.New("参数错误")
		} else {
			parameter := []int{userintId}
			if err = global.App.DB.Preload("UserList").Preload("Roles").Preload("Permissions").Find(&users, parameter).Error; err != nil {
				err = errors.New("查询失败")
			}
		}
	}
	return
}

func (userServices userServices) GetUserInfoById(userId int64) (err error, user models.User) {
	if err = global.App.DB.Preload("Roles").First(&user, userId).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	return
}

// 创建管理员 CreateUser
func (userServices userServices) CreateUser(params *request.Resister) (err error, user models.User) {
	// 首先判断用户角色
	errT := global.App.DB.First(&user, "telnumber = ?", params.Mobile).Error
	errN := global.App.DB.First(&user, "name = ?", params.Name).Error
	if errT != nil && errN != nil {
		user = models.User{
			Name:      params.Name,
			Telnumber: params.Mobile,
			Password:  utils.BcryptMake([]byte(params.Password)),
		}
		// 获取角色列表
		var roles []models.Role
		global.App.DB.Find(&roles)
		isrightrole := false
		for _, value := range roles {
			if value.Name == params.Role {
				isrightrole = true
				break
			}
		}
		if !isrightrole {
			err = errors.New("角色不存在")
			return
		}
		UniqueCode := utils.GenerateRandomString(12)
		user.IdentificationCode = UniqueCode

		user.ManagerID = params.ManagerID
		err = global.App.DB.Create(&user).Error
		role := models.Role{
			Name: params.Role,
		}
		err = global.App.DB.Model(&user).Association("Roles").Append(&role)
		// 增加关联
		global.App.DB.Model(&user).Association("UserList").Append(&user)

	} else {
		err = errors.New("用户已存在")
	}
	return
}

// 重置密码
func (userServices userServices) ChangePwd(params *request.ChangePwd) (err error) {
	var user models.User
	if err = global.App.DB.Where("telnumber= ?", params.Mobile).First(&user).Error; err != nil {
		err = errors.New("用户不存在")
	} else {
		user.Password = utils.BcryptMake([]byte(params.Password))
		global.App.DB.Save(&user)
	}
	return
}

// 删除用户
func (userServices userServices) DeleteUser(params *request.Deleteuser) (err error) {
	var user models.User
	if err = global.App.DB.Model(&user).Where("telnumber= ?", params.Mobile).First(&user).Error; err != nil {
		err = errors.New("用户不存在")
		return
	} else {
		// 删除流程，用户关联 权限，优惠券，店铺，订单
		if err = global.App.DB.Select("Roles", "Coupons", "Shops", "Orders").Delete(&user).Error; err != nil {
			err = errors.New("删除失败")
		}
		err = errors.New("删除成功")
		return
	}
}

// 创建角色
func (userServices userServices) CreateRole(params *request.CreateRole) (err error, role models.Role) {
	err = global.App.DB.First(&models.Role{}, "name = ?", params.Name).Error
	if params.Option == "create" {
		if err != nil {
			role = models.Role{Name: params.Name}
			err = global.App.DB.Create(&role).Error
		} else {
			err = errors.New("角色名已存在")
		}
	} else if params.Option == "delete" {
		if err == nil {
			err = global.App.DB.Where("name", params.Name).Delete(&role).Error
			if err != nil {
				err = errors.New("删除失败")
			} else {
				err = errors.New("删除成功")
			}
		} else {
			err = errors.New("角色不存在")
		}
	} else {
		err = errors.New("option参数错误")
	}
	return
}

// 创建权限
func (userServices userServices) CreatePermission(params *request.CreatePermission) (err error, per models.Permission) {
	err = global.App.DB.First(&models.Permission{}, "name = ?", params.Name).Error
	if params.Option == "create" {
		if err != nil {
			per = models.Permission{Name: params.Name}
			err = global.App.DB.Create(&per).Error
		} else {
			err = errors.New("权限已存在")
		}
	} else if params.Option == "delete" {
		if err == nil {
			err = global.App.DB.Where("name", params.Name).Delete(&per).Error
			if err != nil {
				err = errors.New("删除失败")
			} else {
				err = errors.New("删除成功")
			}
		} else {
			err = errors.New("权限不存在")
		}
	} else {
		err = errors.New("option参数错误")
	}
	return
}

// 编辑角色权限信息
func (userServices userServices) EditUserPermission(params *request.EditUserPermission) (err error, user models.User) {
	userID := params.UserID
	pname := params.Permissionname
	var permission models.Permission
	pItem := global.App.DB.First(&permission, "name = ?", pname)
	rItem := global.App.DB.First(&user, userID)
	if pItem.Error != nil || rItem.Error != nil {
		err = errors.New("用户/权限不存在")
	} else {
		if err = global.App.DB.Preload("Permissions").First(&user).Error; err != nil {
			return
		} else {
			isHasP := false
			for _, item := range user.Permissions {
				if item.Name == pname {
					isHasP = true
					break
				}
			}
			if params.Option == "create" {
				if isHasP {
					err = errors.New("权限已存在")
					return
				}
				err = global.App.DB.Model(&user).Association("Permissions").Append(&permission)
			} else if params.Option == "delete" {
				if !isHasP {
					err = errors.New("权限不存在")
					return
				}
				err = global.App.DB.Model(&user).Association("Permissions").Delete(&permission)
			} else {
				err = errors.New("请求参数错误！")
			}
			return
		}
	}
	return
}