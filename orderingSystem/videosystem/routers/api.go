package routers

import (
	"orderingsystem/app/controllers/management"

	"github.com/gin-gonic/gin"
)

func SetApiGroupRouters(router *gin.RouterGroup) {
	// 公用接口
	// 登录，登出，获取用户信息（用户名，电话，注册时间）
	router.POST("/login", management.Login)
	router.POST("/register", management.CommonRegister)
	router.POST("/logout", management.LoginOut)
	router.POST("/getuserinfo", management.GetUserInfo)
	router.GET("/video/category", management.GetCategory) // 查询分类
	router.GET("/video", management.GetVideo)             // 查询视频

	//  用户调用相关接口，中间件对访问次数做限制
	commonUser := router.Group("commonuser/")
	{
		commonUser.POST("/getorderlist")
	}

	// 管理员接口
	manageApi := router.Group("/manager")
	{
		manageApi.GET("/users", management.GetUserInfo)       // 获取当前代理下会员（列表）
		manageApi.POST("/users/edit", management.GetUserInfo) // 编辑会员
	}

	// 超管接口
	// manageApi := router.Group("manage/").Use(middleware.JWTAuth("manager"))
	superAdminApi := router.Group("superadmin/")
	{
		superAdminApi.GET("/users", management.GetUserInfo)                       // 根据用户ID获取用户信息,如果id不存在，则返回用户列表
		superAdminApi.GET("/usersbyjwt", management.GetUserInfo)                  // 根据jwt信息获取用户信息
		superAdminApi.POST("/users/creatmanager", management.CreateManageuser)    // 创建管理员
		superAdminApi.POST("/users/changepwd", management.ChangePwd)              // 更改密码
		superAdminApi.POST("/users/delete", management.DeleteUser)                // 删除管理员
		superAdminApi.POST("/roles/createdit", management.CreateRole)             // 编辑角色
		superAdminApi.POST("/permissions/createdit", management.CreatePermission) // 编辑权限
		superAdminApi.POST("/users/addpremission", management.EditUserPermission) // 管理管理员权限（add or delete）
		superAdminApi.POST("/video/category/create", management.CreateCategory)   // 新建分类
		superAdminApi.POST("/video/category/delete", management.DeleteCategory)   // 删除分类
		superAdminApi.POST("/video/create", management.UploadVideo)               // 新建视频
		superAdminApi.POST("/video/createinfo", management.UploadVideoInfo)       // 新建剧集
		superAdminApi.POST("/video/delete", management.DeleteVideo)               // 删除视频或者剧集
		// 下方为管理收益相关信息

	}
}
