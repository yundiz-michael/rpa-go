package common

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type YamlFile struct {
	Environment struct {
		DataCenter string `yaml:"dataCenter"`
		Cluster    string `yaml:"cluster"`
		Domain     string `yaml:"domain"`
		Production bool   `yaml:"production"`
		LocalIP    string `yaml:"localIP"`
		LocalName  string `yaml:"localName"`
	}
	Pulsar struct {
		Topic string   `yaml:"topic"`
		Hosts []string `yaml:"hosts"`
	}
	Vnc struct {
		AutoClose bool `yaml:"autoClose"`
	}
	Database struct {
		Mysql string `yaml:"mysql"`
		Mongo string `yaml:"mongo"`
	}
	Debug struct {
		AppServer string `yaml:"appServer"`
	}
	Path struct {
		Config string `yaml:"config"`
		Temp   string `yaml:"temp"`
		Python string `yaml:"python"`
		Shot   string `yaml:"shot"`
		Logger string `yaml:"logger"`
	}
	Consul struct {
		Development []string `yaml:"development"`
		Production  []string `yaml:"production"`
	}
	Version string `yaml:"version"`
}

func LoadYamlConfig() (file *YamlFile) {
	RootPath = os.Getenv("MerkabaRoot")
	var path string
	if len(RootPath) == 0 {
		RootPath = "/workspace/xpa/go/merkaba/"
		path = RootPath + "server.yaml"
	} else {
		path = RootPath + "conf/server.yaml"
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(f, &file)
	if err != nil {
		panic(err.Error())
	}
	return file
}
