package goja

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"merkaba/common"
)

type ScriptVM struct {
	Context       *common.RunContext
	Command       chan DebugCommand
	Runtime       *Runtime
	RequireModule *RequireModule
	Registry      *Registry
}

func (v *ScriptVM) Init(instance *ScriptInstance) error {
	ctx := instance.Context
	v.Context = ctx
	v.Runtime = New()
	v.Runtime.ScriptHandler = instance
	v.Runtime.Context = ctx
	v.Command = make(chan DebugCommand, 1)
	printer := PrinterFunc(func(level string, s string) {
		isBrowser := ctx.RunMode == common.RunModeBrowserRun
		if isBrowser {
			common.SendMessage(level, ctx, s)
		}
		switch level {
		case "warn":
			ctx.Warn(s)
			break
		case "error":
			ctx.Error(s)
			break
		default:
			ctx.Info(s)
			break
		}
	})

	InitConsole(v.Runtime, printer)
	v.Registry = NewRegistry(WithGlobalFolders("."), WithLoader(func(url string) ([]byte, error) {
		content, _ := instance.DB.ReadScriptByUri(url[1:])
		if len(content) > 0 {
			return []byte(content), nil
		} else {
			return nil, errors.New(fmt.Sprintf("require [%s] does not exist", url))
		}
	}))
	v.RequireModule = v.Registry.Enable(v.Runtime)
	return nil
}

// Clear todo michael xiao 需要优化，仅仅删除更新的脚本
func (v *ScriptVM) Clear() {
	v.Runtime.ClearInterrupt()
	v.RequireModule.clear()
	v.Registry.clear()
}
func (v *ScriptVM) Snapshot(fileName string) {
	if v.Runtime.WebClient != nil {
		v.Runtime.WebClient.Snapshot(fileName + ".jpg")
	}
}

func (v *ScriptVM) CloseWebClients() {
	if v.Runtime.WebClient == nil {
		return
	}
	common.LoggerStd.Info("Close WebClient", zap.String("taskName", v.Context.TaskName))
	RemoveWebClient(v.Runtime.WebClient.Option.Domain, v.Context.TaskName)
	v.Runtime.WebClient.Close()
}
