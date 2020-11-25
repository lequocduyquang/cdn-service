package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApp Bootstraping user profile service
func StartApp() {
	mapUrls()
	log.Fatal(router.Run(os.Getenv("PORT")))
}
