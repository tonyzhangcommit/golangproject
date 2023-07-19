package services

import (
	"errors"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"orderingsystem/utils"

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
	err = global.App.DB.Model(&user).Where("telnumber = ?", params.Mobile).First(&user).Error
	if err != nil || utils.BcryptMakeCheck([]byte(user.Password), params.Password) {
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

// 创建用户 CreateUser
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

		err = global.App.DB.Create(&user).Error
		role := models.Role{
			Name: params.Role,
		}
		err = global.App.DB.Model(&user).Association("Roles").Append(&role)

	} else {
		err = errors.New("用户已存在")
	}
	return
}

// 删除用户
func (userServices userServices) Edituser(params *request.Deleteuser) (err error) {
	var user models.User
	if err = global.App.DB.Model(&user).Where("telnumber= ?", params.Mobile).First(&user).Error; err != nil {
		err = errors.New("用户不存在")
		return
	} else {
		// 删除流程，用户关联 权限，优惠券，店铺，订单
		if params.Option == "delete" {
			global.App.DB.Select("Roles", "Coupons", "Shops", "Orders").Delete(&user)
			if err != nil {
				err = errors.New("删除失败")
			}
		} else if params.Option == "edit" {
			err = global.App.DB.Model(&models.User{}).Where("id", user.ID.ID).Update("status", !user.Status).Error
		} else {
			err = errors.New("参数错误")
		}
		return
	}
}

// 创建角色
func (userServices userServices) CreateRole(params *request.CreateRole) (err error, role models.Role) {
	err = global.App.DB.First(&models.Role{}, "name = ?", params.Name).Error
	if err != nil {
		role = models.Role{Name: params.Name}
		err = global.App.DB.Create(&role).Error
	} else {
		err = errors.New("角色名已存在")
	}
	return
}

// 创建权限
func (userServices userServices) CreatePermission(params *request.CreatePermission) (err error, per models.Permission) {
	err = global.App.DB.First(&models.Permission{}, "name = ?", params.Name).Error
	if err != nil {
		per = models.Permission{Name: params.Name}
		err = global.App.DB.Create(&per).Error
	} else {
		err = errors.New("权限已存在")
	}
	return
}

// 编辑角色权限信息
func (userServices userServices) EditRolePermission(params *request.EditRolePermission, method string) (err error, role models.Role) {
	rname := params.Rolename
	pname := params.Permissionname

	var permission models.Permission
	var permissions []models.Permission

	pItem := global.App.DB.First(&permission, "name = ?", pname)
	rItem := global.App.DB.First(&role, "name = ?", rname)

	if pItem.Error != nil || rItem.Error != nil {
		err = errors.New("角色/权限不存在")
	} else {
		methodMap := map[string]string{
			"GET":    "",
			"POST":   "",
			"DELETE": "",
		}
		if _, ok := methodMap[method]; !ok {
			err = errors.New("请求方式有误")
		} else {
			if method == "GET" {
				err = global.App.DB.Preload("Permissions").Where("name= ?", rname).First(&role).Error
			} else if method == "POST" {
				// 首先查询当前角色是否存在指定权限
				err = global.App.DB.Model(&role).Association("Permissions").Find(&permissions)
				if err != nil {
					return
				} else {
					containPermission := false
					for _, p := range permissions {
						if p.ID == permission.ID {
							containPermission = true
							break
						}
					}
					if containPermission {
						err = errors.New("该权限已存在")
					} else {
						err = global.App.DB.Model(&role).Association("Permissions").Append(&permission)
						global.App.DB.Preload("Permissions").Where("name= ?", rname).First(&role)
					}
				}

			} else {
				// 删除权限
				err = global.App.DB.Model(&role).Association("Permissions").Find(&permissions)
				if err != nil {
					return
				} else {
					containPermission := false
					for _, p := range permissions {
						if p.ID == permission.ID {
							containPermission = true
							break
						}
					}
					if !containPermission {
						err = errors.New("权限不存在")
					} else {
						err = global.App.DB.Model(&role).Association("Permissions").Delete(&permission)
						global.App.DB.Preload("Permissions").Where("name= ?", rname).First(&role)
					}
				}
			}
			return
		}
	}
	return
}
