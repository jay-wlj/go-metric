package prom

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/jay-wlj/go-metric/internal/labels"
	"github.com/jay-wlj/go-metric/internal/metrics"
)

var _ interfaces.Gauge = &GaugeSeries{}

type GaugeMetric struct {
	metricName string
	lock       sync.Mutex
	series     map[uint64]*GaugeSeries
	observer   metric.Float64GaugeObserver
	checker    metrics.Checker
}

func NewGaugeMetric(metricName string, meter metric.Meter, checker metrics.Checker) (*GaugeMetric, error) {
	gm := &GaugeMetric{
		metricName: metricName,
		series:     make(map[uint64]*GaugeSeries),
		checker:    checker,
	}
	observer, err := meter.NewFloat64GaugeObserver(metricName, gm.ObserveAll)
	if err != nil {
		return nil, err
	}
	gm.observer = observer
	return gm, nil
}

// ObserveAll is a metric-level observation to gather all time-series
func (gm *GaugeMetric) ObserveAll(_ context.Context, result metric.Float64ObserverResult) {
	gm.lock.Lock()
	defer gm.lock.Unlock()

	for _, series := range gm.series {
		result.Observe(series.value, series.labels...)
	}
}

func (gm *GaugeMetric) NewGaugeSeries() interfaces.Gauge {
	return &GaugeSeries{
		metric: gm,
	}
}

func (gm *GaugeMetric) Bind(series *GaugeSeries) {
	hash := series.labels.Hash()
	gm.lock.Lock()
	defer gm.lock.Unlock()

	pgs, ok := gm.series[hash]
	if ok {
		// exist before, just update the value
		pgs.value = series.value
	} else {
		// insert into registry
		gm.series[hash] = series
	}
}

type GaugeSeries struct {
	metric *GaugeMetric
	labels labels.Labels
	value  float64
}

// https://help.aliyun.com/document_detail/208902.html

func (pgs *GaugeSeries) Update(v float64) {
	pgs.value = v
	// 超限返回
	if pgs.metric.checker.ExceedThreshold(pgs.metric.metricName, pgs.labels.Hash()) {
		return
	}
	// register it to the metric
	pgs.metric.Bind(pgs)
}

func (pgs *GaugeSeries) AddTag(key, value string) interfaces.Gauge {
	if pgs.value != 0 {
		return pgs
	}
	pgs.labels = append(pgs.labels, attribute.String(key, value))
	return pgs
}

func (pgs *GaugeSeries) WithTags(tags map[string]string) interfaces.Gauge {
	if pgs.value != 0 {
		return pgs
	}
	pgs.labels = pgs.labels[:0]
	if tags == nil {
		return pgs
	}
	// fix: WithTags 会覆盖 metric_type 类型
	baseLabel := config.GetConfig().BaseLabel
	if baseLabel != nil && baseLabel.MetricyType != "" {
		pgs.AddTag(baseLabel.MetricyType, "gauge")
	} else {
		pgs.AddTag(config.GetConfig().PrefixBaseLabel+"metric_type", "gauge")
	}

	for k, v := range tags {
		pgs.AddTag(k, v)
	}
	return pgs
}
