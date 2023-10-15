package services

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"orderingsystem/app/common/request"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"strconv"
	"time"
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
			fmt.Println(err)
			err = errors.New("创建失败")
		}
	} else {
		if err = global.App.DB.First(&product, params.PID).Error; err != nil {
			err = errors.New("产品不存在")
			return
		}
		global.App.DB.Model(&product).Where("id = ?", params.PID).Updates(map[string]interface{}{"name": params.Name, "price": params.Price})
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
	orderType := []string{"videotype", "menbertype"}
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
	if params.UserID != 0 {
		// 判断用户是否存在
		if err = global.App.DB.First(&models.User{}, params.UserID).Error; err != nil {
			err = errors.New("用户不存在，请联系管理员")
			return
		}
		order.UserID = params.UserID
	}
	price := params.Price
	productId := params.Product
	if orderCategory == "menbertype" {
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
		if params.UserID == 0 {
			global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
		}
		order.OrderCategory = params.OrderCategory
		order.OrderStatus = "已下单"
		order.Product = productId
		order.Price = price
		if err = global.App.DB.Create(&order).Error; err != nil {
			err = errors.New("下单失败，未知错误，请联系管理员")
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
		if params.UserID == 0 {
			global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
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
	global.App.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	var orderinfo models.OrderInfo
	orderinfo.OrderID = order.ID.ID
	global.App.DB.Create(&orderinfo)
	return
}

// 临时支付接口,支付成功后更改订单状态
func (orderserver orderService) Pay(params *request.Pay) (err error, isPay bool) {
	var order models.Order
	if err = global.App.DB.First(&order, params.OrderID).Error; err != nil {
		err = errors.New("订单不存在！")
	} else {
		// 模拟支付开始
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(10)
		isPay = false
		if randNum > 5 {
			isPay = true
		}
		// 模拟支付结束，pay 是支付状态
		if isPay {
			// 支付成功，更改订单号
			tx := global.App.DB.Begin()
			order.OrderStatus = "已完成"
			if err = tx.Save(&order).Error; err != nil {
				err = errors.New("下单成功，入库order失败")
				tx.Rollback()
				return
			}
			var orderinfo models.OrderInfo
			if err = tx.Where("order_id=?", params.OrderID).First(&orderinfo).Error; err != nil {
				// 发生错误时回滚事务
				tx.Rollback()
				return
			}
			orderinfo.PayType = params.PayType
			if err = tx.Save(&orderinfo).Error; err != nil {
				err = errors.New("下单成功，入库orderinfo失败")
				tx.Rollback()
				return
			}
			tx.Commit()
			if params.UserID != 0 {
				var user models.User
				if err = global.App.DB.First(&user, params.UserID).Error; err != nil {
					err = errors.New("用户不存在!")
					return
				}
				UserList := FindAncestors(&user)
				err = FinuishInCome(UserList, float64(order.Price), order)
			}

		} else {
			err = errors.New("支付调用失败")
		}
	}
	return
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
func FinuishInCome(userlist []*models.User, Price float64, Order models.Order) (err error) {
	managerlist := userlist[:len(userlist)-1]
	lenManage := len(managerlist)
	// 直属代理分70%，剩下30% 为其代理的上层代理平分
	if lenManage == 1 {
		// 用户在超管下
		var incomeinfo models.InComeInfo
		incomeinfo.UserID = managerlist[0].ID.ID
		incomeinfo.OrderType = Order.OrderCategory
		incomeinfo.OrderID = Order.ID.ID
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
		for i, _ := range managerlist {
			incomeinfo := models.InComeInfo{}
			actualIndex := len(managerlist) - 1 - i
			incomeinfo.OrderType = Order.OrderCategory
			incomeinfo.OrderID = Order.ID.ID

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

func save2d(num float64) float64 {
	return float64(int(num*100)) / 100
}

type Orderinfo struct {
	Id          int
	UserId      int
	UserName    string
	OrderType   string
	OrderStatus string
	Product     string
	Price       float64
	PayType     string
	Create_time string
	Pay_time    string
}

func filterByUserId(ordersl []Orderinfo, userid int) (ol []Orderinfo) {
	for _, item := range ordersl {
		if item.UserId == userid {
			ol = append(ol, item)
		}
	}
	return
}

func filterByOrderId(ordersl []Orderinfo, orderid int) (ol []Orderinfo) {
	for _, item := range ordersl {
		if item.Id == orderid {
			ol = append(ol, item)
		}
	}
	return
}

func filterByDate(ordersl []Orderinfo, startD string, endD string) ([]Orderinfo, error) {
	var olEndD []Orderinfo
	var endDate time.Time
	var err error
	if endD == "" {
		currentTime := time.Now()
		nextDay := currentTime.Add(24 * time.Hour)
		endD = nextDay.Format("2006-01-02 15:04:05")
	} else {
		endD = endD + " 23:59:59"
		endDate, err = time.Parse("2006-01-02 15:04:05", endD)
		if err != nil {
			return nil, err
		}
	}
	for _, item := range ordersl {
		temp_startD, err := time.Parse("2006-01-02 15:04:05", item.Create_time)
		if err != nil {
			return nil, err
		}
		if temp_startD.Before(endDate) {
			olEndD = append(olEndD, item)
		}
	}
	if startD == "" {
		return olEndD, nil
	} else {
		var olBeginD []Orderinfo
		startD = startD + " 00:00:00"
		for _, item := range olEndD {
			beginDate, err := time.Parse("2006-01-02 15:04:05", startD)
			if err != nil {
				return nil, err
			}
			temp_startD, err := time.Parse("2006-01-02 15:04:05", item.Create_time)
			if err != nil {
				return nil, err
			}
			if temp_startD.After(beginDate) {
				olBeginD = append(olBeginD, item)
			}
		}
		return olBeginD, nil
	}
}

func GetAllSubUsers(userID int) ([]models.User, error) {
	var users []models.User
	if err := global.App.DB.Where("manager_id = ?", userID).Preload("UserList").Preload("OrderList").Find(&users).Error; err != nil {
		return nil, err
	}
	for i := range users {
		subUsers, err := GetAllSubUsers(int(users[i].ID.ID))
		if err != nil {
			return nil, err
		}
		users = append(users, subUsers...)
	}
	return users, nil
}

func (orderserver orderService) GetOrders(userId string, targetUserId string, beginDate string, endDate string, orderId string, page string, page_size string) (orderslist []Orderinfo, count int, err error) {
	var userid int
	var targetuserid int
	var orderid int
	var userlist []models.User
	var pageN int
	var page_sizeN int
	if pageN, err = strconv.Atoi(page); err != nil {
		err = errors.New("参数错误")
		return
	}
	if page_sizeN, err = strconv.Atoi(page_size); err != nil {
		err = errors.New("参数错误")
		return
	}
	// 检查userid,必填项
	if userid, err = strconv.Atoi(userId); err != nil || userId == "" {
		err = errors.New("参数错误")
		return
	} else {
		userlist, _ = GetAllSubUsers(userid)
		// 服务端处理后续筛选条件
		for _, item := range userlist {
			tempOrderList := item.OrderList
			if len(tempOrderList) > 0 {
				// fmt.Println(item)
				for _, itemOrder := range tempOrderList {
					var orderif models.OrderInfo
					global.App.DB.First(&orderif, itemOrder.ID.ID)
					var pName string
					if itemOrder.OrderCategory == "menbertype" {
						var product models.Products
						global.App.DB.First(&product, itemOrder.Product)
						pName = product.Name
					} else {
						var product models.Video
						global.App.DB.First(&product, itemOrder.Product)
						pName = product.Name
					}
					count += 1
					orderslist = append(orderslist, Orderinfo{
						Id:          int(itemOrder.ID.ID),
						UserId:      int(item.ID.ID),
						UserName:    item.Name,
						OrderType:   itemOrder.OrderCategory,
						OrderStatus: itemOrder.OrderStatus,
						Product:     pName,
						Price:       float64(itemOrder.Price),
						PayType:     orderif.PayType,
						Create_time: itemOrder.CreatedAt.Format("2006-01-02 15:04:05"),
						Pay_time:    orderif.CreatedAt.Format("2006-01-02 15:04:05"),
					})
				}
			}
		}
	}
	// targetUserId 查询指定用户的订单
	if targetUserId != "" {
		if targetuserid, err = strconv.Atoi(targetUserId); err != nil {
			err = errors.New("参数错误")
			return
		}

		orderslist = filterByUserId(orderslist, targetuserid)
	}
	// orderId 根据订单号查询
	if orderId != "" {
		if orderid, err = strconv.Atoi(orderId); err != nil {
			err = errors.New("参数错误")
			return
		}
		orderslist = filterByOrderId(orderslist, orderid)
	}
	if beginDate != "" || endDate != "" {
		orderslist, err = filterByDate(orderslist, beginDate, endDate)
	}
	result := GetPageData(orderslist, page_sizeN, pageN)
	orderslist = result.([]Orderinfo)
	return
}

// 获取收益
func (orderserver orderService) GetIncomes(userId string, page string, page_size string) (incomelist []models.InComeInfo, count int, err error) {
	var pageN int
	var page_sizeN int
	if pageN, err = strconv.Atoi(page); err != nil {
		err = errors.New("参数错误")
		return
	}
	if page_sizeN, err = strconv.Atoi(page_size); err != nil {
		err = errors.New("参数错误")
		return
	}
	if userid, err := strconv.Atoi(userId); err != nil || userId == "" {
		err = errors.New("参数错误")
	} else {
		var user models.User
		if err = global.App.DB.First(&user, userid).Preload("Roles").Error; err != nil {
			err = errors.New("用户不存在")
		}
		global.App.DB.Where("user_id = ?", userid).Find(&incomelist)
		count = len(incomelist)
		result := GetPageData(incomelist, page_sizeN, pageN)
		incomelist = result.([]models.InComeInfo)
	}

	return
}
