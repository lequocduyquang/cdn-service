package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lequocduyquang/cdn-service/db"
	"github.com/lequocduyquang/cdn-service/services"
	"github.com/lequocduyquang/cdn-service/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	// ImageController exported
	ImageController ImageControllerInterface = &imageController{}
)

// ImageControllerInterface interface
type ImageControllerInterface interface {
	Upload(c *gin.Context)
	GetByName(c *gin.Context)
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

func (i *imageController) GetByName(c *gin.Context) {
	fileName := c.Param("filename")
	db := db.Client.Database(os.Getenv("MONGO_DB_NAME"))
	fsFiles := db.Collection("fs.files")
	fsChunks := db.Collection("fs.chunks")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var foundFile, results bson.M
	err := fsFiles.FindOne(ctx, bson.M{"filename": fileName}).Decode(&foundFile)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	filter := bson.M{"files_id": foundFile["_id"]}
	if errFind := fsChunks.FindOne(ctx, filter).Decode(&results); errFind != nil {
		restErr := utils.NewBadRequestError(errFind.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	fmt.Printf("Type of %v", reflect.TypeOf(results["data"]))
	c.JSON(http.StatusOK, gin.H{
		"data": results["data"],
	})
}
