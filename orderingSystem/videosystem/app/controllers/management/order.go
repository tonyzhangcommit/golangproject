package management

import (
	"encoding/json"
	"errors"
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/models"
	"orderingsystem/app/services"
	"orderingsystem/global"
	"time"

	"github.com/gin-gonic/gin"
)

func ProductsList(c *gin.Context) {
	pid := c.DefaultQuery("id", "0")
	if p, err := services.OrderService.ProductsList(pid); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, p)
	}
}

func Products(c *gin.Context) {
	var form request.Product
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if product, err := services.OrderService.Products(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, product)
	}
}

func DelProducts(c *gin.Context) {
	var form request.Delp
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err := services.OrderService.DelProducts(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, "删除成功")
	}
}

// 模拟支付接口
func TempPay() {

}

// 回调查询支付结果
func CallBackRes() {

}

// 用户登录后更新订单信息,在 UpdateOrderInfo 之前
func UpdateOrderInfoByUser(orderID uint, userid uint) (user models.User, order models.Order, err error) {
	if err = global.App.DB.First(&order, orderID).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}
	if err = global.App.DB.First(&user, userid).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	order.UserID = userid
	if err = global.App.DB.Save(&order).Error; err != nil {
		global.App.Log.Error(err.Error() + "UpdateOrderInfoByUser" + string(userid))
	}
	return
}

// 支付成功后更新订单信息
// 不管任何步骤出错，必须将完成的订单信息入到总的支付订单表中
// 总表只保存支付信息和订单号
func UpdateOrderInfo(userID uint, orderID uint, payType string, createT time.Time) (err error) {
	var order models.Order
	var user models.User
	if user, order, err = UpdateOrderInfoByUser(orderID, userID); err != nil {
		return
	}
	var orderInfo models.OrderInfo
	orderInfo.OrderID = orderID
	orderInfo.PayType = payType
	orderInfo.CreatedAt = createT
	if err = global.App.DB.Create(&orderInfo).Error; err != nil {
		err = errors.New("未知错误")
		jsonData, errj := json.Marshal(order)
		if errj != nil {
			global.App.Log.Error(err.Error())
			return
		}
		global.App.Log.Error(string(jsonData))
	}
	// 获取上层所有代理
	UserList := findAncestors(&user)
	// 计算收益逻辑
	err = FinuishInCome(UserList, float64(order.Price), order.OrderCategory)
	return
}

// 支持不登录下单
// 创建订单表，调起支付，回调查询结果，生成订单详情表，更新收益信息
func CreateOrder(c *gin.Context) {
	var form request.CreateOrder
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if order, err := services.OrderService.CreateOrder(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		// 下单成功，开始模拟调起支付...
		// 发起支付成功，等待用户支付
		response.Success(c, order)
	}

}

// 返回订单支持的类型
func GetCategoryOrder(c *gin.Context) {
	data := gin.H{
		"type": []string{"videotype", "menbershiptype"},
	}
	c.JSON(200, data)
}

func GetCategoryPay(c *gin.Context) {
	data := gin.H{
		"type": []string{
			"支付宝",
			"微信"},
	}
	c.JSON(200, data)
}

// 根据用户找出所有上层代理
func findAncestors(user *models.User) (userlist []*models.User) {
	if user == nil || user.ManagerID == 0 {
		userlist = []*models.User{user}
		return
	}
	managerid := user.ManagerID
	var manager models.User
	global.App.DB.First(&manager, managerid)
	uplevelM := findAncestors(&manager)
	return append(uplevelM, user)
}

func TestUser(c *gin.Context) {
	var User models.User
	// var UserList []*models.User
	userId := 17
	global.App.DB.First(&User, userId)
	UserList := findAncestors(&User)
	err := FinuishInCome(UserList, 30.0, "test")
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, "测试成功")
}

func save2d(num float64) float64 {
	return float64(int(num*100)) / 100
}

// 完善收益表
func FinuishInCome(userlist []*models.User, Price float64, OrderCategory string) (err error) {
	managerlist := userlist[:len(userlist)-1]
	lenManage := len(managerlist)
	// 直属代理分70%，剩下30% 为其代理的上层代理平分
	if lenManage == 1 {
		// 用户在超管下
		var incomeinfo models.InComeInfo
		incomeinfo.UserID = managerlist[0].ID.ID
		incomeinfo.InComeType = OrderCategory
		incomeinfo.InComeNum = Price
		if err = global.App.DB.Create(&incomeinfo).Error; err != nil {
			global.App.Log.Error(err.Error() + "superadmin" + "FinuishInCome")
		}
		return
	} else {
		gf_manager_percent := 100 / (len(managerlist) - 1)
		laveP := Price
		all_Income := 0.0
		for i := range managerlist {
			incomeinfo := models.InComeInfo{}
			actualIndex := len(managerlist) - 1 - i
			incomeinfo.InComeType = OrderCategory
			if i == 0 {
				temp_income := save2d(laveP * 70 / 100)
				laveP = laveP * 30 / 100
				incomeinfo.UserID = managerlist[actualIndex].ID.ID
				incomeinfo.InComeNum = temp_income
				all_Income += temp_income
			} else if actualIndex == 0 {
				// 超级管理员的情况，因为float计算存在少量误差，这里将所有小数都精确两位，超级管理员得剩下的金额数
				temp_income := save2d(Price - all_Income)
				incomeinfo.UserID = managerlist[actualIndex].ID.ID
				incomeinfo.InComeNum = temp_income
				all_Income += temp_income
				fmt.Println("superadmin", all_Income, Price)
			} else {
				temp_income := save2d(laveP * float64(gf_manager_percent) / 100)
				// laveP = laveP - temp_income
				incomeinfo.UserID = managerlist[actualIndex].ID.ID
				incomeinfo.InComeNum = temp_income
				all_Income += temp_income
			}
			if err = global.App.DB.Create(&incomeinfo).Error; err != nil {
				global.App.Log.Error(err.Error() + "FinuishInCome")
				return
			}
		}
	}
	return
}
