package request

// 不同请求方式对应不同的功能，这里参数在函数逻辑中进行判断
type CreateCategory struct {
	FirstLevel  string `form:"firstlevel" json:"firstlevel"`
	SecondLevel string `form:"secondlevel" json:"secondlevel"`
	Name        string `form:"name" json:"name" binding:"required"`
	Intro       string `form:"intro" json:"intro"`
}

func (category CreateCategory) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required": "分类名称不能为空",
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

type CreateVideoItem struct {
	VideoID  uint   `form:"videoid" json:"videoid" binding:"required"`
	Episodes uint   `form:"episodes" json:"episodes" binding:"required"`
	Url      string `form:"url" json:"url" binding:"required"`
	Intro    string `form:"intro" json:"intro"`
}

func (createvideitem CreateVideoItem) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"videoid:required":  "剧名不能为空",
		"episodes:required": "剧集不能为空",
		"url:required":      "地址不能为空",
	}
}

type UploadVideo struct {
	Name       string   `form:"name" json:"name" binding:"required"`
	Cover      string   `form:"cover" json:"cover" binding:"required"`
	Intro      string   `form:"intro" json:"intro"`
	Categories []string `form:"categories" json:"categories" binding:"required"`
}

func (createvideo UploadVideo) Getmessage() ValidatorMessages {
	return ValidatorMessages{
		"name:required":       "剧名不能为空",
		"cover:required":      "封面不能为空",
		"Categories:required": "分类不能为空",
	}
}

type Deletevideo struct {
	VideoID     uint   `form:"videoid" json:"videoid"`
	VideoItemIDList []uint `form:"videoitemid" json:"videoitemid"`
}
