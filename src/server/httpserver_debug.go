package server

import (
	"github.com/gin-gonic/gin"
	"merkaba/common"
	"merkaba/goja"
)

func (server *HttpServer) registerDebug() {
	server.instance.POST("/debugScript", func(c *gin.Context) {
		data := server.parseData(c)
		if data == nil {
			return
		}
		command := data["command"].(string)
		scriptId := data["scriptId"].(string)
		taskName := data["taskName"].(string)
		instance := server.DB.FindMemInstance(taskName)
		if instance == nil {
			server.writeResponse(c, errorResp("不存在"))
			return
		}
		m := successResp()
		m["scriptId"] = scriptId
		server.writeResponse(c, m)
		dCommand := goja.DebugCommand{
			TaskName: taskName,
			Command:  command,
		}
		dCommand.Line = common.ParseIntFrom(data, "line")
		instance.RunVM.Command <- dCommand
	})

}
