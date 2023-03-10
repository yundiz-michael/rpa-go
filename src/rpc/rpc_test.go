package rpc

import (
	"fmt"
	"io"
	"merkaba/common"
	"os"
	"testing"
)

func TestPaasService(t *testing.T) {
	common.InitEnviroment()
	client := UsePaasClient("196.168.1.62", "Merkaba", "xiaoming")
	params := make(map[string]any)
	params["name"] = "michael"
	result, err := client.Invoke("onStartScript", params, make([]byte, 0))
	client.Dispose()
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		fmt.Printf("%s ===>1", result)
	}
}

func TestPythonService_Image_Eye(t *testing.T) {
	common.InitEnviroment()
	file, _ := os.Open("/workspace/temp/8.jpg")
	bytes, err := io.ReadAll(file)
	params := make(map[string]any)
	params["detectText"] = true
	result, err := See("196.168.1.3", params, bytes)
	if err == nil && result["isSuccess"].(bool) {
		fmt.Print("识别文本成功成功")
	}
}
