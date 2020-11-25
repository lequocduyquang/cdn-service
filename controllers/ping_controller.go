package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	pong = "PONG"
)

var (
	// PingController export controller
	PingController PingControllerInterface = &pingController{}
)

// PingControllerInterface interface
type PingControllerInterface interface {
	Ping(c *gin.Context)
}

type pingController struct{}

func (p *pingController) Ping(c *gin.Context) {
	c.String(http.StatusOK, pong)
}
