package server

import (
	"github.com/gin-gonic/gin"
)

func (server *HttpServer) registerReadScriptCount() {
	server.instance.POST("/readScriptCount", func(c *gin.Context) {
		data := server.parseData(c)
		includeDetail := false
		if data != nil {
			includeDetail = data["includeDetail"].(bool)
		}
		json := successResp()
		runCount := 0
		idleCount := 0
		items := make([]any, 0)
		for _, i := range server.DB.Instances {
			if i.Status == "Running" {
				runCount += 1
			} else {
				idleCount += 1
			}
			if includeDetail {
				items = append(items, i.AsMap())
			}
		}
		json["runCount"] = runCount
		json["idleCount"] = idleCount
		if includeDetail {
			json["items"] = items
		}
		server.writeResponse(c, json)
	})
}
