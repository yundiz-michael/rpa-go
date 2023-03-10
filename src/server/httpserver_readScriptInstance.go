package server

import (
	"github.com/gin-gonic/gin"
	"merkaba/common"
)

func (server *HttpServer) registerReadScriptInstance() {
	server.instance.POST("/readScriptInstance", func(c *gin.Context) {
		data := server.parseData(c)
		if data == nil {
			server.writeResponse(c, errorResp("data format error"))
			return
		}
		scriptUri := server.CheckDataField(c, data, "scriptUri")
		taskName := server.CheckDataField(c, data, "taskName")

		if len(scriptUri) == 0 || len(taskName) == 0 {
			server.writeResponse(c, errorResp("scriptUri and taskName must be required"))
			return
		}
		state := server.DB.ExistMemInstance(scriptUri, taskName)
		json := successResp()
		json["hasVNC"] = common.HasVNC(taskName)
		json["state"] = state
		json["totalCount"] = server.DB.InstanceCount()
		server.writeResponse(c, json)
	})
}
