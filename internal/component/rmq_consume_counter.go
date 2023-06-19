package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const rmqConsumeMetricName = "dacs_rabbit_consumer_total"

func newRMQConsumeCounter(
	meter interfaces.BaseMeter,
	metricNamePrefix string,
	queue string,
	resource string,
) interfaces.ComponentCounter {
	metricName := rmqConsumeMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "rabbit_consumer_total"
	}
	counter := meter.NewCounter(metricName)
	counter.AddTag(queueKey, queue)
	counter.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	return newComponentCounter(counter)
}
