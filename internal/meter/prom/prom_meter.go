package prom

import (
	"net/http"

	cli_prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/sdk/metric"

	"sync"
	"sync/atomic"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/component"
	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/jay-wlj/go-metric/internal/global"
	"github.com/jay-wlj/go-metric/internal/metrics/nop"
	"github.com/jay-wlj/go-metric/internal/metrics/prom"
	"github.com/jay-wlj/go-metric/internal/runtime"
)

const (
	sdkVersion          = "0.1"
	PrometheusMeterName = "github.com/jay-wlj/go-metric/prometheus-meter"
	maxMetricCount      = 150  // 最大metric数量
	maxSeriesCount      = 1000 // 最大时间线数量
)

var (
	defaultHistogramBoundaries = []float64{
		0.002, 0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 10,
	}

	_ interfaces.Meter = &PrometheusMeter{}
)

type PrometheusMeter struct {
	cfg              *config.Config
	running          int32
	onCh             chan struct{} // receiving start signal
	offCh            chan struct{} // receiving stop signal
	meter            api.Meter
	runtimeCollector runtime.Collector
	// http server
	servers []HTTPServer
	// pushServer HTTPServer
	// used for checker
	allMetricsLock sync.RWMutex
	allMetrics     map[string]*seriesGroup // metricname->seriesID group
	// internal metrics registry
	gaugesLock     sync.RWMutex
	gaugesRegistry map[string]*prom.GaugeMetric // observe

	handler    http.Handler      // http handler func for metrics
	publicTags map[string]string // public labels
}

func NewPrometheusMeter(cfg *config.Config) (*PrometheusMeter, error) {
	// prometheusCfg := prometheus.Config{
	// 	Registry:                   cli_prom.NewRegistry(),
	// 	DefaultHistogramBoundaries: defaultHistogramBoundaries}

	// if len(cfg.HistogramBoundaries) > 0 {
	// 	prometheusCfg.DefaultHistogramBoundaries = cfg.HistogramBoundaries
	// }
	// ctrl := controller.New(
	// 	processor.NewFactory(
	// 		selector.NewWithHistogramDistribution(
	// 			histogram.WithExplicitBoundaries(prometheusCfg.DefaultHistogramBoundaries),
	// 		),
	// 		aggregation.CumulativeTemporalitySelector(),
	// 		processor.WithMemory(true),
	// 	),
	// 	controller.WithResource(config.DtlResource()),
	// )
	// exporter, err := prometheus.New(prometheusCfg, ctrl)
	registy := cli_prom.NewRegistry()
	exporter, err := prometheus.New(
		prometheus.WithRegisterer(registy),
		// prometheus.WithoutUnits(), // 对于带上单位(时间,字节等)的指标值，指标名是否要自动加上单位后缀，如counter类型的指标名request.duration，带有seconds单位，则指标名会变成request_duration_seconds_total
		// prometheus.WithoutCounterSuffixes(),	// 对counter类型的指标，是否要去掉加上_total后缀
		// prometheus.WithAggregationSelector(metrics)
		// prometheus.WithResourceAsConstantLabels(attribute.NewAllowKeysFilter(config.DtlLabels().Keys()...))
		prometheus.WithoutScopeInfo())

	if err != nil {
		cfg.WriteErrorOrNot("failed to initialize Prometheus Meter: " + err.Error())
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(config.DtlResource()),
		metric.WithReader(exporter),
		metric.WithView(metric.NewView(
			metric.Instrument{Kind: metric.InstrumentKindHistogram}, // histogram类型指标指标统一bucket桶
			metric.Stream{Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: cfg.HistogramBoundaries,
			}},
		)))
	meter := provider.Meter(PrometheusMeterName,
		// api.WithInstrumentationAttributes(config.DtlLabels()...),
		api.WithInstrumentationVersion(sdkVersion))

	handler := promhttp.HandlerFor(registy, promhttp.HandlerOpts{})
	pm := PrometheusMeter{
		cfg:     cfg,
		running: 1,
		onCh:    make(chan struct{}),
		offCh:   make(chan struct{}),
		// server:     newPromHTTPServer(cfg, exporter.ServeHTTP),
		allMetrics:     make(map[string]*seriesGroup),
		meter:          meter,
		gaugesRegistry: make(map[string]*prom.GaugeMetric),
		handler:        handler,
		publicTags:     config.DtlLabels().Map(),
	}

	// push方式不需要exporter
	if cfg.PrometheusPort > 0 {
		pm.servers = append(pm.servers, newPromHTTPServer(cfg, pm.GetHandlerFunc()))
	}
	if cfg.Push != nil {
		pm.servers = append(pm.servers, newPromPushServer(cfg, registy))
	}

	pm.runtimeCollector = runtime.NewCollector(cfg, &pm)
	pm.runtimeCollector.Start()

	for _, server := range pm.servers {
		server.Start()
	}

	go pm.signalListener()
	return &pm, nil
}

func (pm *PrometheusMeter) GetHandler() http.Handler {
	return pm.handler
}

func (pm *PrometheusMeter) GetHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pm.handler != nil {
			pm.handler.ServeHTTP(w, r)
		}
	}
}

func (pm *PrometheusMeter) WithRunning(on bool) {
	if on {
		select {
		case pm.onCh <- struct{}{}:
		default:
			// another thread is starting now
		}
	} else {
		select {
		case pm.offCh <- struct{}{}:
		default:
			// another thread is stopping now
		}
	}
}

func (pm *PrometheusMeter) signalListener() {
	for {
		select {
		case <-pm.onCh:
			if !atomic.CompareAndSwapInt32(&pm.running, 0, 1) {
				return
			}
			pm.cfg.WriteInfoOrNot("WithRunning=true, meter starting...")
			pm.runtimeCollector.Start()
			for _, server := range pm.servers {
				server.Start()
			}
			// replace meter
			var m interfaces.Meter = pm
			global.SetMeter(m)
		case <-pm.offCh:
			if !atomic.CompareAndSwapInt32(&pm.running, 1, 0) {
				return
			}
			// replace meter
			global.SetNopMeter()
			pm.cfg.WriteInfoOrNot("WithRunning=false, meter stopping...")
			pm.runtimeCollector.Stop()
			for _, server := range pm.servers {
				server.Stop()
			}
			// clear internal gauges
			pm.gaugesLock.Lock()
			pm.gaugesRegistry = make(map[string]*prom.GaugeMetric)
			pm.gaugesLock.Unlock()
		}
	}
}

func (pm *PrometheusMeter) NewCounter(metricName string) interfaces.Counter {
	if atomic.LoadInt32(&pm.running) == 0 {
		return &nop.Counter
	}
	c, err := pm.meter.Float64Counter(metricName)
	if err != nil {
		return &nop.Counter
	}
	if pm.cfg.BaseLabel == nil || pm.cfg.BaseLabel.MetricyType == "" {
		return prom.NewCounter(metricName, pm, c).WithTags(pm.publicTags)
	}
	return prom.NewCounter(metricName, pm, c).AddTag(pm.cfg.PrefixBaseLabel+"metric_type", "counter")
}

func (pm *PrometheusMeter) NewGauge(metricName string) interfaces.Gauge {
	if atomic.LoadInt32(&pm.running) == 0 {
		return &nop.Gauge
	}
	pm.gaugesLock.RLock()
	gaugeMetric, ok := pm.gaugesRegistry[metricName]
	pm.gaugesLock.RUnlock()
	if ok {
		if pm.cfg.BaseLabel == nil || pm.cfg.BaseLabel.MetricyType == "" {
			return gaugeMetric.NewGaugeSeries().WithTags(pm.publicTags)
		}
		return gaugeMetric.NewGaugeSeries().AddTag(pm.cfg.PrefixBaseLabel+"metric_type", "gauge")
	}
	// not exist before
	pm.gaugesLock.Lock()
	defer pm.gaugesLock.Unlock()
	// double check for concurrency
	gaugeMetric, ok = pm.gaugesRegistry[metricName]
	if !ok {
		var err error
		gaugeMetric, err = prom.NewGaugeMetric(metricName, pm.meter, pm)
		if err != nil {
			return &nop.Gauge
		}
		pm.gaugesRegistry[metricName] = gaugeMetric
	}
	if pm.cfg.BaseLabel == nil || pm.cfg.BaseLabel.MetricyType == "" {
		return gaugeMetric.NewGaugeSeries().WithTags(pm.publicTags)
	}
	return gaugeMetric.NewGaugeSeries().AddTag(pm.cfg.PrefixBaseLabel+"metric_type", "gauge")
}

func (pm *PrometheusMeter) NewTimer(metricName string) interfaces.Timer {
	if atomic.LoadInt32(&pm.running) == 0 {
		return &nop.Timer
	}
	// t, err := pm.meter.Float64Histogram(metricName, api.WithExplicitBucketBoundaries(pm.cfg.HistogramBoundaries...))
	t, err := pm.meter.Float64Histogram(metricName)
	if err != nil {
		return &nop.Timer
	}
	if pm.cfg.BaseLabel == nil || pm.cfg.BaseLabel.MetricyType == "" {
		return prom.NewTimer(metricName, pm, t).WithTags(pm.publicTags)
	}
	return prom.NewTimer(metricName, pm, t).AddTag(pm.cfg.PrefixBaseLabel+"metric_type", "histogram")
}

func (pm *PrometheusMeter) Components() interfaces.Components {
	return component.NewComponents(pm.cfg, pm)
}

func (pm *PrometheusMeter) ExceedThreshold(metricName string, seriesID uint64) bool {
	// 放行超限埋点
	if metricName == "TooManyMetric" || metricName == "TooManyValue" {
		return false
	}
	tooManyMetric := func() { pm.NewCounter("TooManyMetric").AddTag("name", metricName).IncrOnce() }
	tooManySeries := func() { pm.NewCounter("TooManyValue").AddTag("name", metricName).IncrOnce() }

	pm.allMetricsLock.RLock()
	seriesGroup, metricExist := pm.allMetrics[metricName]
	metricsCount := len(pm.allMetrics)
	pm.allMetricsLock.RUnlock()
	if metricExist {
		// case1, metric存在，但组合超限
		if seriesGroup.ExceedThreshold(seriesID) {
			tooManySeries()
			return true
		}
		// case2, metric存在，组合未超限
		return false
	}
	// case3, metric不存在，但组合已超限
	if metricsCount >= maxMetricCount {
		tooManyMetric()
		return true
	}

	pm.allMetricsLock.Lock()
	defer pm.allMetricsLock.Unlock()
	// double check
	seriesGroup, metricExist = pm.allMetrics[metricName]
	metricsCount = len(pm.allMetrics)
	if metricExist {
		// case4,上互斥锁后，发现metric已被其他线程创建，且组合数已超限
		if seriesGroup.ExceedThreshold(seriesID) {
			tooManySeries()
			return true
		}
		// case5,上互斥锁后，发现metric已被其他线程创建，但组合数未超限
		return false
	}
	// case6,上互斥锁后，metric未被其他线程创建，但metric组合数已超限
	if metricsCount >= maxMetricCount {
		tooManyMetric()
		return true
	}
	// case7,上互斥锁后，metric未被其他线程创建，增加该series组合
	seriesGroup = newSeriesGroup(seriesID)
	pm.allMetrics[metricName] = seriesGroup
	return false
}
