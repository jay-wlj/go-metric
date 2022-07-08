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

type Config struct {
	PrometheusPort   int // prometheus related config
	Consul           *ConsulCfg
	Push             *PushCfg // push cfg
	MeterProvider    MeterProviderType
	LocalIP          string
	Env              string
	AppId            string
	ReadRuntimeStats bool
	InfoLogWrite     func(s string)
	ErrorLogWrite    func(s string)
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
