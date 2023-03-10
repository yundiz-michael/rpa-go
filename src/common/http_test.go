package common

import (
	"fmt"
	"testing"
)

func TestHttpPost(t *testing.T) {
	InitEnviroment()
	resp, err := HttpGet("http://193.168.1.30:8500/v1/kv/beijing.xpa.sqldb.xpa", nil)
	fmt.Println(err)
	fmt.Println(resp)
}

func TestHttpOCRImage(t *testing.T) {
	InitEnviroment()
	values := make(map[string]any)
	resp, err := HttpOCRImage(values)
	fmt.Println(err)
	fmt.Println(resp)
}
