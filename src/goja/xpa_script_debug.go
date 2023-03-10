package goja

import (
	"merkaba/common"
)

type DebugCommand struct {
	InstanceId string
	TaskName   string
	Command    string
	Line       int
}

func (dc *DebugCommand) printVariables(t *ScriptInstance) {
	values := make(map[string]any)
	for _, name := range t.Variables {
		values[name] = t.RunVM.Runtime.Get(name).String()
	}
	common.SendData("watch", "map", t.Context, values)
}

//func (dc *DebugCommand) Execute(t *ScriptInstance, dbg *Debugger) {
//switch dc.Command {
//case "addBreakPoint":
//	err := dbg.SetBreakpoint(dc.TaskName, dc.Line)
//	if err != nil {
//		common.Pulsar.SendMessage("info", t.cookieId, err.Error())
//	}
//
//case "removeBreakPoint":
//	err := dbg.ClearBreakpoint(dc.TaskName, dc.Line)
//	if err != nil {
//		common.Pulsar.SendMessage("info", t.cookieId, err.Error())
//	}
//
//case "next":
//	err := dbg.Next()
//	if err != nil {
//		common.Pulsar.SendMessage("info", t.cookieId, err.Error())
//	} else {
//		dc.printVariables(t)
//	}
//
//case "continue":
//	dbg.Continue()
//	dc.printVariables(t)
//
//case "step":
//	err := dbg.StepIn()
//	if err != nil {
//		common.Pulsar.SendMessage("info", t.cookieId, err.Error())
//	} else {
//		dc.printVariables(t)
//	}
//
//}
//}

func (dc *DebugCommand) IsDone() bool {
	return dc.Command == "quit"
}
