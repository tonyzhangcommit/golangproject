package request

// 不同请求方式对应不同的功能，这里参数在函数逻辑中进行判断
type CreateEditShop struct {
	Id      uint   `form:"id" json:"id"`
	UserId  uint   `form:"userid" json:"userid"`
	Name    string `form:"name" json:"name"`
	Address string `form:"address" json:"address"`
	Option  string `form:"option" json:"option"`
}

type Category struct {
	UserId     uint   `form:"userid" json:"userid" binding:"required"`
	ShopId     uint   `form:"shopid" json:"shopid" binding:"required"`
	CategoryId uint   `form:"categoryid" json:"categoryid"`
	Name       string `form:"name" json:"name" binding:"required"`
	Option     string `form:"option" json:"option" binding:"required"`
}

func (category Category) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"userid:required": "用户信息不能为空",
		"shopid:required": "店铺信息不能为空",
		"name:required":   "分类名不能为空",
		"option:required": "option不能为空",
	}
}
