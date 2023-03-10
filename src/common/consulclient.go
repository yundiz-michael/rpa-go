package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"merkaba/common/round"
	"time"
)

type ConsulClient struct {
	Env      *YamlFile
	instance *api.Client
	config   *api.Config
}

type SqlAccount struct {
	Encrypt  bool   `json:"encrypt"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SqlDBConfig struct {
	ID      string      `json:"id"`
	DBType  string      `json:"type"`
	Hosts   []string    `json:"hosts"`
	Account *SqlAccount `json:"account"`
}

func (c *ConsulClient) Init() {
	if c.Env.Environment.Production {
		c.config = &api.Config{Address: c.Env.Consul.Production[0], Datacenter: c.Env.Environment.DataCenter}
	} else {
		c.config = &api.Config{Address: c.Env.Consul.Development[0], Datacenter: c.Env.Environment.DataCenter}
	}
	var err error
	c.instance, err = api.NewClient(c.config)
	if err != nil {
		LoggerStd.Panic("consul client error ", zap.String("address", c.config.Address))
	}
}

func (c *ConsulClient) readConsulConfig(nodeType string, nodeName string) *api.KVPair {
	var buffer bytes.Buffer
	buffer.WriteString(c.Env.Environment.DataCenter)
	buffer.WriteString(".")
	buffer.WriteString(c.Env.Environment.Cluster)
	buffer.WriteString(".")
	buffer.WriteString(nodeType)
	buffer.WriteString(".")
	buffer.WriteString(nodeName)
	key := buffer.String()
	data, _, err := c.instance.KV().Get(key, nil)
	if err != nil {
		LoggerStd.Panic("consul read Config", zap.String("key", key), zap.NamedError("error", err))
		return nil
	}
	return data
}

func (c *ConsulClient) ReadDbConfig(name string) *SqlDBConfig {
	data := c.readConsulConfig("sqldb", name)
	var result = SqlDBConfig{}
	json.Unmarshal(data.Value, &result)
	return &result
}

func (c *ConsulClient) RegisterMerkaba() {
	LoggerStd.Info("Register merkaba node ", zap.String("localIP", LocalIP), zap.Int("localPort", LocalPort))
	registration := new(api.AgentServiceRegistration)
	registration.ID = LocalIP
	registration.Name = "merkaba"
	registration.Address = LocalIP
	registration.Port = LocalPort
	registration.Tags = []string{LocalName, LocalIP, time.Now().Format("01-02 15:04:05")}
	/* 增加consul健康检查回调函数 */
	check := new(api.AgentServiceCheck)
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Interval = "5s"
	check.Timeout = "10s"
	check.DeregisterCriticalServiceAfter = "10s"
	registration.Check = check
	c.instance.Agent().ServiceRegister(registration)
}

func (c *ConsulClient) NextPlatoServer() string {
	catalogServices, _, _ := c.instance.Catalog().Service("plato", "", nil)
	size := len(catalogServices)
	if size == 0 {
		return ""
	}
	index := NewHash().Seed(fmt.Sprintf("%d", time.Now().Unix())).DeHash() % size
	service := catalogServices[index]
	return service.ServiceAddress + ":" + fmt.Sprintf("%d", service.ServicePort)
}

func (c *ConsulClient) ReadServiceNodes(name string) *ServerNode {
	name = fmt.Sprintf("%s.%s.%s", Env.Environment.DataCenter, Env.Environment.Cluster, name)
	catalogServices, _, err := c.instance.Catalog().Service(name, "", nil)
	if err != nil || len(catalogServices) == 0 {
		return nil
	}
	nodeRound, _ := round.BuildRound(catalogServices)
	var service *api.CatalogService
	service = nodeRound.Next().Service
	return &ServerNode{ID: service.ServiceID, Port: service.ServicePort}
}
