package app

import "github.com/lequocduyquang/cdn-service/controllers"

func mapUrls() {
	router.GET("/ping", controllers.PingController.Ping)
}
