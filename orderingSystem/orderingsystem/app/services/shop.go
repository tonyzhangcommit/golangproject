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

func (ShopServices shopServices) CreateEditShop(parame *request.CreateEditShop, method string) (user models.User, err error) {
	var shop models.Shop
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
			err := global.App.DB.Preload("Shops").Preload("Roles").First(&user, parame.UserId).Error
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
			if err = global.App.DB.Save(&shop).Error; err != nil {
				err = errors.New("店铺更新失败")
				return
			}

		} else if parame.Option == "delete" {
			if err = global.App.DB.First(&shop, parame.Id).Error; err != nil {
				err = errors.New("店铺不存在")
				return
			}
			global.App.DB.Delete(&shop, parame.Id)
		} else {
			err = errors.New("请求参数有误！")
		}
		err = global.App.DB.Preload("Shops").Preload("Roles").First(&user, parame.UserId).Error
	}
	return
}

func (ShopServices shopServices) CreatEditCategory(parame *request.Category) (shop models.Shop, err error) {
	var user models.User
	var category models.Catagory
	if err = global.App.DB.Preload("Shops").First(&user, parame.UserId).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	isRightShop := false
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
		category.ShopID = parame.ShopId
		errCreate := global.App.DB.Create(&category).Error
		err = global.App.DB.Preload("Catagories").First(&shop, parame.ShopId).Error
		if errCreate != nil || err != nil {
			err = errors.New("创建失败")
		}
	} else if parame.Option == "edit" {
		if parame.CategoryId == 0 {
			err = errors.New("缺少分类信息")
			return
		}
		if err = global.App.DB.First(&category, parame.CategoryId).Error; err != nil {
			err = errors.New("分类不存在")
			return
		}
		category.ID.ID = parame.CategoryId
		category.Name = parame.Name
		global.App.DB.Save(&category)
		err = global.App.DB.Preload("Catagories").First(&shop, parame.ShopId).Error
	} else {
		err = errors.New("参数错误")
	}
	return
}

func (ShopServices shopServices) CreatCuisine(parame *request.Cuisine, imageurls []string) (cuisine models.Cuisine, err error) {
	var category models.Catagory
	if err = global.App.DB.First(&category, parame.CatagoryId).Error; err != nil {
		err = errors.New("分类不存在")
		return
	}
	cuisine.CatagoryID = parame.CatagoryId
	cuisine.Name = parame.Name
	cuisine.Price = parame.Price
	cuisine.Peculiarity = parame.Peculiarity
	cuisineObj := global.App.DB.Create(&cuisine)
	if cuisineObj.Error != nil {
		err = cuisineObj.Error
		return
	}
	for _, url := range imageurls {
		imageurl := models.Image{
			CuisineID: cuisine.ID.ID,
			ImageUrl:  url,
		}
		global.App.DB.Create(&imageurl)
	}
	global.App.DB.Preload("Images").First(&cuisine)
	return
}
