package prom

import (
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"

	"github.com/jay-wlj/go-metric/internal/labels"
	"github.com/jay-wlj/go-metric/internal/metrics"
)

type Base struct {
	name      string
	labels    labels.Labels
	checker   metrics.Checker
	completed int32 // 标识是否完成了初始化
}

func (pb *Base) finish() bool {
	return atomic.CompareAndSwapInt32(&pb.completed, 0, 1)
}

func (pb *Base) AddTag(key, value string) {
	pb.labels = append(pb.labels, attribute.String(key, value))
}

func (pb *Base) WithTags(tags map[string]string) {
	if tags == nil || len(tags) == 0 {
		return
	}
	for k, v := range tags {
		pb.AddTag(k, v)
	}
}

func (pb *Base) ExceedThreshold() bool {
	return pb.checker.ExceedThreshold(pb.name, pb.labels.Hash())
}
