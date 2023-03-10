package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"go.uber.org/zap"
	"time"
)

type PulsarClient struct {
	instance pulsar.Client
	producer pulsar.Producer
}

type PulsarConfig struct {
	ID    string   `json:"id"`
	Hosts []string `json:"hosts"`
}

func (p *PulsarClient) Init(topic string, host string) {
	var err error
	p.instance, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:               fmt.Sprintf("pulsar://%s", host),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		LoggerStd.Panic("init", zap.NamedError("error", err))
	}
	properties := make(map[string]string)
	p.producer, err = p.instance.CreateProducer(pulsar.ProducerOptions{
		Schema: pulsar.NewStringSchema(properties),
		Topic:  topic,
		Name:   LocalName,
	})
	if err != nil {
		LoggerStd.Panic("init topic", zap.NamedError("error", err))
	} else {
		LoggerStd.Info("init topic", zap.String("name", topic))
	}
}

func SendMessage(level string, ctx *RunContext, message string) {
	for _, client := range PulsarMessageClients {
		client._SendMessage(level, ctx.CookieId, ctx.TaskName, ctx.ScriptVersion, message)
	}
}

func (p *PulsarClient) _SendMessage(level string, _cookieId string, _taskName string, _scriptVersion string, message string) {
	info := make(map[string]any)
	info["type"] = "message"
	info["server"] = LocalName
	info["ip"] = LocalIP
	info["level"] = level
	info["cookieId"] = _cookieId
	info["taskName"] = _taskName
	info["scriptVersion"] = _scriptVersion
	info["message"] = message
	info["timestamp"] = time.Now().UnixMilli()
	bytes, _ := json.Marshal(info)
	_, err := p.producer.Send(context.Background(), &pulsar.ProducerMessage{
		Value: string(bytes),
	})
	if err != nil {
		LoggerStd.Error("SendMessage", zap.NamedError("error", err))
	}
}

func SendLog(level string, ctx *RunContext, info map[string]any) {
	for _, client := range PulsarMessageClients {
		client._SendLog(level, ctx, info)
	}
}

func (p *PulsarClient) _SendLog(level string, ctx *RunContext, info map[string]any) {
	info["type"] = "message"
	info["cookieId"] = ctx.CookieId
	info["taskName"] = ctx.TaskName
	info["server"] = LocalName
	info["ip"] = LocalIP
	info["level"] = level
	info["timestamp"] = time.Now().UnixMilli()
	bytes, _ := json.Marshal(info)
	_, err := p.producer.Send(context.Background(), &pulsar.ProducerMessage{
		Value: string(bytes),
	})
	if err != nil {
		LoggerStd.Error("SendLog", zap.NamedError("error", err))
	}
}

func Notify(action string, info map[string]any) error {
	info["type"] = "notify"
	info["action"] = action
	info["server"] = LocalName
	info["ip"] = LocalIP
	info["timestamp"] = time.Now().UnixMilli()
	bytes, err := json.Marshal(info)
	if err != nil {
		LoggerStd.Error("notify", zap.NamedError("error", err))
		return err
	}
	for _, client := range PulsarMessageClients {
		_, err = client.producer.Send(context.Background(), &pulsar.ProducerMessage{
			Value: string(bytes),
		})
		if err != nil {
			LoggerStd.Error("notify", zap.NamedError("error", err))
			return err
		}
	}
	return nil
}

func SendData(_type string, _format string, ctx *RunContext, values map[string]any) error {
	for _, client := range PulsarDataClients {
		err := client._SendData(_type, _format, ctx, values)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PulsarClient) _SendData(_type string, _format string, ctx *RunContext, values map[string]any) error {
	data := make(map[string]any)
	data["type"] = _type
	data["format"] = _format
	data["server"] = LocalName
	data["ip"] = LocalIP
	data["cookieId"] = ctx.CookieId
	data["taskName"] = ctx.TaskName
	data["data"] = values
	data["timestamp"] = time.Now().UnixMilli()
	bytes, err := json.Marshal(data)
	if err != nil {
		LoggerStd.Error("SendData", zap.NamedError("error", err))
		return err
	}
	_, err = p.producer.Send(context.Background(), &pulsar.ProducerMessage{
		Value: string(bytes),
	})
	if err != nil {
		LoggerStd.Error("SendData", zap.NamedError("error", err))
		return err
	}
	return nil
}
