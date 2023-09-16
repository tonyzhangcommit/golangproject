package services

import (
	"errors"
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"strconv"
	"strings"
)

type videoServices struct {
}

var VideoServices = new(videoServices)

// 创建分类
func (videoServices videoServices) CreateCategory(params *request.CreateCategory) (category models.Category, err error) {
	category.Name = params.Name
	category.Intro = params.Intro
	if params.FirstLevel != "" {
		var ferr error
		var firstl int
		var categoryF models.Category
		if firstl, ferr = strconv.Atoi(params.FirstLevel); ferr != nil {
			err = errors.New("参数错误")
			return
		}
		if err = global.App.DB.First(&categoryF, firstl).Error; err != nil {
			err = errors.New("一级分类不存在")
			return
		}
		if params.SecondLevel != "" {
			var serr error
			var secondl int
			var categoryS models.Category
			if secondl, serr = strconv.Atoi(params.SecondLevel); serr != nil {
				err = errors.New("二级参数错误")
				return
			}
			if err = global.App.DB.First(&categoryS, secondl).Error; err != nil {
				err = errors.New("二级分类不存在")
				return
			}
			category.CategoryID = uint(secondl)
			if err = global.App.DB.Create(&category).Error; err != nil {
				err = errors.New("创建失败")
			}
			return
		}
		category.CategoryID = uint(firstl)
		if err = global.App.DB.Create(&category).Error; err != nil {
			err = errors.New("创建失败")
		}
	} else {
		category.CategoryID = 0
		global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
		if err = global.App.DB.Create(&category).Error; err != nil {
			err = errors.New("创建失败")
			return
		}
		global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}

	return
}

// 获取分类
func (videoServices videoServices) GetCategory(categoryID string, categoryType string) (categories []models.Category, err error) {
	var builder strings.Builder
	builder.WriteString("%")
	builder.WriteString(categoryType)
	builder.WriteString("%")
	condition := builder.String()
	if categoryID == "0" {
		if err = global.App.DB.Preload("Categorylist").Where("intro LIKE ?", condition).Find(&categories).Error; err != nil {
			err = errors.New("查询失败")
		}
	} else {
		if categoryId, err := strconv.Atoi(categoryID); err != nil {
			err = errors.New("参数错误")
		} else {
			if err = global.App.DB.Preload("Categorylist").Where("intro LIKE ?", condition).Find(&categories, categoryId).Error; err != nil {
				err = errors.New("查询分类详情失败")
			}
		}
	}
	return
}
func (videoServices videoServices) GetVideo(videoname string, category string) (videos []models.Video, err error) {

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
		if itemInt, err := strconv.Atoi(item); err != nil {
			err = errors.New("分类信息有误")
			var category models.Category
			if dberr = global.App.DB.First(&category, itemInt).Error; dberr != nil {
				dberr = errors.New("分类不存在")
			}
		}
		if err != nil || dberr != nil {
			return
		}
	}
	fmt.Println(video)
	if err = global.App.DB.Where("name = ?", params.Name).First(&video).Error; err == nil {
		err = errors.New("视频已存在")
		return
	}
	if err = global.App.DB.Create(&video).Error; err != nil {
		err = errors.New("创建视频失败，请联系管理员")
	}
	// 添加关联
	for _, item := range params.Categories {
		itemInt, _ := strconv.Atoi(item)
		var category models.Category
		global.App.DB.First(&category, itemInt)
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
	if params.VideoID != 0 {
		// 删除影视
		if err = global.App.DB.First(&video, params.VideoID).Error; err != nil {
			err = errors.New("视频不存在!")
			return
		}
		// 清空关联
		global.App.DB.Model(&video).Association("Categories").Clear()
		global.App.DB.Where("video_id = ?", params.VideoID).Delete(&models.VideoInfo{})
		global.App.DB.Where(&video).Delete(&video)
	} else if len(params.VideoItemIDList) != 0 {
		for _, videoItemId := range params.VideoItemIDList {
			var videoitem models.VideoInfo

			if err = global.App.DB.First(&videoitem, videoItemId).Error; err != nil {
				err = errors.New("剧集不存在!")
			}
			videoID := videoitem.VideoID
			if err = global.App.DB.First(&video, videoID).Error; err != nil {
				err = errors.New("该剧集对应视频不存在!")
			}
			// 删除剧集
			global.App.DB.Delete(&videoitem)
		}

	} else {
		err = errors.New("请填写正确参数!")
	}
	return
}
