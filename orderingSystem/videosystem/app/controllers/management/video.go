package management

import (
	"errors"
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/models"
	"orderingsystem/app/services"
	"orderingsystem/global"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(c *gin.Context) {
	var form request.CreateCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := services.VideoServices.CreateCategory(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, "创建成功")
	}
}

func GetCategory(c *gin.Context) {
	categoryid := c.DefaultQuery("categoryid", "0")
	if categories, err := services.VideoServices.GetCategory(categoryid); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		fmt.Println(categories)
		response.Success(c, categories)
	}
}

// 根据类型或者剧名获取剧情详情
// 两个查询条件只能同时存在一个，若同时存在，以videoname为准
func GetVideo(c *gin.Context) {
	videoname := c.DefaultQuery("videoname", "")
	category := c.DefaultQuery("category", "")
	page := c.DefaultQuery("page", "")
	page_size := c.DefaultQuery("page_size", "")
	var pageN int
	var page_sizeN int
	var err error
	if pageN, err = strconv.Atoi(page); err != nil {
		err = errors.New("参数错误")
		response.BusinessFail(c, err.Error())
		return
	}
	if page_sizeN, err = strconv.Atoi(page_size); err != nil {
		err = errors.New("参数错误")
		response.BusinessFail(c, err.Error())
		return
	}
	if videos, count, err := services.VideoServices.GetVideo(videoname, category, pageN, page_sizeN); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		type VideoInfo struct {
			Videolist []models.Video `json:"videos"`
			Count     int            `json:"count"`
		}
		var videinfo VideoInfo
		videinfo.Videolist = videos
		videinfo.Count = count
		response.Success(c, videinfo)
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
