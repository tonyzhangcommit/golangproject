package routers

import (
	"orderingsystem/app/controllers/management"
	"orderingsystem/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiGroupRouters(router *gin.RouterGroup) {
	// 公用接口
	// 获取菜品分类，分类下菜品信息，获取店铺信息
	router.POST("/login", management.Login)

	router.GET("getcategory/")
	router.GET("getmenu/")
	router.GET("getshopinfo/")
	//  小程序相关接口，中间件对访问次数做限制
	miniAppApi := router.Group("miniapp/")
	// 小程序
	// 登录（注册），下单，查看订单状态，订单列表
	{
		miniAppApi.POST("/login")
		miniAppApi.POST("/getorderlist")
		miniAppApi.POST("/takeorder")
	}
	// shopApi := router.Group("shop").Use(middleware.JWTAuth("shopkeeper"))
	shopApi := router.Group("shop")
	// 商户接口
	// 创建店铺, 生成桌号二维码，菜品编辑
	{
		shopApi.POST("/getuserinfo", management.GetUserInfo)
		shopApi.Handle("GET", "/createditshop", management.CreateEditShop)  // 店铺信息
		shopApi.Handle("POST", "/createditshop", management.CreateEditShop) // 新增或修改店铺信息
		shopApi.GET("/categorylist/:shopid", management.GetCategoryList)    // 菜品分类信息，包含菜品信息
		shopApi.POST("/createditcategory", management.CreatEditCategory)    // 新增或修改菜品分类信息
		shopApi.POST("/upload", management.UploadImages)                    // 上传图片
		shopApi.Handle("GET", "/menu", management.CreateCuisine)            // 获取菜品列表
		shopApi.Handle("POST", "/createcuisine", management.CreateCuisine)  // 新建菜品
		shopApi.GET("/tableslist")                                          // 获取座位列表
		shopApi.GET("/tableinfo")                                           // 获取座位详细信息
		shopApi.POST("/createstables")                                      // 创建座位列表
		shopApi.POST("/edittables")                                         // 编辑座位（add edit delete）
		shopApi.POST("/editorderinfo")                                      // 编辑订单，主要是加减菜，使用优惠券等
		shopApi.GET("/getorderlist")                                        // 获取订单列表
		// 运营模块  日，周，月，分析
	}

	manageApi := router.Group("manage/").Use(middleware.JWTAuth("manager"))
	// manageApi := router.Group("manage/")
	// 超管接口
	// 商户编辑，获取商户列表/信息
	{
		manageApi.POST("/getuserinfo", management.GetUserInfo)
		manageApi.POST("/logout", management.LoginOut)
		manageApi.POST("/createuser", management.Createuser)
		manageApi.POST("/edituser", management.Edituser)
		manageApi.GET("/getuserlist") // 代理列表
		manageApi.POST("/editrole", management.CreateRole)
		manageApi.POST("/editper", management.CreatePermission)
		manageApi.Handle("GET", "/editroleper", management.EditRolePermission)
		manageApi.Handle("POST", "/editroleper", management.EditRolePermission)
		manageApi.Handle("DELETE", "/editroleper", management.EditRolePermission)
	}
}
