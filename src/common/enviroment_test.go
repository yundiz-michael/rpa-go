package common

import (
	"fmt"
	"testing"
)

func TestReadMerkabaScript(t *testing.T) {
	InitEnviroment()
}

func TestLoadConfig(t *testing.T) {
	file := LoadYamlConfig()
	fmt.Println(file.Consul.Development[0])
}
