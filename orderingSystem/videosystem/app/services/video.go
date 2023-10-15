package services

import (
	"errors"
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"

	// "orderingsystem/app/services"
	"orderingsystem/global"
	"strconv"
)

type videoServices struct {
}

var VideoServices = new(videoServices)

func isValidCategory(firstL, secondL, thirdL string) bool {
	if thirdL != "" {
		return firstL != "" && secondL != ""
	} else if secondL != "" {
		return firstL != ""
	} else {
		return firstL != ""
	}
}

// 创建分类,参数只有三种情况
func (videoServices videoServices) CreateCategory(params *request.CreateCategory) (err error) {
	// 首先判断参数正确性
	if !isValidCategory(params.FirstLevel, params.SecondLevel, params.Thirdlevel) {
		err = errors.New("参数错误")
		return
	}
	var categoryF models.Category
	var categoryS models.Category
	var categoryT models.Category
	// 首先创建一级分类
	if params.FirstLevel != "" {
		if err = global.App.DB.Where("name = ?", params.FirstLevel).First(&categoryF).Error; err != nil {
			categoryF.CategoryID = 0
			global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
			categoryF.Name = params.FirstLevel
			if err = global.App.DB.Create(&categoryF).Error; err != nil {
				err = errors.New("创建一级分类失败")
				return
			}
			global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
		}
	} else {
		err = errors.New("参数错误")
		return
	}
	// 创建二级分类
	if params.SecondLevel != "" {
		if err = global.App.DB.Where("name = ?", params.SecondLevel).First(&categoryS).Error; err != nil {
			categoryS.CategoryID = categoryF.ID.ID
			categoryS.Name = params.SecondLevel
			if err = global.App.DB.Create(&categoryS).Error; err != nil {
				err = errors.New("创建二级分类失败")
				return
			}
		}
	}
	// 创建三级分类
	if params.Thirdlevel != "" {
		if err = global.App.DB.Where("name = ?", params.Thirdlevel).First(&categoryS).Error; err != nil {
			categoryT.CategoryID = categoryS.ID.ID
			categoryT.Name = params.Thirdlevel
			if err = global.App.DB.Create(&categoryT).Error; err != nil {
				err = errors.New("创建三级分类失败")
				return
			}
		} else {
			err = errors.New("三级分类已存在")
			return
		}
	}
	return
}

func GetAllSubcategories(categoryID uint, categories *[]models.Category) {
	global.App.DB.Where("category_id = ?", categoryID).Find(categories)
	for i := range *categories {
		GetAllSubcategories((*categories)[i].ID.ID, &(*categories)[i].Categorylist)
	}
}

// 获取分类
func (videoServices videoServices) GetCategory(categoryID string) (categories []models.Category, err error) {

	if categoryId, err := strconv.Atoi(categoryID); err != nil {
		err = errors.New("参数错误")
	} else {
		if categoryId == 0 {
			// 获取一级分类
			if err = global.App.DB.Preload("Categorylist").Where("category_id = ?", categoryId).Find(&categories).Error; err != nil {
				err = errors.New("查询分类详情失败")
			}
			for i := range categories {
				GetAllSubcategories(categories[i].ID.ID, &categories[i].Categorylist)
			}
		} else {
			if err = global.App.DB.Preload("Categorylist").Find(&categories, categoryId).Error; err != nil {
				err = errors.New("查询分类详情失败")
			}
		}

	}
	return
}
func (videoServices videoServices) GetVideo(videoname string, category string, pageN, page_sizeN int) (videos []models.Video, count int, err error) {
	fmt.Println(pageN, page_sizeN)

	if videoname != "" {
		global.App.DB.Where("name like ?", "%"+videoname+"%").Preload("Videolist").Find(&videos)
	} else if category != "" {
		var tempVideo []models.Video
		global.App.DB.Preload("Videolist").Preload("Categories").Find(&tempVideo)
		for _, item := range tempVideo {
			categories := item.Categories

			for _, categoryitem := range categories {
				if categoryitem.Name == category {
					videos = append(videos, item)
				}
			}
		}
	} else {
		global.App.DB.Preload("Videolist").Preload("Categories").Find(&videos)
	}
	count = len(videos)
	result := GetPageData(videos, page_sizeN, pageN)
	videos = result.([]models.Video)
	return
}

func (videoServices videoServices) UploadVideo(params *request.UploadVideo) (video models.Video, err error) {
	video.Name = params.Name
	video.Cover = params.Cover
	video.Intro = params.Intro

	if params.Intro == "" {
		video.Intro = "暂无简介"
	}
	for _, item := range params.Categories {
		var dberr error
		var category models.Category
		if dberr = global.App.DB.First(&category, item).Error; dberr != nil {
			dberr = errors.New("分类不存在")
		}
		if err != nil || dberr != nil {
			return
		}
	}
	if err = global.App.DB.Where("name = ?", params.Name).First(&video).Error; err == nil {
		err = errors.New("视频已存在")
		return
	}
	if err = global.App.DB.Create(&video).Error; err != nil {
		err = errors.New("创建视频失败，请联系管理员")
	}
	// 添加关联
	for _, item := range params.Categories {
		var category models.Category
		global.App.DB.First(&category, item)
		global.App.DB.Model(&video).Association("Categories").Append(&category)
	}
	return
}

// 新建剧集
func (videoServices videoServices) UploadVideoItem(params *request.CreateVideoItem) (video models.Video, err error) {
	var videoitem models.VideoInfo
	if err = global.App.DB.First(&video, params.VideoID).Error; err != nil {
		err = errors.New("视频不存在，请新建视频")
	} else {
		videoitem.Episodes = int(params.Episodes)
		videoitem.Url = params.Url
		videoitem.Intro = params.Intro
		videoitem.VideoID = params.VideoID
		if err = global.App.DB.Create(&videoitem).Error; err != nil {
			err = errors.New("新建剧集失败")
			return
		}
	}
	global.App.DB.Model(&models.Video{}).Preload("Videolist").First(&video)
	return
}

// 删除影视/影视中的某一集
func (videoServices videoServices) DeleteVideo(params *request.Deletevideo) (err error) {
	var video models.Video
	var videoitem models.VideoInfo
	if params.VideoItemID != 0 {
		if err = global.App.DB.First(&videoitem, params.VideoItemID).Error; err != nil {
			err = errors.New("剧集不存在!")
			return
		}
		global.App.DB.Delete(&videoitem)
		// 删除影视，通过判断当前影视下是否还包含剧集，如果没有则删除影视

		if err = global.App.DB.Preload("Videolist").First(&video, params.VideoID).Error; err != nil {
			err = errors.New("视频不存在!")
			return
		}
		if len(video.Videolist) == 0 {
			global.App.DB.Where(&video).Delete(&video)
			global.App.DB.Model(&video).Association("Categories").Clear()
		}
		// global.App.DB.Model(&video).Association("Categories").Clear()
		// global.App.DB.Where("video_id = ?", params.VideoID).Delete(&models.VideoInfo{})
		// global.App.DB.Where(&video).Delete(&video)
		// } else if len(params.VideoItemID) != 0 {
		// 	for _, videoItemId := range params.VideoItemID  {

		// 		if err = global.App.DB.First(&videoitem, videoItemId).Error; err != nil {
		// 			err = errors.New("剧集不存在!")
		// 		}
		// 		videoID := videoitem.VideoID
		// 		if err = global.App.DB.First(&video, videoID).Error; err != nil {
		// 			err = errors.New("该剧集对应视频不存在!")
		// 		}
		// 		// 删除剧集
		// 		// global.App.DB.Delete(&videoitem)
		// 	}

		// } else {
		// 	err = errors.New("请填写正确参数!")
	} else {
		if err = global.App.DB.Preload("Videolist").First(&video, params.VideoID).Error; err != nil {
			err = errors.New("视频不存在!")
			return
		}
		if len(video.Videolist) == 0 {
			global.App.DB.Where(&video).Delete(&video)
			global.App.DB.Model(&video).Association("Categories").Clear()
		}
	}
	return
}
