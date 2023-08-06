package management

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"orderingsystem/app/common/request"
	"orderingsystem/app/common/response"
	"orderingsystem/app/services"
	"orderingsystem/global"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateEditShop(c *gin.Context) {
	// 不同请求方式对应不同的处理方式
	var form request.CreateEditShop
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	} else {
		method := c.Request.Method
		if user, err := services.ShopServices.CreateEditShop(&form, method); err != nil {
			response.BusinessFail(c, err.Error())
			return
		} else {
			response.Success(c, user)
		}
		return
	}
}

func GetCategoryList(c *gin.Context) {

}

// 菜品分类相关操作
func CreatEditCategory(c *gin.Context) {
	var form request.Category
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if category, err := services.ShopServices.CreatEditCategory(&form); err != nil {
			response.BusinessFail(c, err.Error())
			return
		} else {
			response.Success(c, category)
		}
		return
	}
}

func UploadImages(c *gin.Context) {
	file, err := c.FormFile("image")
	fmt.Println(file)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件上传失败",
		})
		return
	}

	if !isAllowedImageFormat(file) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件格式不支持",
		})
		return
	}

	// 保存到服务器
	fileSuffix := file.Filename
	new_file_name := uuid.New().String() + fileSuffix

	err = c.SaveUploadedFile(file, "static/images/"+new_file_name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件上传失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "文件上传成功",
		})
	}
	return
}

func UploadImagesTools(c *gin.Context) (err error, imageurls []string) {

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("请上传图片")
		return
	}
	files := form.File["images"]
	for _, file := range files {
		if !isAllowedImageFormat(file) {
			err = errors.New("图片格式不支持")
			return
		}
	}

	for _, file := range files {
		// 保存到服务器
		fileSuffix := file.Filename
		new_file_name := uuid.New().String() + fileSuffix
		err = c.SaveUploadedFile(file, "static/images/"+new_file_name)
		if err != nil {
			err = errors.New("上传失败")
		} else {
			imageurl := global.App.Config.App.AppUrl + ":" + global.App.Config.App.Port + "/" + "static/images/" + new_file_name
			imageurls = append(imageurls, imageurl)
		}
	}
	return
}

func isAllowedImageFormat(file *multipart.FileHeader) bool {
	ext := strings.ToLower(path.Ext(file.Filename))
	allowedFormats := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, format := range allowedFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// 菜品操作
func CreateCuisine(c *gin.Context) {
	var form request.Cuisine
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err, imageurls := UploadImagesTools(c); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		if Cuisine, err := services.ShopServices.CreatCuisine(&form, imageurls); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, Cuisine)
		}
	}
	return
}
