package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const redisMetricName = "dacs_redis_request_seconds"

func newRedisTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix string,
	cmd string,
	resource string,
	hasError bool,
) interfaces.ComponentTimer {
	metricName := redisMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "redis_request_seconds"
	}
	timer := meter.NewTimer(metricName)
	timer.AddTag(cmdKey, cmd)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
