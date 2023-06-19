package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const kafkaConsumeMetricName = "dacs_kafka_consumer_total"

func newKafkaConsumeCounter(
	meter interfaces.BaseMeter,
	metricNamePrefix string,
	topic string,
	resource string,
) interfaces.ComponentCounter {
	metricName := kafkaConsumeMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "kafka_consumer_total"
	}
	counter := meter.NewCounter(metricName)
	counter.AddTag(topicKey, topic)
	counter.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	return newComponentCounter(counter)
}
