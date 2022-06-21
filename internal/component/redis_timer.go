package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const redisMetricName = "dtlci_redis_request_seconds"

func newRedisTimer(
	meter interfaces.BaseMeter,
	cmd string,
	resource string,
	hasError bool,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(redisMetricName)
	timer.AddTag(cmdKey, cmd)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
