package gometric

import (
	"errors"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/jay-wlj/go-metric/internal/global"
	"github.com/jay-wlj/go-metric/internal/meter/nop"
	"github.com/jay-wlj/go-metric/internal/meter/prom"
)

func NewMeter(options ...interfaces.Option) interfaces.Meter {
	var cfg = config.GetConfig()
	for _, option := range options {
		option.ApplyConfig(cfg)
	}
	// 测试环境
	if cfg.IsTest() {
		cfg.WriteInfoOrNot("under test environment, using NopMeter")
		m := nop.NewNopMeter(cfg, nil)
		global.SetMeter(m)
		return m
	}

	cfg.WriteInfoOrNot("global labels for all metrics: " + config.DtlLabels().String())
	switch cfg.MeterProvider {
	case PrometheusMeterProvider:
		cfg.WriteInfoOrNot("you are using PrometheusMeter now!")
		bm, err := prom.NewPrometheusMeter(cfg)
		if err != nil {
			return nop.NewNopMeter(cfg, err)
		}
		global.SetMeter(bm)
		return bm
	default:
		m := nop.NewNopMeter(cfg, errors.New("不支持该 MeterProvider"))
		global.SetMeter(m)
		return m
	}
}
