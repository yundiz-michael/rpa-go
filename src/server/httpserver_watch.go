package server

import (
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) registerWatch() {
	server.instance.POST("/watchScript", func(c *gin.Context) {
		data := server.parseData(c)
		if data == nil {
			return
		}
		scriptId := data["scriptId"].(string)
		taskName := data["taskName"].(string)
		instance := server.DB.FindMemInstance(taskName)
		if instance == nil {
			server.writeResponse(c, errorResp("不存在"))
			return
		}
		variables := data["variables"].([]string)
		instance.Variables = variables
		values := make(map[string]any)
		for _, name := range variables {
			values[name] = instance.RunVM.Runtime.Get(name).String()
		}
		m := successResp()
		m["scriptId"] = scriptId
		m["values"] = values
		server.writeResponse(c, m)
	})
}
