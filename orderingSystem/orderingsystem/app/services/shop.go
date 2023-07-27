package services

import (
	"errors"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
)

type shopServices struct {
}

var ShopServices = new(shopServices)

func (ShopServices shopServices) CreateEditShop(parame *request.CreateEditShop, method string) (shop models.Shop, err error) {

	if method == "GET" {
		// 需要店铺ID
		if parame.Id == 0 {
			err = errors.New("请求参数有误！")
			return
		}
		err = global.App.DB.First(&shop, parame.Id).Error
		if err != nil {
			err = errors.New("店铺不存在，请首先新建店铺")
		}
	} else if method == "POST" {
		if parame.Option == "create" {
			// 必须有userid
			if parame.UserId == 0 {
				err = errors.New("请求参数有误！")
				return
			}
			shop = models.Shop{UserID: parame.UserId, Name: parame.Name, Address: parame.Address}
			result := global.App.DB.Create(&shop)

			err = global.App.DB.Model(&models.User{}).Where("id=?", parame.UserId).Association("Shops").Append(&shop)
			if result.Error != nil || err != nil {
				err = errors.New("创建失败")
			}

		} else if parame.Option == "edit" {
			// 必须有shopid
			err = global.App.DB.First(&shop, parame.Id).Error
			if err != nil {
				err = errors.New("店铺不存在")
				return
			}
			shop.Address = parame.Address
			shop.Name = parame.Name
			err = global.App.DB.Save(&shop).Error
			if err != nil {
				err = errors.New("店铺更新失败")
			}
		} else if parame.Option == "delete" {
			err = errors.New("请求参数有误！")
		} else {
			err = errors.New("请求参数有误！")
		}
	}
	return
}

func (ShopServices shopServices) CreatEditCategory(parame *request.Category) (category models.Catagory, err error) {
	var user models.User
	if err = global.App.DB.Preload("Shops").First(&user, parame.UserId).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	isRightShop := false
	var shop models.Shop
	for _, s := range user.Shops {
		if s.ID.ID == parame.ShopId {
			isRightShop = true
			shop = s
			break
		}
	}
	if !isRightShop {
		err = errors.New("店铺不存在")
		return
	}
	if parame.Option == "create" {
		category.Name = parame.Name
		if err = global.App.DB.Model(&shop).Association("Catagories").Append(&category); err != nil {
			err = errors.New("创建失败")
			return
		}
	} else if parame.Option == "edit" {
		if parame.CategoryId == 0 {
			err = errors.New("缺少分类信息")
			return
		}
		if err = global.App.DB.First(&category, parame.CategoryId).Error; err != nil {
			err = errors.New("分类不存在")
		}
		category.ID.ID = parame.CategoryId
		category.Name = parame.Name
		if err = global.App.DB.Save(&user).Error; err != nil {
			err = errors.New("更新失败")
		}
	} else {
		err = errors.New("参数错误")
	}
	return
}
