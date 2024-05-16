package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultEnv            = "unknown_env"
	defaultAppID          = "unknown_app_id"
	defaultConsulPort     = 8500
	defaultPrometheusPort = 16670
)

var (
	singletonConfig      *Config
	once4SingletonConfig sync.Once
)

func GetConfig() *Config {
	once4SingletonConfig.Do(initSingletonConfig)
	return singletonConfig
}

type MeterProviderType int

type PushCfg struct {
	PushAddress string        // push url
	PushPeriod  time.Duration // push period
}

type ConsulCfg struct {
	ConsulAddress string // prometheus related config
	ConsulToken   string // prometheus related config
}

// base label names in all time series
type BaseLabelCfg struct {
	Appid       string
	AppVer      string
	Env         string
	IP          string
	DataType    string
	MetricyType string
}

type Config struct {
	PrometheusPort      int // prometheus related config
	Consul              *ConsulCfg
	Push                *PushCfg // push cfg
	MeterProvider       MeterProviderType
	BaseLabel           *BaseLabelCfg
	HistogramBoundaries []float64 // histgoram的分桶配置，默认0.002, 0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 10,
	PrefixBaseLabel     string
	PrefixMetricName    string
	LocalIP             string
	Env                 string
	AppId               string
	AppVer              string
	ReadRuntimeStats    bool
	InfoLogWrite        func(s string)
	ErrorLogWrite       func(s string)
}

func initSingletonConfig() {
	singletonConfig = new(Config)
	singletonConfig.Env = os.Getenv("dtl.env")

	if singletonConfig.Env == "" {
		singletonConfig.Env = defaultEnv
	}
	singletonConfig.AppId = os.Getenv("dtl.app.id")
	if singletonConfig.AppId == "" {
		singletonConfig.AppId = defaultAppID
	}
	var consulPort = defaultConsulPort
	consulPortInt64, err := strconv.ParseInt(os.Getenv("dtl.consul.port"), 10, 64)
	if err == nil {
		consulPort = int(consulPortInt64)
	}

	if os.Getenv("dtl.consul.host") != "" {
		if singletonConfig.Consul == nil {
			singletonConfig.Consul = new(ConsulCfg)
		}
		singletonConfig.Consul.ConsulAddress = fmt.Sprintf("%s:%d", os.Getenv("dtl.consul.host"), consulPort)
	}
	if os.Getenv("HOST_IP") != "" {
		if singletonConfig.Consul == nil {
			singletonConfig.Consul = new(ConsulCfg)
		}
		singletonConfig.Consul.ConsulAddress = fmt.Sprintf("%s:%d", os.Getenv("HOST_IP"), consulPort)
	}

	if os.Getenv("dtl.consul.token") != "" {
		if singletonConfig.Consul != nil {
			singletonConfig.Consul.ConsulToken = os.Getenv("dtl.consul.token")
		}
	}
	if os.Getenv("CONSUL_TOKEN") != "" {
		if singletonConfig.Consul != nil {
			singletonConfig.Consul.ConsulToken = os.Getenv("CONSUL_TOKEN")
		}
	}

	prometheusPort, err := strconv.ParseInt(os.Getenv("dtl.monitor.port"), 10, 64)
	if err != nil {
		singletonConfig.PrometheusPort = defaultPrometheusPort
	} else {
		singletonConfig.PrometheusPort = int(prometheusPort)
	}
	singletonConfig.LocalIP = singletonConfig.getLocalIP()

	singletonConfig.BaseLabel = &BaseLabelCfg{
		Appid:       "appid",
		Env:         "env",
		IP:          "ip",
		DataType:    "data_type",
		MetricyType: "metricy_type",
	}
}

func (cfg *Config) IsTest() bool {
	return strings.Contains(strings.ToLower(cfg.Env), "test")
}

func (cfg *Config) GetEnv() string {
	if cfg.Env == "" {
		return defaultEnv
	}
	return strings.ToLower(cfg.Env)
}

func (cfg *Config) WriteInfoOrNot(s string) {
	if cfg.InfoLogWrite == nil {
		_, _ = os.Stdout.WriteString("[go-metric][info]: " + s + "\n")
	} else {
		cfg.InfoLogWrite("[go-metric] " + s)
	}
}

func (cfg *Config) WriteErrorOrNot(s string) {
	if cfg.ErrorLogWrite == nil {
		_, _ = os.Stdout.WriteString("[go-metric][error]: " + s + "\n")
	} else {
		cfg.ErrorLogWrite("[go-metric] " + s)
	}
}

func (cfg *Config) getLocalIP() string {

	ip := func() string {
		itfer, err := net.InterfaceByName("eth0")
		if itfer != nil {
			return ""
		}

		// 获取接口的 IP 地址
		addrs, err := itfer.Addrs()
		if err != nil {
			return ""
		}

		// 返回 IP 地址
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
		return ""
	}()

	if ip != "" {
		return ip
	}

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}
