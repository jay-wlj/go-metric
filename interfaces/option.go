package interfaces

import "github.com/jay-wlj/go-metric/internal/config"

// Option 用于控制初始化MeterProvider
// 比如 WithMeterProvider(PrometheusMeter) 切换监控存储为Prometheus
//     WithRuntimeStatsCollectorInEvery(time.Second * 30) 激活runtime数据采集
type Option interface {
	ApplyConfig(c *config.Config)
}
