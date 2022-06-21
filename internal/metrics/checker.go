package metrics

// Checker 用于检查是否超限
type Checker interface {
	// ExceedThreshold 判定是否超限，如果超限，则上报对应的错误
	ExceedThreshold(metricName string, seriesID uint64) bool
}
