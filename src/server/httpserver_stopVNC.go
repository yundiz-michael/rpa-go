package server

import (
	"github.com/gin-gonic/gin"
	"merkaba/common"
)

func (server *HttpServer) registerStopVNC() {
	server.instance.POST("/stopVNC", func(c *gin.Context) {
		data := server.parseData(c)
		if data == nil {
			server.writeResponse(c, errorResp("data format error"))
			return
		}
		taskName := server.CheckDataField(c, data, "taskName")
		if len(taskName) == 0 {
			server.writeResponse(c, errorResp("taskName must be required"))
			return
		}
		common.StopVNC(taskName, "WebClose")
		json := successResp()
		server.writeResponse(c, json)
	})
}
