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
		return
	}

	if err, role := services.UserServices.CreateRole(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, role)
	}
}

func CreateRole(c *gin.Context) {
	var form request.CreateRole
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		fmt.Println(err)
		return
	}
	if err, role := services.UserServices.CreateRole(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, role)
	}
}

func CreatePermission(c *gin.Context) {
	var form request.CreatePermission
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		fmt.Println(err)
		return
	}
	if err, role := services.UserServices.CreatePermission(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, role)
	}
}

func EditRolePermission(c *gin.Context) {
	// 不同请求方式对应不同的处理方式
	var form request.EditRolePermission
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
	} else {
		method := c.Request.Method
		if err, role := services.UserServices.EditRolePermission(&form, method); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, role)
		}
	}
	return
}
