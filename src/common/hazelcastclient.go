package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/serialization"
	"go.uber.org/zap"
	"time"
)

type HazelcastClient struct {
	Env      *YamlFile
	instance *hazelcast.Client
	runtime  string
	merkaba  *hazelcast.Map
}

type HazelcastConfig struct {
	Mode    int      `json:"mode"`
	Members []string `json:"members"`
}

type MerkabaNode struct {
	Name    string
	IP      string
	Runtime string
}

func (n MerkabaNode) String() string {
	return fmt.Sprintf("Name: %s,IP:%s,RunTime:%s", n.Name, n.IP, n.Runtime)
}

func (n MerkabaNode) WriteData(output serialization.DataOutput) {
	output.WriteString(n.Name)
	output.WriteString(n.IP)
	output.WriteString(n.Runtime)
}

func (n *MerkabaNode) ReadData(input serialization.DataInput) {
	n.Name = input.ReadString()
	n.IP = input.ReadString()
	n.Runtime = input.ReadString()
}

func (h *HazelcastClient) Init() {
	ctx := context.Background()
	var err error
	h.runtime = time.Now().Format("01-02 15:04:05")
	c := hazelcast.Config{
		ClientName: LocalName + "_" + h.runtime,
	}
	c.Cluster.Name = h.Env.Environment.DataCenter + "." + h.Env.Environment.Cluster
	data := Consul.readConsulConfig("appserver", "hazelcast")
	var config = HazelcastConfig{}
	json.Unmarshal(data.Value, &config)
	c.Cluster.Network.SetAddresses(config.Members...)
	h.instance, err = hazelcast.StartNewClientWithConfig(ctx, c)
	if err != nil {
		LoggerStd.Panic(err.Error())
	}
	h.merkaba, err = h.instance.GetMap(ctx, "Merkaba")
	if err != nil {
		LoggerStd.Panic(err.Error())
	}
}

func (h *HazelcastClient) RegisterNode() {
	node := &MerkabaNode{Name: LocalName, IP: LocalIP, Runtime: h.runtime}
	b, _ := json.Marshal(node)
	jsonValue := serialization.JSON(b)
	err := h.merkaba.Set(context.Background(), LocalName, jsonValue)
	if err != nil {
		LoggerStd.Info(LocalName, zap.NamedError("error", err))
	}
}
