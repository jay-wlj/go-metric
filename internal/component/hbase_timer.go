package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const hbaseTimerMetricName = "dacs_hbase_request_seconds"

func newHBaseTimer(
	meter interfaces.BaseMeter,
	cmd, resource string,
	hasError bool,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(hbaseTimerMetricName)
	timer.AddTag(cmdKey, cmd)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
