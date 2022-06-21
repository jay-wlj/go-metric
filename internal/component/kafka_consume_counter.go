package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const kafkaConsumeMetricName = "dtlci_kafka_consumer_total"

func newKafkaConsumeCounter(
	meter interfaces.BaseMeter,
	topic string,
	resource string,
) interfaces.ComponentCounter {
	counter := meter.NewCounter(kafkaConsumeMetricName)
	counter.AddTag(topicKey, topic)
	counter.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	return newComponentCounter(counter)
}
