package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const hbaseTimerMetricName = "dacs_hbase_request_seconds"

func newHBaseTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix,
	cmd, resource string,
	hasError bool,
) interfaces.ComponentTimer {
	metricName := hbaseTimerMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "hbase_request_seconds"
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
