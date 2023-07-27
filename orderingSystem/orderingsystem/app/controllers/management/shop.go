package management

import (
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/services"

	"github.com/gin-gonic/gin"
)

func CreateEditShop(c *gin.Context) {
	// 不同请求方式对应不同的处理方式
	var form request.CreateEditShop
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	} else {
		method := c.Request.Method
		if shop, err := services.ShopServices.CreateEditShop(&form, method); err != nil {
			response.BusinessFail(c, err.Error())
			return
		} else {
			response.Success(c, shop)
		}
		return
	}
}

func GetCategoryList(c *gin.Context){
	
}

// 菜品分类相关操作
func CreatEditCategory(c *gin.Context){
	var form request.Category
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form,err))
		return
	} else {
		if category, err := services.ShopServices.CreatEditCategory(&form); err != nil {
			response.BusinessFail(c, err.Error())
			return
		} else {
			response.Success(c, category)
		}
		return
	}
}

// 菜品操作
func CreatEditMenu(c *gin.Context){

}



