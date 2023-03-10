package server

import (
	"github.com/gin-gonic/gin"
	"merkaba/common"
	"strings"
)

func (server *HttpServer) CheckDataField(c *gin.Context, data map[string]interface{}, fieldName string) string {
	if v, ok := data[fieldName]; ok {
		return v.(string)
	} else {
		server.writeResponse(c, errorResp("data miss "+fieldName))
		return ""
	}
}

func (server *HttpServer) registerRunScript() {
	server.instance.POST("/runScript", func(c *gin.Context) {
		data := server.parseData(c)
		if data == nil {
			server.writeResponse(c, errorResp("data format error"))
			return
		}
		scriptId := server.CheckDataField(c, data, "scriptId")
		if len(scriptId) == 0 {
			return
		}
		scriptUri := server.CheckDataField(c, data, "scriptUri")
		if len(scriptUri) == 0 {
			return
		}
		taskName := server.CheckDataField(c, data, "taskName")
		if len(taskName) == 0 {
			return
		}
		var parameters map[string]interface{}
		if v, ok := data["parameters"]; ok {
			parameters = v.(map[string]interface{})
		}
		localNode := server.DB.ReadLocalMerkabaNode()
		/*如果网站有验证的脚本，编译后使用*/
		names := strings.Split(scriptUri, "/")
		if len(names) == 0 {
			server.writeResponse(c, errorResp("uri can't find siteName"))
			return
		}
		siteName := names[0]
		/*读取脚本内容和版本*/
		var builder strings.Builder
		name := siteName + "/timeout"
		script, _ := server.DB.ReadScriptByUri(name)
		if len(script) > 0 {
			builder.WriteString(script)
		}
		var scriptVersion string
		script, scriptVersion = server.DB.ReadScript(scriptId)
		builder.WriteString(script)
		scriptContent := builder.String()
		instance := server.buildScriptInstance(c, localNode, siteName, scriptId, scriptUri, scriptVersion, scriptContent, parameters, taskName)
		if instance == nil {
			return
		}
		server.DB.UseMemInstance(instance)
		if vv, ok := data["breakPoints"]; ok && len(vv.([]any)) > 0 {
			instance.BreakPoints = make([]map[string]any, 0)
			for _, v := range vv.([]any) {
				instance.BreakPoints = append(instance.BreakPoints, v.(map[string]any))
			}
		}
		if v, ok := data["variables"]; ok && len(v.([]any)) > 0 {
			instance.Variables = common.CopyArray[string](v.([]string))
		}
		if v, ok := data["cookieId"]; ok {
			instance.Context.CookieId = v.(string)
		}
		if v, ok := data["appServerIP"]; ok {
			instance.Context.AppServerIP = v.(string)
		}
		if v, ok := data["appServerPort"]; ok {
			instance.Context.AppServerPort = v.(string)
		}
		if v, ok := data["maxWaitTime"]; ok {
			instance.Context.MaxWaitTime = common.ParseInt64(v)
		}
		if instance.Context.MaxWaitTime == 0 {
			instance.Context.MaxWaitTime = 10
		}
		instance.Context.RunMode = common.RunModeAppServer
		if v, ok := data["runMode"]; ok {
			instance.Context.RunMode = common.ParseRunMode(v.(string))
		}

		Queue.Enqueue(instance)
		json := successResp()
		json["scriptId"] = scriptId
		json["scriptUri"] = scriptUri
		json["scriptVersion"] = scriptVersion
		json["taskName"] = taskName
		json["ip"] = common.LocalIP
		server.writeResponse(c, json)
	})
}
