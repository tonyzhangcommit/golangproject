package management

import (
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/models"
	"orderingsystem/app/services"
	"orderingsystem/global"

	"github.com/gin-gonic/gin"
)

var db = global.App.DB

func CreateCategory(c *gin.Context) {
	var form request.CreateCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if category, err := services.VideoServices.CreateCategory(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, category)
	}
}

func DeleteCategory(c *gin.Context) {
	var form request.DeleteCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	var category models.Category
	if err := db.First(&category, form.ID).Error; err != nil {
		response.BusinessFail(c, "分类不存在")
	} else {
		db.Delete(&category)
		response.Success(c, "删除成功")
	}
}

// 上传影视
func CreateVideo(c *gin.Context) {
	var form request.CreateCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if category, err := services.VideoServices.CreateCategory(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, category)
	}
}

// 删除影视/影视中的某一集
func DeleteVideo(c *gin.Context) {
	var form request.CreateCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if category, err := services.VideoServices.CreateCategory(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, category)
	}
}
