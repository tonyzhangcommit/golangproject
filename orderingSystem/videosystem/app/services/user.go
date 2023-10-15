package services

import (
	"errors"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"orderingsystem/utils"
	"reflect"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	if !user.Status {
		err = errors.New("用户已封禁")
		return
	}
	return
}

// if userName != "" {
// 	query = query.Where("name = ?", userName)
// }
// if veriCode != "" {
// 	query = query.Where("identificationcode = ?", veriCode)
// }

func getAllSubUsers(userID uint) []models.User {
	var users []models.User
	var user models.User
	global.App.DB.Preload("UserList").Preload("Roles").First(&user, userID)
	users = append(users, user)

	for _, i := range user.UserList {
		subUsers := getAllSubUsers(i.ID.ID)
		users = append(users, subUsers...)
	}
	return users
}

type UserInfo struct {
	ID                 uint
	Name               string
	Telnumber          string
	Status             bool
	Roles              string
	Permissions        string
	IdentificationCode string
	CountSelfUser      int
	ManagerID          uint
	CreateTime         string
}

func filterByName(userlist []UserInfo, username string) (ol []UserInfo) {
	for _, item := range userlist {
		if item.Name == username {
			ol = append(ol, item)
		}
	}
	return
}

func filterByCode(userlist []UserInfo, code string) (ol []UserInfo) {
	for _, item := range userlist {
		if item.IdentificationCode == code {
			ol = append(ol, item)
		}
	}
	return
}

func filterByTel(userlist []UserInfo, tel string) (ol []UserInfo) {
	for _, item := range userlist {
		if item.Telnumber == tel {
			ol = append(ol, item)
		}
	}
	return
}

func GetPageData(data interface{}, pageCount, pageIndex int) interface{} {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic("data must be a slice")
	}

	// 确定总数据量
	total := value.Len()

	// 计算总页数
	totalPages := total / pageCount
	if total%pageCount != 0 {
		totalPages++
	}

	// 确定要获取的页数
	if pageIndex <= 0 {
		pageIndex = 1
	} else if pageIndex > totalPages {
		pageIndex = totalPages
	}

	// 计算起始索引和结束索引
	startIdx := (pageIndex - 1) * pageCount
	endIdx := startIdx + pageCount
	if endIdx > total {
		endIdx = total
	}

	// 获取指定页的数据
	result := reflect.MakeSlice(value.Type(), endIdx-startIdx, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		result.Index(i - startIdx).Set(value.Index(i))
	}

	return result.Interface()
}

func (userServices userServices) GetUserInfoID(c *gin.Context) (err error, userlist []UserInfo, totalcount int) {
	var user models.User
	userID := c.DefaultQuery("id", "")
	var userId int
	if userId, err = strconv.Atoi(userID); err != nil {
		err = errors.New("参数错误")
		return
	}
	if err = global.App.DB.First(&user, userId).Error; err != nil {
		err = errors.New("用户不存在")
		return
	} else {
		userName := c.DefaultQuery("name", "")
		veriCode := c.DefaultQuery("veriCode", "")
		telNum := c.DefaultQuery("telNum", "")
		page := c.DefaultQuery("page", "")
		page_size := c.DefaultQuery("page_size", "")
		var pageN int
		var page_sizeN int
		if pageN, err = strconv.Atoi(page); err != nil {
			err = errors.New("参数错误")
			return
		}
		if page_sizeN, err = strconv.Atoi(page_size); err != nil {
			err = errors.New("参数错误")
			return
		}
		users := getAllSubUsers(user.ID.ID)
		for _, item := range users {
			var count int
			if item.UserList != nil {
				count = len(item.UserList)
			} else {
				count = 0
			}
			var role string
			if len(item.Roles) == 0 {
				role = ""
			} else {
				role = item.Roles[0].Name
			}
			userlist = append(userlist, UserInfo{
				ID:                 item.ID.ID,
				Name:               item.Name,
				Telnumber:          item.Telnumber,
				Status:             item.Status,
				Roles:              role,
				IdentificationCode: item.IdentificationCode,
				CountSelfUser:      count,
				ManagerID:          item.ManagerID,
				CreateTime:         item.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		if userName != "" {
			userlist = filterByName(userlist, userName)
		}
		if veriCode != "" {
			userlist = filterByCode(userlist, veriCode)
		}
		if telNum != "" {
			userlist = filterByTel(userlist, telNum)
		}
		totalcount = len(userlist)
		result := GetPageData(userlist, page_sizeN, pageN)
		userlist = result.([]UserInfo)
	}
	return
}

// 管理员获取当前下所有用户
func (userServices userServices) GetManagerUsers(c *gin.Context, userId string) (err error, users []models.User) {
	var manager models.User
	userTel, isExit := c.Get("userTel")
	if err = global.App.DB.Where("telnumber=?", userTel).First(&manager).Error; err != nil && !isExit {
		err = errors.New("非法Token")
		return
	}
	if userId == "0" {
		// 获取当前管理员下所有的用户信息
		if err = global.App.DB.Where("manager_id = ?", manager.ID.ID).Preload("UserList").Preload("Roles").Find(&users).Error; err != nil {
			err = errors.New("查询失败")
		}
	} else {
		if userintId, err := strconv.Atoi(userId); err != nil {
			err = errors.New("参数错误")
		} else {
			parameter := []int{userintId}
			if err = global.App.DB.Where("manager_id = ?", manager.ID.ID).Preload("UserList").Preload("Roles").Preload("Permissions").Find(&users, parameter).Error; err != nil {
				err = errors.New("查询失败")
			}
		}
	}
	return
}

func (userServices userServices) GetRoles(userId string) (err error, roles []models.Role) {
	var userid int
	var user models.User
	if userid, err = strconv.Atoi(userId); err != nil || userId == "" {
		err = errors.New("参数错误")
	} else {
		if err = global.App.DB.Preload("Roles").First(&user, userid).Error; err != nil {
			err = errors.New("用户不存在")
		} else {
			global.App.DB.Find(&roles)
			isSuperM := false
			for _, item := range user.Roles {
				if item.Name == "superadmin" {
					isSuperM = true
					break
				}
			}
			if !isSuperM {
				var index int
				for i, item := range roles {
					if item.Name == "superadmin" {
						index = i
						break
					}
				}
				roles = append(roles[:index], roles[index+1:]...)
			}
		}
	}
	return
}

// 管理员对用户开启代理权限
func (userServices userServices) ChangeProxy(c *gin.Context, params *request.EditProxy) (user models.User, err error) {
	var manager models.User
	userTel, isExit := c.Get("userTel")
	if err = global.App.DB.Where("telnumber=?", userTel).First(&manager).Error; err != nil && !isExit {
		err = errors.New("非法Token")
		return
	}
	if err = global.App.DB.Where("id=?", params.UserId).First(&user).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	// 更改用户角色
	if params.Option == "add" {
		global.App.DB.Model(&user).Association("Roles").Clear()
		role := models.Role{
			Name: "manager",
		}
		global.App.DB.Model(&user).Association("Roles").Append(&role)
	} else if params.Option == "delete" {
		user.Status = false
		global.App.DB.Model(&user).Association("Roles").Clear()
		global.App.DB.Save(&user)
	} else {
		err = errors.New("参数错误")
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
func (userServices userServices) CreateManageuser(params *request.Resister) (err error, user models.User) {
	// 首先判断用户角色
	errT := global.App.DB.First(&user, "telnumber = ?", params.Mobile).Error
	if errT != nil {
		user = models.User{
			Name:      params.Name,
			Telnumber: params.Mobile,
			Password:  utils.BcryptMake([]byte(params.Password)),
		}
		// 获取角色列表
		var roles []models.Role
		var role models.Role
		global.App.DB.Find(&roles)
		isrightrole := false
		for _, value := range roles {
			if value.Name == params.Role {
				isrightrole = true
				role = value
				break
			}
		}

		if !isrightrole {
			err = errors.New("角色不存在")
			return
		}
		// 推荐码
		UniqueCode := utils.GenerateRandomString(6)
		user.IdentificationCode = UniqueCode
		var tempUser models.User
		if err = global.App.DB.First(&tempUser, params.ManagerID).Error; err != nil {
			err = errors.New("用户不存在")
			return
		}
		if params.ManagerID == 0 {
			user.ManagerID = 0
			global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
			if err = global.App.DB.Create(&user).Error; err != nil {
				err = errors.New("创建失败")
				return
			}
			global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

		} else {
			user.ManagerID = params.ManagerID
			err = global.App.DB.Create(&user).Error
		}
		err = global.App.DB.Model(&user).Association("Roles").Append(&role)
	} else {
		err = errors.New("用户已存在")
	}
	return
}

func (userServices userServices) EditUserStatus(param *request.EditUser) (err error) {
	if param.Id == param.TargetId {
		err = errors.New("不能操作自身")
	} else {
		var user models.User
		var tuser models.User
		if err = global.App.DB.First(&user, param.Id).Error; err != nil {
			err = errors.New("参数错误")
			return
		}
		if err = global.App.DB.First(&tuser, param.TargetId).Error; err != nil {
			err = errors.New("用户不存在")
			return
		}
		var status bool

		if param.Status == 0 {
			status = false
		} else if param.Status == 1 {
			status = true
		} else {
			err = errors.New("参数错误")
			return
		}
		tuser.Status = status
		err = global.App.DB.Save(&tuser).Error
	}
	return
}

// 普通用户注册
func (userServices userServices) CreateCommonuser(params *request.CommonRegister) (user models.User, err error) {
	var manager models.User
	if err = global.App.DB.Where("identificationcode =?", params.IdentificationCode).First(&manager).Error; err != nil {
		err = errors.New("推广码不存在")
	} else {
		var role models.Role
		if err = global.App.DB.Where("name = ?", "normalusers").First(&role).Error; err != nil {
			return
		}
		user.Telnumber = params.Mobile
		user.Password = utils.BcryptMake([]byte(params.Password))
		user.ManagerID = manager.ID.ID
		user.Name = params.Name
		if params.Name == "" {
			user.Name = "normal" + utils.GenerateRandomIntString(4)
		}

		if err = global.App.DB.Create(&user).Error; err != nil {
			err = errors.New("创建失败")
		} else {
			err = global.App.DB.Model(&user).Association("Roles").Append(&role)
		}
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

// 创建删除角色
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
	rItem := global.App.DB.Preload("Permissions").Preload("Roles").First(&user, userID)
	if pItem.Error != nil || rItem.Error != nil {
		err = errors.New("用户/权限不存在")
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
				return
			}
			err = global.App.DB.Model(&user).Association("Permissions").Delete(&permission)
		} else {
			err = errors.New("请求参数错误！")
		}
	}
	return
}
