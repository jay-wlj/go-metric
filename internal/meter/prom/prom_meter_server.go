package prom

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/jay-wlj/go-metric/internal/config"

	consulapi "github.com/hashicorp/consul/api"
)

type HTTPServer interface {
	// Start 启动 consul register，并定期重试保证再注册
	Start()
	// Stop 关闭 register，解注册
	Stop()
}

type promHTTPServer struct {
	cfg          *config.Config
	server       *http.Server
	consulClient *consulapi.Client
	running      int32
	CloseCh      chan struct{}
}

func newPromHTTPServer(cfg *config.Config, exporterHandler http.HandlerFunc) HTTPServer {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.PrometheusPort),
		Handler: mux,
	}
	server := promHTTPServer{
		cfg:     cfg,
		server:  srv,
		running: 0,
		CloseCh: make(chan struct{}),
	}

	logRoute := func(route string) string {
		cfg.WriteInfoOrNot(fmt.Sprintf("http handler, GET  http://%s:%d%s",
			cfg.LocalIP,
			cfg.PrometheusPort,
			route,
		))
		return route
	}
	// for consul
	mux.HandleFunc(logRoute("/actuator/health"), server.healthCheck)
	// See https://dtl.feishu.cn/docs/doccnW3H5g3uMeGbREggwq3aNyf
	mux.HandleFunc(logRoute("/metrics"), exporterHandler)
	mux.HandleFunc(logRoute("/debug/pprof/"), pprof.Index)
	mux.HandleFunc(logRoute("/debug/pprof/cmdline"), pprof.Cmdline)
	mux.HandleFunc(logRoute("/debug/pprof/profile"), pprof.Profile)
	mux.HandleFunc(logRoute("/debug/pprof/symbol"), pprof.Symbol)
	mux.HandleFunc(logRoute("/debug/pprof/trace"), pprof.Trace)

	go server.startHTTPServer()
	return &server
}

func (p *promHTTPServer) Start() {
	// register is already running
	if !(atomic.CompareAndSwapInt32(&p.running, 0, 1)) {
		return
	}
	go p.Register()
}

func (p *promHTTPServer) Stop() {
	// register is already stopped
	if !(atomic.CompareAndSwapInt32(&p.running, 1, 0)) {
		return
	}
	p.CloseCh <- struct{}{}
}

func (p *promHTTPServer) Register() {
	const registerRetryDuration = time.Minute * 5
	p.cfg.WriteInfoOrNot("consul register is running, it registers in every " + registerRetryDuration.String())
	registerTicker := time.NewTicker(registerRetryDuration)
	defer registerTicker.Stop()

	// first try, always log
	if p.checkConsulClient() {
		registration := p.consulRegistration()
		bs, _ := json.Marshal(*p.consulRegistration())
		if err := p.consulClient.Agent().ServiceRegister(registration); err == nil {
			p.cfg.WriteInfoOrNot("register service on consul successfully")
		} else {
			p.cfg.WriteErrorOrNot(fmt.Sprintf(
				"failed to register service on consul: %s, service: %s", err.Error(), string(bs)))
		}
	}
	// loop check
	for {
		select {
		case <-registerTicker.C:
			if !p.checkConsulClient() {
				continue
			}
			registration := p.consulRegistration()
			bs, _ := json.Marshal(*p.consulRegistration())
			if err := p.consulClient.Agent().ServiceRegister(registration); err != nil {
				p.cfg.WriteErrorOrNot(fmt.Sprintf(
					"failed to register service on consul: %s, service: %s",
					err.Error(), string(bs),
				))
			}
		// 解注册、并关闭
		case <-p.CloseCh:
			if !p.checkConsulClient() {
				return
			}
			registration := p.consulRegistration()
			if err := p.consulClient.Agent().ServiceDeregister(registration.ID); err == nil {
				p.cfg.WriteInfoOrNot("deregister service on consul successfully")
			} else {
				p.cfg.WriteErrorOrNot("failed to deregister service on consul: " + err.Error())
			}
			return
		}
	}
}

func (p *promHTTPServer) checkConsulClient() bool {
	if p.cfg.Consul == nil {
		return false
	}
	if p.cfg.Consul != nil && p.consulClient == nil {
		consulCfg := consulapi.DefaultConfig()
		consulCfg.Address = p.cfg.Consul.ConsulAddress
		consulCfg.Token = p.cfg.Consul.ConsulToken
		client, err := consulapi.NewClient(consulCfg)
		if err != nil {
			p.cfg.WriteErrorOrNot(fmt.Sprintf(
				"failed to initialize consule client: %s",
				err.Error()))
			return false
		}
		p.consulClient = client
	}
	return true
}
func (p *promHTTPServer) consulRegistration() *consulapi.AgentServiceRegistration {
	return &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s", p.cfg.AppId, p.cfg.LocalIP),
		Name:    "go-sdk-exporter",
		Port:    p.cfg.PrometheusPort,
		Address: p.cfg.LocalIP,
		Tags:    []string{"go"},
		// 供 prometheus 使用查看对应的服务 appid
		Meta: map[string]string{"appid": config.GetConfig().AppId},
		Check: &consulapi.AgentServiceCheck{
			HTTP: fmt.Sprintf("http://%s:%d/actuator/health",
				p.cfg.LocalIP, p.cfg.PrometheusPort),
			Timeout:  "5s",
			Interval: "5s",
		},
	}
}

func (p *promHTTPServer) startHTTPServer() {
	p.cfg.WriteInfoOrNot("prom http server listen and server on: " + strconv.Itoa(p.cfg.PrometheusPort))
	err := p.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		p.cfg.WriteErrorOrNot(fmt.Sprintf("faield to start prom http server on : %d with error: %s ",
			p.cfg.PrometheusPort, err.Error()))
	}
}

func (p *promHTTPServer) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/json")
	msg, _ := json.Marshal(map[string]interface{}{"status": "UP"})
	_, _ = w.Write(msg)
}
