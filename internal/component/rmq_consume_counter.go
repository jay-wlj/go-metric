package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const rmqConsumeMetricName = "dtlci_rabbit_consumer_total"

func newRMQConsumeCounter(
	meter interfaces.BaseMeter,
	queue string,
	resource string,
) interfaces.ComponentCounter {
	counter := meter.NewCounter(rmqConsumeMetricName)
	counter.AddTag(queueKey, queue)
	counter.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	return newComponentCounter(counter)
}
