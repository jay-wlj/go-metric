package interfaces

import "time"

type BaseMeter interface {
	WithRunning(on bool) // 设置为false，SDK切换为空实现，关闭指标的收集功能
	NewCounter(metricName string) Counter
	NewGauge(metricName string) Gauge
	NewTimer(metricName string) Timer
}

// Meter 用于管理、创建指标
type Meter interface {
	BaseMeter
	Components() Components // 返回中间件埋点方法
}

// Counter 计数器，适用于 PV/UV，总请求数、总耗时等统计
type Counter interface {
	Incr(delta float64) // Incr(1)
	IncrOnce()          // +1
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key, value string) Counter
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) Counter
}

// Gauge 可增可减，适用于内存用量、cpu利用率等统计
type Gauge interface {
	Update(v float64)
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key, value string) Gauge
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) Gauge
}

// Timer 计时器，底层数据结构为histogram，可以生成99、95线等
// 由于 open-telemetry的设计，暂不支持Buckets的指定
type Timer interface {
	Time(f func())               // 记录函数执行的耗时
	Update(d time.Duration)      // 记录一段耗时
	UpdateSince(start time.Time) // 记录从起始时间的耗时
	UpdateInMillis(m float64)    // 记录一段毫秒单位的耗时
	UpdateInSeconds(s float64)   // 记录一段秒单位的耗时
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key, value string) Timer
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) Timer
}
