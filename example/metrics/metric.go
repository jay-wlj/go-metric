package metric

import (
	"errors"
	"strings"
	"time"

	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/interfaces"
)

type MetricPushConfig struct {
	// unit second, every `period` second will push metrics to mars_executor
	Period int `json:"period" yaml:"period"`
	// mars_executor push address, default is http://127.0.0.1:7072
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// http push timeout, second, default is 5s
	Timeout int `json:"timeout" yaml:"timeout"`
}

// InitMetrics metric 模块初始化
/*
	参数说明:
	sname: 服务名,eg: config_server|auth_server
	sver: 服务二进制版本
	prometheusPort: 新开监控侦听端口(默认0不开启),方便 prometheus 来抓取指标,eg: http://127.0.0.1:16700/metrics
	MetricPushConfig: 主动推送指标到push_gateway的地址及间隔
*/
func InitMetrics(sname string, metricsPort int, pushCfg *MetricPushConfig) error {
	if sname == "" {
		meili_zap.Logger.Sugar().Errorf("InitMetrics fail! sname is empty!")
		return errors.New("InitMetrics fail!")
	}

	// vers := strings.Split(version.GetVersion(), "-")
	ops := []interfaces.Option{
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
		gometric.WithPrometheusPort(metricsPort),
		gometric.WithAppID(sname),
		// gometric.WithAppVer(vers[0]),
		gometric.WithHistogramBoundaries([]float64{0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 10, 15}), // histogram默认分桶序列,单位:秒, 不要随便变动
		gometric.WithBaseLabelNames("dc_sname", "", "dc_ip", "", "dc_sver"),
	}

	if pushCfg != nil && pushCfg.Endpoint != "" {
		var period = time.Duration(pushCfg.Period) * time.Second
		if period == 0 {
			period = time.Duration(15 * time.Second)
		}
		ops = append(ops, gometric.WithPush(pushCfg.Endpoint, period))
	}
	_ = gometric.NewMeter(
		ops...,
	)

	// 服务初始化接入指标
	// NewGuage("dacs_server_up").AddTag("up_time", time.Now().String()).Update(1)
	NewCounter("dacs_server_up_total").AddTag("up_time", time.Now().String()).IncrOnce()

	return nil
}

func GetMeter() interfaces.Meter {
	return gometric.GetGlobalMeter()
}

// NewTimer 上报histogram指标，关注该指标耗时性能等使用histogram
/*
	参数说明:
	metricName: 指标名
*/
func NewTimer(metricName string) interfaces.Timer {
	return gometric.GetGlobalMeter().NewTimer(metricName)
}

// NewTimer 上报Counter指标, 不关注耗时，只关注qps及分布情况等，请使用counter
/*
	参数说明:
	metricName: 指标名
*/
func NewCounter(metricName string) interfaces.Counter {
	return gometric.GetGlobalMeter().NewCounter(metricName)
}

// NewTimer 上报guage指标，类似cpu/mem这些实时变动的数值请使用guage类型
/*
	参数说明:
	metricName: 指标名
*/
func NewGuage(metricName string) interfaces.Gauge {
	return gometric.GetGlobalMeter().NewGauge(metricName)
}
