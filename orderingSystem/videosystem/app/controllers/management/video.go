package management

import (
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/models"
	"orderingsystem/app/services"
	"orderingsystem/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

func GetCategory(c *gin.Context) {
	categoryid := c.DefaultQuery("categoryid", "0")
	categoryT := c.DefaultQuery("categorytype", "")
	if categories, err := services.VideoServices.GetCategory(categoryid, categoryT); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, categories)
	}
}

// 根据类型或者剧名获取剧情详情
// 两个查询条件只能同时存在一个，若同时存在，以videoname为准
func GetVideo(c *gin.Context) {
	videoname := c.DefaultQuery("videoname", "")
	category := c.DefaultQuery("category", "")
	if videos, err := services.VideoServices.GetVideo(videoname, category); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, videos)
	}
}

func deleteCategoryAndChildren(db *gorm.DB, category *models.Category) {
	var children []models.Category
	db.Model(category).Association("Categorylist").Find(&children)
	for _, child := range children {
		deleteCategoryAndChildren(db, &child)
	}
	db.Delete(category)
}

func DeleteCategory(c *gin.Context) {
	var form request.DeleteCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	var category models.Category
	if err := global.App.DB.First(&category, form.ID).Error; err != nil {
		response.BusinessFail(c, "分类不存在")
	} else {
		// 递归删除
		deleteCategoryAndChildren(global.App.DB, &category)
		response.Success(c, "删除成功")

	}
}

// 新建剧集
// 上传剧名，封面，分类，简介
func UploadVideo(c *gin.Context) {
	var form request.UploadVideo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if video, err := services.VideoServices.UploadVideo(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, video)
	}
}

// 更新剧集详情，视频链接，剧集，视频地址，集数
func UploadVideoInfo(c *gin.Context) {
	var form request.CreateVideoItem
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if video, err := services.VideoServices.UploadVideoItem(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, video)
	}
}

// 删除影视/影视中的某一集
func DeleteVideo(c *gin.Context) {
	var form request.Deletevideo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err := services.VideoServices.DeleteVideo(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, "删除成功")
	}
}
