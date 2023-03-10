package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"merkaba/common"
	"strconv"
	"strings"
	"time"
)

type PaasClient struct {
	serverIP    string
	instance    GrpcAgentServiceClient
	connect     *GrpcConnect
	cookieId    string
	serviceName string
}

func UsePaasClient(serverIP string, serviceName string, cookieId string) *PaasClient {
	if !common.Env.Environment.Production && len(common.Env.Debug.AppServer) > 0 {
		serverIP = common.Env.Debug.AppServer
	}
	conn := getGrpcConnect(serverIP)
	result := PaasClient{
		serverIP:    serverIP,
		connect:     conn,
		instance:    NewGrpcAgentServiceClient(conn.Instance),
		cookieId:    cookieId,
		serviceName: serviceName,
	}
	return &result
}

func UsePaasClient2(name string, cookieId string) *PaasClient {
	serverNode := common.Consul.ReadServiceNodes(name)
	if serverNode == nil {
		common.LoggerStd.Error(name + "没有发现服务节点")
		return nil
	}
	ids := strings.Split(serverNode.ID, "@")
	serviceName := ids[0]
	var serverIP string
	if !common.Env.Environment.Production && len(common.Env.Debug.AppServer) > 0 {
		serverIP = common.Env.Debug.AppServer
	} else {
		serverIP = ids[1]
	}
	conn := getGrpcConnect(serverIP)
	result := PaasClient{
		serverIP:    serverIP,
		connect:     conn,
		instance:    NewGrpcAgentServiceClient(conn.Instance),
		cookieId:    cookieId,
		serviceName: serviceName,
	}
	return &result
}

func (client *PaasClient) Dispose() {
	retGrpcConnect(client.connect)
}

func (client *PaasClient) Invoke(methodName string, params map[string]any, stream []byte) (any, error) {
	return client.DoInvoke(methodName, params, stream, 2*time.Minute)
}

func (client *PaasClient) DoInvoke(methodName string, params map[string]any, stream []byte, timeOut time.Duration) (any, error) {
	if client.connect == nil {
		return nil, errors.New(client.serverIP + " connect is null")
	}
	var traceId string
	if v, ok := params["traceId"]; ok {
		traceId = v.(string)
	} else {
		traceId = ""
	}
	ctx := &AgentContext{
		CookieId:  client.cookieId,
		Server:    common.Env.Environment.LocalName,
		UserId:    common.Env.Environment.LocalIP,
		UserName:  "",
		TraceId:   traceId,
		Timestamp: time.Now().UnixNano(),
	}
	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	request := &AgentRequest{
		Ctx:         ctx,
		LangType:    2,
		Module:      "",
		ServiceName: client.serviceName,
		MethodName:  methodName,
		Params:      bytes,
		Stream:      stream,
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	resp, err := client.instance.Invoke(timeOutCtx, request)
	if err != nil {
		msg := "grpc remote call[" + client.serviceName + "." + methodName + "]&AppServer=" + client.serverIP
		scriptUri := params["scriptUri"]
		if scriptUri == nil {
			common.LoggerStd.Error(msg, zap.Error(err))
		} else {
			taskName := params["taskName"]
			common.LoggerBy(scriptUri.(string)+"@"+taskName.(string)).Error(msg, zap.Error(err))
		}
		return nil, err
	}
	dataType := resp.GetDataType()
	if !resp.IsSuccess {
		if dataType == 0 {
			return nil, nil
		} else {
			return nil, errors.New(resp.DataBody)
		}
	}
	body := resp.GetDataBody()
	switch dataType {
	case 0: //void
		return nil, nil
	case 1: //string
		return body, nil
	case 2: //integer
		i, e := strconv.Atoi(body)
		return i, e
	case 3: //long
		i, e := strconv.ParseInt(body, 10, 64)
		return i, e
	case 4: //double
		i, e := strconv.ParseFloat(body, 64)
		return i, e
	case 5: //float
		i, e := strconv.ParseFloat(body, 32)
		return i, e
	case 6: //boolean
		return body == "true", nil
	case 16: //list
		result := new([]interface{})
		err = json.Unmarshal([]byte(body), &result)
		return result, err
	default:
		result := make(map[string]any)
		err = json.Unmarshal([]byte(body), &result)
		return result, err
	}
}
