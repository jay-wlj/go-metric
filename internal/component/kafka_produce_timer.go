package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const kafkaProduceMetricName = "dtlci_kafka_producer_seconds"

func newKafkaTimer(
	meter interfaces.BaseMeter,
	topic string,
	resource string,
	hasError bool,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(kafkaProduceMetricName)
	timer.AddTag(topicKey, topic)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
