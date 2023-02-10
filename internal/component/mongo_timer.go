package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const mongoTimerMetricName = "dacs_mongo_request_seconds"

func newMongoTimer(
	meter interfaces.BaseMeter,
	command, collection, resource string,
	hasError bool,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(mongoTimerMetricName)
	timer.AddTag(commandKey, command)
	timer.AddTag(collectionKey, collection)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
