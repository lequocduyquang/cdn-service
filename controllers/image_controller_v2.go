package controllers

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lequocduyquang/cdn-service/clients"
	"github.com/lequocduyquang/cdn-service/services"
	"github.com/lequocduyquang/cdn-service/utils"
	"google.golang.org/appengine"
)

var (
	// ImageControllerV2 exported
	ImageControllerV2 ImageControllerV2Interface = &imageControllerV2{}
)

// ImageControllerV2Interface interface
type ImageControllerV2Interface interface {
	Upload(c *gin.Context)
	GetByName(c *gin.Context)
	Delete(c *gin.Context)
	UploadMultiple(c *gin.Context)
}

type imageControllerV2 struct{}

func (i *imageControllerV2) UploadMultiple(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	storageClient, _ := clients.InitiateGCPClient(ctx)

	form, _ := c.MultipartForm()
	files := form.File["files"]

	/*
		image service
	*/

	for _, file := range files {
		sw := storageClient.Bucket(os.Getenv("GCP_BUCKET")).Object(file.Filename).NewWriter(ctx)
		data, err := file.Open()

		if _, err := io.Copy(sw, data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		if err := sw.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}

		_, err = url.Parse("/" + os.Getenv("GCP_BUCKET") + "/" + sw.Attrs().Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"Error":   true,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
	})
}

func (i *imageControllerV2) Upload(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	storageClient, err := clients.InitiateGCPClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	defer file.Close()

	resultURL, err := services.GCPService.Upload(ctx, storageClient.Bucket(os.Getenv("GCP_BUCKET")).Object(header.Filename), file)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": resultURL.EscapedPath(),
	})
}

func (i *imageControllerV2) GetByName(c *gin.Context) {
	var err error
	fileName := c.Param("filename")
	ctx := appengine.NewContext(c.Request)
	storageClient, err := clients.InitiateGCPClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	data, err := services.GCPService.Get(ctx, storageClient.Bucket(os.Getenv("GCP_BUCKET")).Object(fileName))
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "get file successfully",
		"data":    data,
	})
}

func (i *imageControllerV2) Delete(c *gin.Context) {
	fileName := c.Param("filename")
	ctx := appengine.NewContext(c.Request)
	storageClient, err := clients.InitiateGCPClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	if err := services.GCPService.Delete(ctx, storageClient.Bucket(os.Getenv("GCP_BUCKET")).Object(fileName)); err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete file successfully",
	})
}
