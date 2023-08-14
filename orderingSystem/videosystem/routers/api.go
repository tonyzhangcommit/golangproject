package routers

import (
	"orderingsystem/app/controllers/management"

	"github.com/gin-gonic/gin"
)

func SetApiGroupRouters(router *gin.RouterGroup) {
	// 公用接口
	// 登录，登出，获取用户信息（用户名，电话，注册时间）
	router.POST("/login", management.Login)
	router.POST("/logout", management.LoginOut)
	router.POST("/getuserinfo", management.GetUserInfo)
	//  小程序相关接口，中间件对访问次数做限制
	commonUser := router.Group("commonuser/")
	{
		commonUser.POST("/getorderlist")
	}

	// 管理员接口
	manageApi := router.Group("/manage")
	{
		manageApi.POST("/getuserinfo", management.GetUserInfo)
	}

	// 超管接口
	// manageApi := router.Group("manage/").Use(middleware.JWTAuth("manager"))
	superAdminApi := router.Group("superadmin/")
	{
		superAdminApi.POST("/getuserinfo", management.GetUserInfo)
		superAdminApi.POST("/createuser", management.Createuser)
		superAdminApi.POST("/edituser", management.Edituser)
		superAdminApi.GET("/getuserlist") // 代理列表
		superAdminApi.POST("/editrole", management.CreateRole)
		superAdminApi.POST("/editper", management.CreatePermission)
		superAdminApi.Handle("GET", "/editroleper", management.EditRolePermission)
		superAdminApi.Handle("POST", "/editroleper", management.EditRolePermission)
		superAdminApi.Handle("DELETE", "/editroleper", management.EditRolePermission)
	}
}
