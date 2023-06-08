package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const httpServerTimeMetricName = "dacs_api_request_seconds"

func newHTTPServerTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix string,
	route string,
	ret string,
	statusCode int,
) interfaces.ComponentTimer {
	metricName := esTimerMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "api_request_seconds"
	}
	timer := meter.NewTimer(metricName)
	timer.AddTag(routeKey, labels.Filter.FilterRoute(route))
	timer.AddTag(retKey, labels.Filter.FilterRet(ret))
	timer.AddTag(statusKey, labels.Filter.FilterStatusCode(statusCode))
	return newComponentTimer(timer)
}
