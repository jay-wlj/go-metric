package prom

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/metrics"
)

var _ interfaces.Timer = &Timer{}

type Timer struct {
	base      Base
	timerImpl metric.Float64Histogram
}

func NewTimer(name string, checker metrics.Checker, timerImpl metric.Float64Histogram) interfaces.Timer {
	return &Timer{
		base: Base{
			name:    name,
			checker: checker,
		},
		timerImpl: timerImpl,
	}
}

func (pt *Timer) Time(f func()) {
	start := time.Now()
	f()

	pt.UpdateSince(start)
}

func (pt *Timer) Update(d time.Duration) {
	pt.UpdateInSeconds(d.Seconds())
}

func (pt *Timer) UpdateSince(start time.Time) {
	elapsed := time.Now().Sub(start).Seconds()
	pt.UpdateInSeconds(elapsed)
}

func (pt *Timer) UpdateInMillis(m float64) {
	pt.UpdateInSeconds(m / 1000)
}

func (pt *Timer) UpdateInSeconds(s float64) {
	if !pt.base.finish() {
		return
	}
	if pt.base.ExceedThreshold() {
		return
	}
	pt.timerImpl.Record(context.TODO(), s, pt.base.labels...)
}

func (pt *Timer) AddTag(key, value string) interfaces.Timer {
	pt.base.AddTag(key, value)
	return pt
}

func (pt *Timer) WithTags(tags map[string]string) interfaces.Timer {
	pt.base.WithTags(tags)
	return pt
}
