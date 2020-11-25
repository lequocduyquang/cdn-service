package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lequocduyquang/cdn-service/db"
)

var (
	router = gin.Default()
)

// StartApp Bootstraping user profile service
func StartApp() {
	db.InitiateMongoClient()
	mapUrls()
	log.Fatal(router.Run(os.Getenv("PORT")))
}
