package gometric

import (
	"time"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
)

const (
	// PrometheusMeterProvider
	// 现阶段只提供一种存储，需要显式的定义
	PrometheusMeterProvider config.MeterProviderType = iota + 1
)

// WithMeterProvider 用于选择对应的监控存储
func WithMeterProvider(mpt config.MeterProviderType) interfaces.Option {
	return &meterProviderOption{mpt: mpt}
}

type envOption struct{ env string }

func (eo *envOption) ApplyConfig(c *config.Config) { c.Env = eo.env }

// WithEnv 用于设置环境, 如 test, dev, stg, pre, gray, prod
func WithEnv(env string) interfaces.Option { return &envOption{env: env} }

type consulOption struct {
	address string
	token   string
}

func (co *consulOption) ApplyConfig(c *config.Config) {
	c.Consul = &config.ConsulCfg{
		ConsulAddress: co.address,
		ConsulToken:   co.token,
	}
}

// WithConsul 用于设置consul的地址, 如 "consul-dev.my.cn:8500"
func WithConsul(address string, token string) interfaces.Option {
	return &consulOption{address: address, token: token}
}

type pushOption struct {
	address string
	period  time.Duration
}

func (co *pushOption) ApplyConfig(c *config.Config) {
	c.Push = &config.PushCfg{
		PushAddress: co.address,
		PushPeriod:  co.period,
	}
}

func WithPush(address string, period time.Duration) interfaces.Option {
	return &pushOption{address: address, period: period}
}

type appidOption struct{ appid string }

func (ao *appidOption) ApplyConfig(c *config.Config) {
	c.AppId = ao.appid
}

// WithAppID 用于设置 用户的 appid
func WithAppID(appid string) interfaces.Option { return &appidOption{appid: appid} }

type prometheusPortOption struct{ port int }

func (ppo *prometheusPortOption) ApplyConfig(c *config.Config) { c.PrometheusPort = ppo.port }

// WithPrometheusPort 用于设置 prometheus http 端口，默认为 16670
func WithPrometheusPort(port int) interfaces.Option { return &prometheusPortOption{port: port} }

type meterProviderOption struct{ mpt config.MeterProviderType }

func (mpo *meterProviderOption) ApplyConfig(c *config.Config) {
	c.MeterProvider = mpo.mpt
}

// WithRuntimeStatsCollector 用于开启采集 runtime 信息
// 不设置则不采集
func WithRuntimeStatsCollector() interfaces.Option {
	return &runtimeStatsOption{}
}

type runtimeStatsOption struct{}

func (rso *runtimeStatsOption) ApplyConfig(c *config.Config) {
	c.ReadRuntimeStats = true
}

// WithInfoLogWrite 用于指定输出 info 格式日志的方法
// sdk内会尽可能抑制日志出现的频次，因此不会对业务产生较大影响
// case1. 假如，你使用zap.Logger，可以使用闭包指定
//
//	WithInfoLogWrite(func(infoLine string){
//		logger.Info(infoLine)
//	})
//
// case2. 未指定时，默认使用了标准输出，可以定向到os.Stdout
//
//			WithInfoLogWrite(func(infoLine string){
//				_, _ = os.Stdout.WriteString(infoLine)
//	         _, _ = os.Stdout.WriteString("\n")
//			})
func WithInfoLogWrite(infoWriter func(infoLine string)) interfaces.Option {
	return &infoLogWriteOption{writer: infoWriter}
}

type infoLogWriteOption struct{ writer func(s string) }

func (lwo *infoLogWriteOption) ApplyConfig(c *config.Config) { c.InfoLogWrite = lwo.writer }

// WithErrorLogWrite 用法与 WithInfoLogWrite 一致
func WithErrorLogWrite(errWriter func(errorLine string)) interfaces.Option {
	return &errLogWriterOption{writer: errWriter}
}

type errLogWriterOption struct{ writer func(s string) }

func (lwo *errLogWriterOption) ApplyConfig(c *config.Config) { c.ErrorLogWrite = lwo.writer }

type labelNameOption struct {
	Appid       string
	Env         string
	IP          string
	DataType    string
	MetricyType string
}

// WithBaseLabelNames 用于设置 基础标签名
func WithBaseLabelNames(appid, env, ip, data_type string) interfaces.Option {
	return &labelNameOption{
		Appid:    appid,
		Env:      env,
		IP:       ip,
		DataType: data_type,
	}
}

func (t *labelNameOption) ApplyConfig(c *config.Config) {
	c.BaseLabel = &config.BaseLabelCfg{
		Appid:       t.Appid,
		Env:         t.Env,
		IP:          t.IP,
		DataType:    t.DataType,
		MetricyType: t.MetricyType,
	}
}

type prefixlabelNameOption struct {
	prefixBaseLablename string
}

// WithBaseLabelNames 用于设置 基础标签名前缀
func WithPrefixBaseLabelName(prefixLabelName string) interfaces.Option {
	return &prefixlabelNameOption{
		prefixBaseLablename: prefixLabelName,
	}
}

func (t *prefixlabelNameOption) ApplyConfig(c *config.Config) {
	c.PrefixBaseLabel = t.prefixBaseLablename
	c.BaseLabel = &config.BaseLabelCfg{
		Appid:       t.prefixBaseLablename + "appid",
		Env:         t.prefixBaseLablename + "env",
		IP:          t.prefixBaseLablename + "ip",
		DataType:    t.prefixBaseLablename + "data_type",
		MetricyType: t.prefixBaseLablename + "metricy_type",
	}
}

type metricNameprefixOption struct {
	prefixMetricName string
}

// WithPrefixMetricName 用于设置 中间件指标名前缀
func WithPrefixMetricName(prefixMetricName string) interfaces.Option {
	return &metricNameprefixOption{
		prefixMetricName: prefixMetricName,
	}
}

func (t *metricNameprefixOption) ApplyConfig(c *config.Config) {
	c.PrefixMetricName = t.prefixMetricName

}
