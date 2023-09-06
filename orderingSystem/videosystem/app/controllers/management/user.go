package management

import (
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 主要包含每个请求中的调用逻辑，是gin中的handler
func Test(c *gin.Context) {
	var form request.CreateRole
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		fmt.Println(err)
	}

	if err, role := services.UserServices.CreateRole(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, role)
	}
}

// 登录
// 这里登录只返回生成的JWT信息，然后根据jwt信息获取用户详情
func Login(c *gin.Context) {
	var form request.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err, user := services.UserServices.Login(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		// 创建jwt
		tokenData, err, _ := services.JwtService.CreateToken(services.AppFuardName, user)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		response.Success(c, tokenData)
	}
}

// 登出接口
func LoginOut(c *gin.Context) {
	fmt.Println(c.Keys["token"])
	fmt.Println("asdfasdfasdf")
	err := services.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.BusinessFail(c, "登出失败")
		return
	}
	response.Success(c, nil)
}

func CommonRegister(c *gin.Context) {
	var form request.CommonRegister
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if user, err := services.UserServices.CreateCommonuser(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

// 创建用户
// 根据传参的不同，选择性创建manager 和 shopkeeper
func CreateManageuser(c *gin.Context) {
	var form request.Resister
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err, user := services.UserServices.CreateManageuser(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

func GetUserInfo(c *gin.Context) {
	userid := c.DefaultQuery("id", "0")
	if err, users := services.UserServices.GetUserInfoID(userid); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, users)
	}
}

func GetUserInfoByJwt(c *gin.Context) {
	var form request.GetUserInfo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err, user := services.UserServices.GetUserInfo(&form); err != nil {
		response.TokenFail(c, err.Error())
		return
	} else {
		var roles []string
		type responseuser struct {
			Id    uint
			Name  string
			Roles []string
		}
		for _, item := range user.Roles {
			roles = append(roles, item.Name)
		}
		response.Success(c, responseuser{user.ID.ID, user.Name, roles})
	}
}

// 更改管理员密码
func ChangePwd(c *gin.Context) {
	var form request.ChangePwd
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err := services.UserServices.ChangePwd(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, "操作成功")
	}
	return
}

// 删除管理员 并清除关联关系
func DeleteUser(c *gin.Context) {
	var form request.Deleteuser
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err := services.UserServices.DeleteUser(&form); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		response.Success(c, "操作成功")
	}
	return
}

func CreateRole(c *gin.Context) {
	var form request.CreateRole
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err, role := services.UserServices.CreateRole(&form); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		response.Success(c, role)
	}
	return
}

func CreatePermission(c *gin.Context) {
	var form request.CreatePermission
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err, role := services.UserServices.CreatePermission(&form); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		response.Success(c, role)
	}
	return
}

func EditUserPermission(c *gin.Context) {
	// 不同请求方式对应不同的处理方式
	var form request.EditUserPermission
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if err, role := services.UserServices.EditUserPermission(&form); err != nil {
			response.BusinessFail(c, err.Error())
			return
		} else {
			response.Success(c, role)
		}
		return
	}
}
