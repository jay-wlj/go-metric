package prom

import (
	"context"

	"go.opentelemetry.io/otel/metric"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/metrics"
)

var _ interfaces.Counter = &Counter{}

type Counter struct {
	base        Base
	counterImpl metric.Float64Counter
}

func NewCounter(name string, checker metrics.Checker, counterImpl metric.Float64Counter) interfaces.Counter {
	return &Counter{
		base: Base{
			name:    name,
			checker: checker,
		},
		counterImpl: counterImpl,
	}
}

func (pc *Counter) Incr(delta float64) {
	if !pc.base.finish() {
		return
	}
	// 超限返回
	if pc.base.ExceedThreshold() {
		return
	}
	pc.counterImpl.Add(context.TODO(), delta, metric.WithAttributes(pc.base.labels...))
}

func (pc *Counter) IncrOnce() {
	pc.Incr(1)
}

func (pc *Counter) AddTag(key, value string) interfaces.Counter {
	pc.base.AddTag(key, value)
	return pc
}

func (pc *Counter) WithTags(tags map[string]string) interfaces.Counter {
	pc.base.WithTags(tags)
	return pc
}
