package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const esTimerMetricName = "dacs_es_request_seconds"

func newESTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix,
	api, index, resource string,
	hasError bool,
) interfaces.ComponentTimer {
	metricName := esTimerMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "es_request_seconds"
	}

	timer := meter.NewTimer(metricName)
	timer.AddTag(apiKey, labels.Filter.FilterRoute(api))
	timer.AddTag(indexKey, index)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
