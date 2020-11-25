package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lequocduyquang/cdn-service/db"
	"github.com/lequocduyquang/cdn-service/utils"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

var (
	// ImageController exported
	ImageController ImageControllerInterface = &imageController{}
)

// ImageControllerInterface interface
type ImageControllerInterface interface {
	Upload(c *gin.Context)
}

type imageController struct{}

func (i *imageController) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	defer file.Close()

	bucket, err := gridfs.NewBucket(db.Client.Database("jungleDB"))
	if err != nil {
		log.Fatal(err)
	}
	uploadStream, err := bucket.OpenUploadStream(header.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer uploadStream.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	fileSize, err := uploadStream.Write(data)

	c.JSON(http.StatusOK, gin.H{
		"file_name": header.Filename,
		"file_size": fileSize,
	})
}
