package server

import (
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) registerInfo() {
	server.instance.GET("/info", func(c *gin.Context) {

		server.writeResponse(c, successResp())
	})
}
