package server

import (
	"github.com/gin-gonic/gin"
	"merkaba/common"
)

func (server *HttpServer) registerStartVNC() {
	server.instance.POST("/startVNC", func(c *gin.Context) {
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
		vi := common.StartVNC(common.VncWidth, common.VncHeight, taskName)
		json := successResp()
		json["port"] = vi.Port
		json["width"] = common.VncWidth
		json["height"] = common.VncHeight
		server.writeResponse(c, json)
	})
}
