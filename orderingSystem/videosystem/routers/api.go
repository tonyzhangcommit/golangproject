package routers

import (
	"orderingsystem/app/controllers/management"
	"orderingsystem/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiGroupRouters(router *gin.RouterGroup) {
	// 公用接口
	// 登录，登出，获取用户信息（用户名，电话，注册时间）
	router.POST("/login", management.Login)
	router.POST("/register", management.CommonRegister)         // 普通注册
	router.POST("/usersinfobyjwt", management.GetUserInfoByJwt) // 根据jwt信息获取用户信息
	router.POST("/getuserinfo", management.GetUserInfo)
	router.GET("/video/category", management.GetCategory) // 查询分类
	router.GET("/video", management.GetVideo)             // 查询视频
	// 订单相关
	router.GET("/ordercategory", management.GetCategoryOrder) // 获取订单类型(已完成)
	router.GET("/paycategory", management.GetCategoryPay)     // 获取支付类型(已完成)
	router.POST("/createorder", management.CreateOrder)       // 下单(待测试)
	router.GET("/test", management.TestUser)                  // 测试接口

	//  用户调用相关接口，中间件对访问次数做限制
	commonUser := router.Group("common/").Use(middleware.JWTAuth(""))
	{
		commonUser.POST("/getorderlist")
		commonUser.POST("/logout", management.LoginOut)
	}

	// 管理员接口
	// 获取当前自代理列表
	// 获取当前所有订单（包含自代理订单）
	// manageApi := router.Group("/manager").Use(middleware.JWTAuth("manager"))
	manageApi := router.Group("/manager")
	{
		manageApi.GET("/users", management.GetManagerUsers)   // 获取当前代理下会员（列表）
		manageApi.POST("/users/edit", management.ChangeProxy) // 编辑用户，给用户增加代理，封禁用户等

	}

	// 超管接口
	// superAdminApi := router.Group("superadmin/").Use(middleware.JWTAuth("superadmin"))
	superAdminApi := router.Group("superadmin/")
	{
		superAdminApi.GET("/users", management.GetUserInfo)                       // 根据用户ID获取用户信息,如果id不存在，则返回用户列表(已完成)
		superAdminApi.POST("/users/creatmanager", management.CreateManageuser)    // 创建管理员(已完成)
		superAdminApi.POST("/users/changepwd", management.ChangePwd)              // 更改密码(已完成)
		superAdminApi.POST("/users/delete", management.DeleteUser)                // 删除管理员(已完成)
		superAdminApi.POST("/roles/createdit", management.CreateRole)             // 编辑角色(已完成)
		superAdminApi.POST("/permissions/createdit", management.CreatePermission) // 编辑权限(已完成)
		superAdminApi.POST("/users/addpremission", management.EditUserPermission) // 管理管理员权限（add or delete）(已完成)
		superAdminApi.POST("/video/category/create", management.CreateCategory)   // 新建分类(已完成)
		superAdminApi.POST("/video/category/delete", management.DeleteCategory)   // 删除分类(已完成)
		superAdminApi.POST("/video/create", management.UploadVideo)               // 新建视频(已完成)
		superAdminApi.POST("/video/createinfo", management.UploadVideoInfo)       // 新建剧集(已完成)
		superAdminApi.POST("/video/delete", management.DeleteVideo)               // 删除视频或者剧集(已完成)
		superAdminApi.Handle("GET", "/order/product", management.ProductsList)    // 查询当前所有产品(已完成)
		superAdminApi.Handle("POST", "/order/product", management.Products)       // 编辑（创建或者编辑）产品(已完成)
		superAdminApi.Handle("POST", "/order/delproduct", management.DelProducts) // 删除产品(已完成)
	}
}
