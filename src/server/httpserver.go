package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"merkaba/common"
	"merkaba/common/queue"
	"merkaba/goja"
	"net/http"
	"time"
)

var Queue *queue.Queue

type HttpServer struct {
	instance *gin.Engine
	DB       goja.ScriptDb
}

func NewHttpServer(db goja.ScriptDb) *HttpServer {
	r := gin.New()
	r.Use(gin.Recovery())
	result := &HttpServer{
		instance: r,
		DB:       db,
	}
	return result
}

func successResp() gin.H {
	return gin.H{
		"serverName": common.LocalName,
		"serverIP":   common.LocalIP,
		"isSuccess":  true,
	}
}

func errorResp(message string) gin.H {
	return gin.H{
		"serverName": common.LocalName,
		"serverIP":   common.LocalIP,
		"isSuccess":  false,
		"message":    message,
	}
}

func (server *HttpServer) parseData(c *gin.Context) map[string]interface{} {
	var data map[string]interface{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		common.LoggerStd.Error(err.Error())
		return nil
	}
	return data
}

func (server *HttpServer) buildScriptInstance(c *gin.Context, localNode *goja.MerkabaNode, siteName string,
	scriptId string, scriptUri string, scriptVersion string, scriptContent string, parameters map[string]any,
	taskName string) *goja.ScriptInstance {
	instance := server.DB.FindMemInstance(taskName)
	if instance == nil {
		common.LoggerStd.Info("创建新的VM", zap.String("taskName", taskName))
		runCount, idleCount := server.DB.ReadInstanceCount()
		if localNode.MaxCount > 0 && (runCount+idleCount) > localNode.MaxCount {
			server.writeResponse(c, errorResp(fmt.Sprintf("can't lanuch new instance,exceed %d", localNode.MaxCount)))
			return nil
		}
		context := &common.RunContext{
			TaskName:      taskName,
			ScriptId:      scriptId,
			ScriptUri:     scriptUri,
			ScriptVersion: scriptVersion,
			SiteName:      siteName,
		}
		instance = &goja.ScriptInstance{
			Context:   context,
			DB:        &server.DB,
			StartTime: time.Now().UnixMilli(),
		}
		err := instance.InitVM()
		if err != nil {
			server.writeResponse(c, errorResp(err.Error()))
			return nil
		}
		server.DB.AddMemInstance(instance)
	} else {
		if instance.Status == "Running" {
			server.writeResponse(c, errorResp("实例正在运行中，请先终止运行"))
			return nil
		} else {
			common.LoggerStd.Info("使用缓存VM", zap.String("taskName", taskName))
		}
	}
	/*更新实列的运行参数*/
	instance.Context.Init(parameters)
	instance.ScriptContent = scriptContent
	if v, ok := parameters["proxy"]; ok {
		instance.Context.Proxy = v.(bool)
	}
	if v, ok := parameters["headless"]; ok {
		instance.Context.Headless = v.(bool)
	}
	return instance
}

func (server *HttpServer) writeResponse(c *gin.Context, json gin.H) {
	c.JSON(http.StatusOK, json)
}

func (server *HttpServer) Start() {
	server.registerRunScript()
	server.registerDebug()
	server.registerInfo()
	server.registerWatch()
	server.registerStopScript()
	server.registerReadScriptCount()
	server.registerReadScriptInstance()
	server.registerVideo()
	server.registerStartVNC()
	server.registerStopVNC()
	common.LoggerStd.Info("🍊🍊🍊Merkaba start success", zap.String("version", "1.3.16"))
	server.instance.Run(fmt.Sprintf("0.0.0.0:%d", common.LocalPort))
}
