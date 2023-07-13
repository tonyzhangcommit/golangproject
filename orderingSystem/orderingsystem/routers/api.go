package routers

import "github.com/gin-gonic/gin"

func SetApiGroupRouters(router *gin.RouterGroup){
	// 公用接口
	// 获取菜品分类，分类下菜品信息，获取店铺信息
	router.GET("getcategory/")
	router.GET("getmenu/")
	router.GET("getshopinfo/")
	miniAppApi := router.Group("miniapp/")
	// 小程序
	// 登录（注册），下单，查看订单状态，订单列表
	{
		miniAppApi.POST("/login")
		miniAppApi.POST("/getorderlist")
		miniAppApi.POST("/takeorder")
	}
	shopApi := router.Group("shop")
	// 商户接口
	// 商户登录，生成桌号二维码，菜品编辑
	// 订单编辑，获取餐桌状态
	{
		shopApi.POST("/login")
	}
	manageApi := router.Group("manage/")
	// 超管接口
	// 商户编辑，获取商户列表/信息
	{
		manageApi.POST("/login")
		manageApi.GET("/getuserlist")
	}
}