package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const httpClientTimerMetricName = "dacs_service_http_call_seconds"

func newHTTPClientTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix string,
	toAppId string,
	serverHost string,
	serverPath string,
	ret string,
	statusCode int,
) interfaces.ComponentTimer {
	metricName := httpClientTimerMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "service_http_call_seconds"
	}
	timer := meter.NewTimer(metricName)

	timer.AddTag(fromAppIdKey, config.GetConfig().AppId)
	if toAppId == "" {
		toAppId = "-"
	}
	timer.AddTag(toAppIdKey, toAppId)
	timer.AddTag(clientIPKey, config.GetConfig().LocalIP)
	timer.AddTag(serverDomainKey, labels.Filter.FilterHost(serverHost))
	timer.AddTag(serverURLApiKey, labels.Filter.FilterRoute(serverPath))
	if statusCode >= 400 && statusCode <= 600 {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	timer.AddTag(retKey, labels.Filter.FilterRet(ret))
	timer.AddTag(statusKey, labels.Filter.FilterStatusCode(statusCode))
	return newComponentTimer(timer)
}
