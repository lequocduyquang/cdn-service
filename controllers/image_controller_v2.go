package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lequocduyquang/cdn-service/services"
	"github.com/lequocduyquang/cdn-service/utils"
)

var (
	// ImageControllerV2 exported
	ImageControllerV2 ImageControllerV2Interface = &imageControllerV2{}
)

// ImageControllerV2Interface interface
type ImageControllerV2Interface interface {
	Upload(c *gin.Context)
}

type imageControllerV2 struct{}

func (i *imageControllerV2) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	defer file.Close()

	fileName, fileSize, err := services.ImageService.Upload(file, *header)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"file_name": fileName,
		"file_size": fileSize,
	})
}
