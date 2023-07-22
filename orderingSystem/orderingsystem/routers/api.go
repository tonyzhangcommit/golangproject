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
	router.POST("/logout", management.LoginOut)
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
	shopApi := router.Group("shop").Use(middleware.JWTAuth("shopkeeper"))
	// 商户接口
	// 创建店铺, 生成桌号二维码，菜品编辑
	{
		shopApi.POST("/getuserinfo", management.GetUserInfo)
		shopApi.Handle("GET", "/shopinfo") // 店铺信息
		shopApi.Handle("POST", "/shopinfo")
		shopApi.POST("/getorderlist") // 获取订单列表/单个订单信息
		shopApi.POST("/gettablelist") // 获取座位列表/单个座位信息
		shopApi.Handle("GET", "/menu")
		shopApi.Handle("POST", "/menu")
		shopApi.Handle("DELETE", "/menu")
	}

	manageApi := router.Group("manage/").Use(middleware.JWTAuth("manager"))
	// 超管接口
	// 商户编辑，获取商户列表/信息
	{
		manageApi.POST("/getuserinfo", management.GetUserInfo)
		manageApi.POST("/createuser", management.Createuser)
		manageApi.POST("/edituser", management.Edituser)
		manageApi.GET("/getuserlist") // 代理列表
		manageApi.POST("/createrole", management.CreateRole)
		manageApi.POST("/createper", management.CreatePermission)
		manageApi.Handle("GET", "/editroleper", management.EditRolePermission)
		manageApi.Handle("POST", "/editroleper", management.EditRolePermission)
		manageApi.Handle("DELETE", "/editroleper", management.EditRolePermission)
	}
}
