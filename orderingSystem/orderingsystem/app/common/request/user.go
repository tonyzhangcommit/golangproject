package request

// 这里参数验证有两个，一个为超级管理员，一个为商户

type Resister struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (register Resister) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name.required":     "用户名称不能为空",
		"mobile.required":   "手机号码不能为空",
		"mobile.mobile":     "手机号码格式错误",
		"password.required": "用户密码不能为空",
	}
}
