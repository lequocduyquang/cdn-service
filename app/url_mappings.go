package app

import "github.com/lequocduyquang/cdn-service/controllers"

func mapUrls() {
	router.GET("/ping", controllers.PingController.Ping)

	router.POST("/upload", controllers.ImageController.Upload)
	router.GET("/image/:filename", controllers.ImageController.GetByName)
	router.POST("/v2/upload", controllers.ImageControllerV2.Upload)
}
