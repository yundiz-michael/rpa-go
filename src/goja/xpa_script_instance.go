package goja

import (
	"go.uber.org/zap"
	"merkaba/chromedp"
	"merkaba/common"
	"merkaba/rpc"
	"time"
)

type ScriptInstance struct {
	Context       *common.RunContext
	ScriptContent string
	StartTime     int64
	StopTime      int64
	IsSuccess     bool
	Variables     []string
	BreakPoints   []map[string]interface{}
	Status        string
	ErrorMessage  string
	DB            *ScriptDb
	RunVM         *ScriptVM
	fnTimeout     func(any) (bool, error)
}

// HandleTimeout  如果返回为true,则退出timeout函数，否则一直判断
func (s *ScriptInstance) HandleTimeout(webPage any, err error) (bool, error) {
	vm := s.RunVM.Runtime
	if s.fnTimeout == nil {
		return true, err
	}
	page := vm.CreateWebPageObject(webPage.(*chromedp.WebPage))
	result, _err := s.fnTimeout(page)
	if _err != nil {
		return false, _err
	}
	return result, err
}

func (s *ScriptInstance) InitVM() error {
	s.RunVM = &ScriptVM{}
	s.Context.OnClientClose = s.OnClientClose
	return s.RunVM.Init(s)
}

func (s *ScriptInstance) OnClientClose(ev interface{}) {
	if client, ok := ev.(*chromedp.WebClient); ok {
		delete(WebClients, client.Option.Domain)
		s.DB.RemoveMemInstance(s)
	}
}

func (s *ScriptInstance) Run() {
	fields := make([]zap.Field, 0)
	fields = append(fields, zap.String("taskName", s.Context.TaskName))
	fields = append(fields, zap.String("version", s.Context.ScriptVersion))
	fields = append(fields, zap.String("uri", s.Context.ScriptUri))

	s.DB.UseMemInstance(s)
	/*数据发回监控中心，更新实例状态*/
	s.StopTime = time.Now().UnixMilli()
	vm := s.RunVM.Runtime
	msg := "===>start script"
	s.Context.Info(msg, fields...)
	s.remoteCall("Merkaba", "onStartScript")
	if s.Context.RunMode == common.RunModeBrowserRun {
		if len(s.BreakPoints) > 0 {
			debugger := s.RunVM.Runtime.AttachDebugger()
			for _, breakPoint := range s.BreakPoints {
				_scriptUri := breakPoint["scriptUri"].(string)
				_lines := breakPoint["lines"].([]any)
				for _, _line := range _lines {
					debugger.SetBreakpoint(_scriptUri, common.ParseInt(_line))
				}
			}
		}
	}
	program, err := vm.Compile(s.Context.ScriptUri, s.ScriptContent, false, true, nil)
	if err == nil {
		if v, ok := s.Context.Parameters["enableVNC"]; s.Context.RunMode == common.RunModeBrowserRun && ok && v.(bool) {
			common.StartVNC(common.VncWidth, common.VncHeight, s.Context.TaskName)
			common.LoggerStd.Info("启动VNC", zap.String("任务名称", s.Context.TaskName))
		}
		fnValue := vm.Get("handleTimeout")
		if fnValue != nil {
			vm.ExportTo(fnValue, &s.fnTimeout)
		}
		_, err = vm.RunProgram(program)
	}
	s.RunVM.Clear()
	s.StopTime = time.Now().UnixMilli()
	s.IsSuccess = true
	if err != nil {
		s.RunVM.Snapshot(common.Env.Path.Shot + "/" + s.Context.TaskName)
		s.ErrorMessage = err.Error()
		common.SendMessage("error", s.Context, s.ErrorMessage)
		s.Context.Error(err.Error(), fields...)
		common.LoggerStd.Error(err.Error(), fields...)
		s.IsSuccess = false
	}
	s.DB.FreeMemInstance(s)
	common.SendMessage("stop", s.Context, "运行结束")
	/*数据发回监控中心，更新实例状态*/
	s.remoteCall("Merkaba", "onStopScript")
	msg = "===>stop script"
	s.Context.Info(msg, fields...)
}

func (s *ScriptInstance) remoteCall(service string, funcName string) {
	param := s.AsMap()
	client := rpc.UsePaasClient(s.Context.AppServerIP, service, s.Context.CookieId)
	if client == nil {
		return
	}
	result, err := client.Invoke(funcName, param, make([]byte, 0))
	client.Dispose()
	if err != nil {
		s.Context.Error(err.Error(), zap.String("ip", s.Context.AppServerIP), zap.Error(err))
		common.LoggerStd.Error(err.Error(), zap.String("ip", s.Context.AppServerIP), zap.Error(err))
		return
	}
	common.LoggerStd.Info(service+"."+funcName, zap.String("result", result.(string)), zap.String("scriptUri", s.Context.ScriptUri))
}

func (s *ScriptInstance) Interrupt() {
	info := "用户终止"
	common.SendMessage("stop", s.Context, info)
	s.RunVM.Runtime.Interrupt(info)
	s.DB.FreeMemInstance(s)
	s.RunVM.CloseWebClients()
}

func (s *ScriptInstance) AsMap() map[string]any {
	data := make(map[string]any)
	data["siteName"] = s.Context.SiteName
	data["scriptId"] = s.Context.ScriptId
	data["taskName"] = s.Context.TaskName
	data["scriptUri"] = s.Context.ScriptUri
	data["cookieId"] = s.Context.CookieId
	data["scriptVersion"] = s.Context.ScriptVersion
	data["parameters"] = s.Context.Parameters
	data["isSuccess"] = s.IsSuccess
	data["ip"] = common.LocalIP
	data["server"] = common.LocalName
	data["startTime"] = s.StartTime
	data["hasVNC"] = common.HasVNC(s.Context.TaskName)
	data["stopTime"] = s.StopTime
	data["error"] = s.ErrorMessage
	if s.Status == "Running" {
		data["status"] = 1
	} else {
		data["status"] = 0
	}
	return data
}
