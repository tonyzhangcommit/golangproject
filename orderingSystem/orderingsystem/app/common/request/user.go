package request

// 这里参数验证有两个，一个为超级管理员，一个为商户
// common/request/user  中存放涉及用户请求中除实际逻辑处理中的所有逻辑

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

type Login struct {
	Name   string `form:"name" json:"name" binding:"required"`
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
}

func (login Login) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required":   "用户名不能为空",
		"mobile:required": "手机号不能为空",
		"mobile:mobile":   "手机号格式错误",
	}
}

//  创建角色
type CreateRole struct {
	Name string `json:"name" binding:"required"`
}

func (createRole CreateRole) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required": "角色名不能为空",
	}
}

// 创建权限
type CreatePermission struct {
	Name string `json:"name" binding:"required"`
}

func (createPermission CreatePermission) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required": "权限名不能为空",
	}
}

// 编辑角色权限信息，CRUD,针对不同的请求方式，对权限进行更改type CreatePermission struct {
type EditRolePermission struct {
	Rolename       string `json:"rolename" binding:"required"`
	Permissionname string `json:"permissionname" binding:"required"`
}

func (editRolePermission EditRolePermission) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"rolename:required":       "角色名不能为空",
		"permissionname:required": "权限名不能为空",
	}
}
