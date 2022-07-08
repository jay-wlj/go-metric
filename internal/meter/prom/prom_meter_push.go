package prom

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type promPushServer struct {
	cfg     *config.Config
	pusher  *push.Pusher
	running int32
	CloseCh chan struct{}
}

func newPromPushServer(cfg *config.Config, g prometheus.Gatherer) HTTPServer {
	pushServer := promPushServer{
		cfg:     cfg,
		running: 0,
		CloseCh: make(chan struct{}),
	}
	pushServer.pusher = push.New(cfg.Push.PushAddress, cfg.LocalIP).Gatherer(g)

	return &pushServer
}

func (p *promPushServer) Start() {
	// register is already running
	if !(atomic.CompareAndSwapInt32(&p.running, 0, 1)) {
		return
	}
	go p.push()
}

func (p *promPushServer) Stop() {
	// register is already stopped
	if !(atomic.CompareAndSwapInt32(&p.running, 1, 0)) {
		return
	}
	p.CloseCh <- struct{}{}
}

func (p *promPushServer) push() {
	pushTicker := time.NewTicker(p.cfg.Push.PushPeriod)
	defer pushTicker.Stop()

	now := time.Now()
	if err := p.pusher.Push(); err != nil {
		p.cfg.WriteErrorOrNot(fmt.Sprintf("push fail! tick=%s err=%v\n", time.Now().Sub(now), err))
		// break
	} else {
		p.cfg.WriteInfoOrNot(fmt.Sprintln("push success! tick=", time.Now().Sub(now), " now=", time.Now().Local().String()))
	}

	for {
		select {
		case <-pushTicker.C:
			now := time.Now()
			if err := p.pusher.Push(); err != nil {
				p.cfg.WriteErrorOrNot(fmt.Sprintf("push fail! tick=%s err=%v\n", time.Now().Sub(now), err))
				break
			}
			p.cfg.WriteInfoOrNot(fmt.Sprintln("push success! tick=", time.Now().Sub(now), " now=", time.Now().Local().String()))
		case <-p.CloseCh:
			p.cfg.WriteInfoOrNot("push service exist success")
			return
		}
	}
}
