package main

import (
	"go.uber.org/zap"
	"merkaba/common"
	"merkaba/common/queue"
	"merkaba/goja"
	"merkaba/server"
	"runtime"
)

var db goja.ScriptDb

func main() {
	common.InitEnviroment()
	db = goja.ScriptDb{Client: common.Mysql}
	db.Init()
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		db.UnRegisterMerkabaNode()
		switch err.(type) {
		case runtime.Error:
			common.LoggerStd.Error("merkaba crash", zap.Error(err.(error)))
		}
	}()
	server.Queue = queue.NewQueue(8)
	server.Queue.Start()
	defer server.Queue.Stop()
	common.Consul.RegisterMerkaba()
	for i := 0; i < 100; i++ {
		db.RegisterMerkabaNode()
	}
	server.NewHttpServer(db).Start()
}
