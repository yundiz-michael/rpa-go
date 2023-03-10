package common

import "go.uber.org/zap"

type RunContext struct {
	RunMode       RunMode
	AppServerIP   string
	AppServerPort string
	CookieId      string
	SiteName      string
	TaskName      string /*实例的名称*/
	ScriptId      string
	ScriptUri     string
	ScriptVersion string
	MaxWaitTime   int64
	Headless      bool
	Proxy         bool
	OnClientClose func(ev interface{})
	Parameters    map[string]interface{}
	loggerTask    *zap.Logger
}

func (ctx *RunContext) Init(parameters map[string]any) {
	ctx.Parameters = parameters
	if ctx.loggerTask == nil {
		loggerCore := LoggerCoreBy(ctx.TaskName)
		ctx.loggerTask = zap.New(loggerCore, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())
	}
}

func (ctx *RunContext) logger() *zap.Logger {
	if ctx.RunMode == 0 {
		return LoggerStd
	} else {
		return ctx.loggerTask
	}
}

func (ctx *RunContext) Info(msg string, fields ...zap.Field) {
	ctx.logger().Info(msg, fields...)
}

func (ctx *RunContext) Error(msg string, fields ...zap.Field) {
	ctx.logger().Error(msg, fields...)
}

func (ctx *RunContext) Warn(msg string, fields ...zap.Field) {
	ctx.logger().Warn(msg, fields...)
}

func (ctx *RunContext) AsMap() map[string]any {
	data := make(map[string]any)
	data["siteName"] = ctx.SiteName
	data["scriptId"] = ctx.ScriptId
	data["taskName"] = ctx.TaskName
	data["scriptUri"] = ctx.ScriptUri
	data["scriptVersion"] = ctx.ScriptVersion
	return data
}

func (ctx *RunContext) Account() string {
	if v, ok := ctx.Parameters["userName"]; ok {
		return v.(string)
	} else {
		return ctx.TaskName
	}
}
