package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const httpServerTimeMetricName = "dtlci_api_request_seconds"

func newHTTPServerTimer(
	meter interfaces.BaseMeter,
	route string,
	ret string,
	statusCode int,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(httpServerTimeMetricName)
	timer.AddTag(routeKey, labels.Filter.FilterRoute(route))
	timer.AddTag(retKey, labels.Filter.FilterRet(ret))
	timer.AddTag(statusKey, labels.Filter.FilterStatusCode(statusCode))
	return newComponentTimer(timer)
}
