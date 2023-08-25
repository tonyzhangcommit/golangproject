package services

import (
	"errors"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
)

type videoServices struct {
}

var VideoServices = new(videoServices)
var db = global.App.DB

func (videoServices videoServices) CreateCategory(params *request.CreateCategory) (category models.Category, err error) {
	if err = db.First(&category, "name = ?", params.Name).Error; err != nil {
		category.Name = params.Name
		category.Intro = params.Intro
		if err = db.Create(&category).Error; err != nil {
			err = errors.New("创建失败")
		}
	} else {
		err = errors.New("分类已存在")
	}
	return
}


