package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const rmqProduceTimerMetricName = "dacs_rabbit_producer_seconds"

func newRMQProduceTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix, exchange, resource string,
	hasError bool,
) interfaces.ComponentTimer {
	metricName := esTimerMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "rabbit_producer_seconds"
	}
	timer := meter.NewTimer(metricName)
	timer.AddTag(exchangeKey, exchange)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
