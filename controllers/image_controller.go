package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lequocduyquang/cdn-service/db"
	"github.com/lequocduyquang/cdn-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

	bucket, err := gridfs.NewBucket(db.Client.Database("jungleDB"))
	if err != nil {
		log.Fatal(err)
	}
	randomID := uuid.New()
	fileName := randomID.String() + "-" + header.Filename
	uploadStream, err := bucket.OpenUploadStream(fileName)
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
		"file_name": fileName,
		"file_size": fileSize,
	})
}

func (i *imageController) GetByName(c *gin.Context) {
	fileName := c.Param("filename")
	db := db.Client.Database("jungleDB")
	fsFiles := db.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v\n", dStream)
	ioutil.WriteFile(fileName, buf.Bytes(), 0600)
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
