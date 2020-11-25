package controllers

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	out, err := os.Create("public/" + header.Filename)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	// data, err := ioutil.ReadAll(file)
	c.JSON(http.StatusOK, gin.H{
		"file_name": header.Filename,
	})
}
