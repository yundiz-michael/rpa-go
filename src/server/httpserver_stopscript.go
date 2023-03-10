package server

import (
	"github.com/gin-gonic/gin"
	"merkaba/common"
)

func (server *HttpServer) registerStopScript() {
	server.instance.POST("/stopScript", func(c *gin.Context) {
		data := server.parseData(c)
		if data == nil {
			return
		}
		taskName := data["taskName"].(string)
		instance := server.DB.FindMemInstance(taskName)
		isSuccess := instance != nil
		if isSuccess {
			instance.Interrupt()
		}
		common.StopVNC(taskName, "ScriptStop")
		m := successResp()
		m["taskName"] = taskName
		server.writeResponse(c, m)
	})
}
