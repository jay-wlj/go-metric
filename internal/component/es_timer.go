package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const esTimerMetricName = "dacs_es_request_seconds"

func newESTimer(
	meter interfaces.BaseMeter,
	api, index, resource string,
	hasError bool,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(esTimerMetricName)
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
