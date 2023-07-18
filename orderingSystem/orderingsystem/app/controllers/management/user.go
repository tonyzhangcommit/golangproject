package management

import (
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/services"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//  主要包含每个请求中的调用逻辑，是gin中的handler

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
	}

	if err, role := services.UserServices.Login(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, role)
	}
}

// 创建用户
// 根据传参的不同，选择性创建manager 和 shopkeeper
func Createuser(c *gin.Context) {
	var form request.Resister
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		fmt.Println(err)
		return
	}

	if err,user := services.UserServices.CreateUser(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

// 删除用户并清除关联关系
func DeleteUser(c *gin.Context) {
	var form request.Deleteuser
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err := services.UserServices.Deleteuser(&form); err != nil {
		response.BusinessFail(c, err.Error())
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
	} else {
		response.Success(c, role)
	}
	return
}

func EditRolePermission(c *gin.Context) {
	// 不同请求方式对应不同的处理方式
	var form request.EditRolePermission
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		method := c.Request.Method
		if err, role := services.UserServices.EditRolePermission(&form, method); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, role)
		}
		return
	}
}
