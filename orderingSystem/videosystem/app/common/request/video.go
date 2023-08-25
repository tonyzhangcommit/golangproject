package request

// 不同请求方式对应不同的功能，这里参数在函数逻辑中进行判断
type CreateCategory struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Intro string `form:"intro" json:"intro" binding:"required"`
}

func (category CreateCategory) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required":  "分类名称不能为空",
		"intro:required": "分类信息不能为空",
	}
}

type DeleteCategory struct {
	ID string `form:"id" json:"id" binding:"required"`
}

func (category DeleteCategory) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"id:required": "分类ID不能为空",
	}
}

type CreateVideo struct {
	Name  string  `form:"name" json:"name" binding:"required"`
	Cover float32 `form:"cover" json:"cover" binding:"required"`
	Intro float32 `form:"intro" json:"intro" binding:"required"`
	
}

func (createvideo CreateVideo) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required":  "剧名不能为空",
		"cover:required": "封面不能为空",
		"intro:required": "简介不能为空",
	}
}

type CreateVideoItem struct {
	Episodes string  `form:"episodes" json:"episodes"`
	Url      float32 `form:"url" json:"url"`
	Intro    float32 `form:"intro" json:"intro"`
}
