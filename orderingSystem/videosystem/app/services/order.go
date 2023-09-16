package services

import (
	"errors"
	"math"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"strconv"
)

type orderService struct {
}

var OrderService = new(orderService)

func (orderserver orderService) Products(params *request.Product) (product models.Products, err error) {
	if params.PID == 0 {
		if err = global.App.DB.Where("name = ?", params.Name).First(&product).Error; err == nil {
			err = errors.New("产品已存在")
			return
		}
		product.Name = params.Name
		product.Price = params.Price
		if err = global.App.DB.Create(&product).Error; err != nil {
			err = errors.New("创建失败")
		}
	} else {
		if err = global.App.DB.First(&product, params.PID).Error; err != nil {
			err = errors.New("产品不存在")
			return
		}
		product.Name = params.Name
		product.Price = params.Price
		global.App.DB.Save(&product)
	}
	return
}

func (orderserver orderService) ProductsList(pid string) (res interface{}, err error) {
	if pid == "0" {
		var plist []models.Products
		global.App.DB.Find(&plist)
		res = plist
	} else {
		if pid, err := strconv.Atoi(pid); err != nil {
			err = errors.New("参数错误")
		} else {
			var p models.Products
			if err = global.App.DB.First(&p, pid).Error; err != nil {
				err = errors.New("产品不存在")
			}
			res = p
		}
	}
	return
}

func (orderserver orderService) DelProducts(pid *request.Delp) (err error) {
	var p models.Products
	err = global.App.DB.Delete(&p, pid.Pid).Error
	return
}

func (orderserver orderService) CreateOrder(params *request.CreateOrder) (order models.Order, err error) {
	// 会员订单/视频订单
	orderType := []string{"videotype", "menbershiptype"}
	orderCategory := params.OrderCategory
	isRightOT := false
	for _, itme := range orderType {
		if orderCategory == itme {
			isRightOT = true
			break
		}
	}
	if !isRightOT {
		err = errors.New("参数错误，下单失败")
		return
	}
	price := params.Price
	productId := params.Product
	if orderCategory == "menbershiptype" {
		PItem := models.Products{}
		if err = global.App.DB.First(&PItem, productId).Error; err != nil {
			err = errors.New("商品不存在!")
			return
		}
		if math.Abs(float64(price-PItem.Price)) > 0.1 {
			// 二次验证，为了防止恶意更改请求价格
			err = errors.New("参数错误，下单失败!")
			return
		}
		// 创建订单
		if params.UserID != 0 {
			order.UserID = params.UserID
		}
		order.OrderCategory = params.OrderCategory
		order.OrderStatus = "已下单"
		order.Product = productId
		order.Price = price
		if err = global.App.DB.Create(&order).Error; err != nil {
			err = errors.New("下单失败，未知错误，请联系关联员")
			global.App.Log.Error(err.Error())
		}
	} else {
		// 视频价格暂定每部10元
		StandPrice := 10.0
		PItem := models.Video{}
		if err = global.App.DB.First(&PItem, productId).Error; err != nil {
			err = errors.New("视频不存在!")
			return
		}
		if math.Abs(StandPrice-float64(params.Price)) > 0.1 {
			err = errors.New("参数错误，下单失败!")
			return
		}
		if params.UserID != 0 {
			order.UserID = params.UserID
		}
		order.OrderCategory = params.OrderCategory
		order.OrderStatus = "已下单"
		order.Product = productId
		order.Price = price
		if err = global.App.DB.Create(&order).Error; err != nil {
			err = errors.New("下单失败，未知错误，请联系关联员")
			global.App.Log.Error(err.Error())
		}
	}
	return
}
