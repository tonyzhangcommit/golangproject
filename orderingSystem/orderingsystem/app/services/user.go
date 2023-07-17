package services

import (
	"errors"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
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
