package management

import (
	"errors"
	"fmt"
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/models"
	"orderingsystem/app/services"
	"orderingsystem/global"

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
func TempPay(c *gin.Context) {
	var form request.Pay
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err, _ := services.OrderService.Pay(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, "支付成功")
	}
}

// 回调查询支付结果
func CallBackRes() {

}

// 未登录用户支付成功后，修改订单信息，用于后续计算各级代理分成情况
func Changeorderstatus(c *gin.Context) {
	var form request.ChangeOrderStatus
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if user, order, err := UpdateOrderInfoByUser(form.OrderID, form.UserID); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		UserList := FindAncestors(&user)
		err := FinuishInCome(UserList, float64(order.Price), order.OrderCategory)
		if err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, "更新成功")
		}
	}
}

// 用户登录后更新订单信息,在 UpdateOrderInfo 之前,发生在支付成功之后
func UpdateOrderInfoByUser(orderID uint, userid uint) (user models.User, order models.Order, err error) {
	if err = global.App.DB.First(&order, orderID).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}
	if err = global.App.DB.First(&user, userid).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	if order.UserID == userid {
		// 说明已经更新过状态
		err = errors.New("订单已更新")
		return
	}
	order.UserID = userid
	if err = global.App.DB.Save(&order).Error; err != nil {
		global.App.Log.Error(err.Error() + "UpdateOrderInfoByUser" + fmt.Sprint(userid))
	}
	return
}

// 支付成功后更新订单信息
// 不管任何步骤出错，必须将完成的订单信息入到总的支付订单表中
// 总表只保存支付信息和订单号

// 支持不登录下单,创建订单表，调起支付，回调查询结果，生成订单详情表，更新收益信息
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
		"type": []string{"videotype", "membertype"},
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
func FindAncestors(user *models.User) (userlist []*models.User) {
	if user == nil || user.ManagerID == 0 {
		userlist = []*models.User{user}
		return
	}
	managerid := user.ManagerID
	var manager models.User
	global.App.DB.First(&manager, managerid)
	uplevelM := FindAncestors(&manager)
	return append(uplevelM, user)
}

func TestUser(c *gin.Context) {
	var User models.User
	// var UserList []*models.User
	userId := 17
	global.App.DB.First(&User, userId)
	UserList := FindAncestors(&User)
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
		incomeinfo.OrderType = OrderCategory
		incomeinfo.InComeNum = Price
		incomeinfo.InComeType = "DirectEarnings"
		if err = global.App.DB.Create(&incomeinfo).Error; err != nil {
			global.App.Log.Error(err.Error() + "superadmin" + "FinuishInCome")
		}
		return
	} else {
		gf_manager_percent := 100 / (len(managerlist) - 1)
		laveP := Price
		all_Income := 0.0
		for i, user := range managerlist {
			fmt.Println(i, *user)
			incomeinfo := models.InComeInfo{}
			actualIndex := len(managerlist) - 1 - i
			incomeinfo.OrderType = OrderCategory
			if i == 0 {
				temp_income := save2d(laveP * 70 / 100)
				laveP = laveP * 30 / 100
				incomeinfo.UserID = managerlist[actualIndex].ID.ID
				incomeinfo.InComeNum = temp_income
				incomeinfo.InComeType = "DirectEarnings"
				all_Income += temp_income
			} else if actualIndex == 0 {
				// 超级管理员的情况，因为float计算存在少量误差，这里将所有小数都精确两位，超级管理员得剩下的金额数
				temp_income := save2d(Price - all_Income)
				incomeinfo.UserID = managerlist[actualIndex].ID.ID
				incomeinfo.InComeNum = temp_income
				incomeinfo.InComeType = "AgencyEarnings"
				all_Income += temp_income
			} else {
				temp_income := save2d(laveP * float64(gf_manager_percent) / 100)
				// laveP = laveP - temp_income
				incomeinfo.UserID = managerlist[actualIndex].ID.ID
				incomeinfo.InComeNum = temp_income
				incomeinfo.InComeType = "AgencyEarnings"
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

// 获取订单,get请求
func GetOrders(c *gin.Context) {
	userId := c.Query("user_id")
	targetUserId := c.Query("t_user_id")
	beginDate := c.Query("start_date")
	endDate := c.Query("end_date")
	orderId := c.Query("order_id")
	page := c.Query("page")
	page_size := c.Query("page_size")

	if orderlist, count, err := services.OrderService.GetOrders(userId, targetUserId, beginDate, endDate, orderId, page, page_size); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		type orderl struct {
			Orderlist []services.Orderinfo `json:"orderlist"`
			Count     int                  `json:"count"`
		}
		response.Success(c, orderl{
			Orderlist: orderlist,
			Count:     count,
		})
	}
}

// 获取收益
func GetIncomes(c *gin.Context) {
	userId := c.Query("user_id")
	page := c.Query("page")
	page_size := c.Query("page_size")
	if incomelist, count, err := services.OrderService.GetIncomes(userId, page, page_size); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		type incomel struct {
			IncomeL []models.InComeInfo `json:"incomelist"`
			Count   int                 `json:"count"`
		}
		response.Success(c, incomel{
			IncomeL: incomelist,
			Count:   count,
		})
	}
}
