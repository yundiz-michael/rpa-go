package common

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"os"
	"strings"
	"sync"
)

type RunMode int

const (
	RunModeBrowserRun RunMode = 0
	RunModeAppServer  RunMode = 1
	RunModeNative     RunMode = 2
)

type ServerNode struct {
	ID   string
	Port int
}

type ScriptHandler = interface {
	HandleTimeout(webPage any, err error) (bool, error)
}

func ParseRunMode(value string) RunMode {
	switch value {
	case "BrowserRun":
		return RunModeBrowserRun
	case "AppServer":
		return RunModeAppServer
	case "Native":
		return RunModeNative
	default:
		return RunModeAppServer
	}
}

var RootPath string
var Env *YamlFile
var LoggerStd *zap.Logger
var LoggerCoreMap map[string]zapcore.Core
var LoggerMap map[string]*zap.Logger
var LoggerStdCore zapcore.Core
var LoggerCoreLock sync.Mutex
var LoggerLock sync.Mutex
var Consul *ConsulClient
var ConsulClientPool = make(map[string]sync.Pool)
var ConsulClientLock = sync.RWMutex{}
var PulsarDataClients []*PulsarClient
var PulsarMessageClients []*PulsarClient
var Mysql *MysqlClient
var Hazelcast *HazelcastClient
var ConfigPath string
var _isLoaded = false
var LocalName string
var LocalIP string
var LocalPort = 4320

func initLogger() {
	LoggerCoreMap = make(map[string]zapcore.Core)
	LoggerMap = make(map[string]*zap.Logger)
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)
	LoggerStdCore = zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewCustomEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		atom,
	)
	LoggerStd = zap.New(LoggerStdCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	// defer logger.Sync()
	// 配置 zap 包的全局变量
	zap.ReplaceGlobals(LoggerStd)
	// 运行时安全地更改 logger 日记级别
	atom.SetLevel(zap.InfoLevel)
}

func LoggerCoreBy(key string) zapcore.Core {
	LoggerCoreLock.Lock()
	defer LoggerCoreLock.Unlock()
	result := LoggerCoreMap[key]
	if result != nil {
		return result
	}
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)
	fileName := Env.Path.Logger + strings.Replace(key, "/", "-", -1) + ".log"
	var zapFile = zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10240, // megabytes
		MaxBackups: 10,
		MaxAge:     7, // days
	})
	result = zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewCustomEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapFile),
		atom,
	)
	LoggerCoreMap[key] = result
	return result
}

func LoggerBy(key string) *zap.Logger {
	LoggerLock.Lock()
	defer LoggerLock.Unlock()
	result := LoggerMap[key]
	if result != nil {
		return result
	}
	uriLoggerCore := LoggerCoreMap[key]
	if uriLoggerCore == nil {
		uriLoggerCore = LoggerCoreBy(key)
		LoggerCoreMap[key] = uriLoggerCore
	}
	result = zap.New(uriLoggerCore, zap.AddCaller(), zap.Development())
	LoggerMap[key] = result
	return result
}

func addPulsarDataClient(topic string, host string) {
	pulsarClient := &PulsarClient{}
	pulsarClient.Init(topic+"@data", host)
	PulsarDataClients = append(PulsarDataClients, pulsarClient)
}

func addPulsarMessageClient(topic string, host string) {
	pulsarClient := &PulsarClient{}
	pulsarClient.Init(topic+"@message", host)
	PulsarMessageClients = append(PulsarMessageClients, pulsarClient)
}

func InitEnviroment() {
	if _isLoaded {
		return
	}
	_isLoaded = true
	Env = LoadYamlConfig()
	/*初始化日志*/
	initLogger()
	if len(Env.Environment.LocalName) == 0 {
		LocalName, _ = os.Hostname()
	} else {
		LocalName = Env.Environment.LocalName
	}
	if Env.Environment.Production {
		LoggerStd.Info("Production Mode", zap.String("ip", Env.Consul.Production[0]))
	} else {
		LoggerStd.Info("Develop Mode", zap.String("ip", Env.Consul.Development[0]))
	}
	if len(Env.Environment.LocalIP) == 0 {
		LocalIP, _ = getClientIp()
	} else {
		LocalIP = Env.Environment.LocalIP
	}
	/*初始化consul*/
	ConfigPath = Env.Path.Config
	Consul = &ConsulClient{Env: Env}
	Consul.Init()
	/*初始化mysql*/
	Mysql = &MysqlClient{Env: Env}
	Mysql.Init()

	topic := "merkaba"
	if len(Env.Pulsar.Topic) > 0 {
		topic = Env.Pulsar.Topic
	}
	/*构建消息队列客户端数组，因为一个数据可能要发给多个集群*/
	PulsarDataClients = make([]*PulsarClient, 0)
	PulsarMessageClients = make([]*PulsarClient, 0)
	/*读取当前集群本身的pulsar*/
	data := Consul.readConsulConfig("pulsar", "default")
	var config = PulsarConfig{}
	json.Unmarshal(data.Value, &config)
	addPulsarDataClient(topic, config.Hosts[0])
	addPulsarMessageClient(topic, config.Hosts[0])
	/*读取本地的server.yaml的配置*/
	for _, host := range Env.Pulsar.Hosts {
		if config.Hosts[0] == host {
			continue
		}
		addPulsarDataClient(topic, host)
		addPulsarMessageClient(topic, host)
	}
}

func InitHazelClient() {
	Hazelcast = &HazelcastClient{Env: Env}
	Hazelcast.Init()
}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		LoggerStd.Error(err.Error())
		return "127.0.0.1", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "127.0.0.1", nil
}
