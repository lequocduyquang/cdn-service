package controllers

import (
	"io"
	"net/http"
	"net/url"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/lequocduyquang/cdn-service/utils"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
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

var (
	storageClient *storage.Client
)

func (i *imageControllerV2) Upload(c *gin.Context) {
	var err error
	bucket := "cdn_lotus_bucket"
	ctx := appengine.NewContext(c.Request)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("lotus-fitness.json"))
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

	sw := storageClient.Bucket(bucket).Object(header.Filename).NewWriter(ctx)

	if _, err := io.Copy(sw, file); err != nil {
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

	resultURL, err := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": resultURL.EscapedPath(),
	})

	/*
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
	*/
}
