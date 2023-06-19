package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const kafkaProduceMetricName = "dacs_kafka_producer_seconds"

func newKafkaTimer(
	meter interfaces.BaseMeter,
	metricNamePrefix string,
	topic string,
	resource string,
	hasError bool,
) interfaces.ComponentTimer {
	metricName := kafkaProduceMetricName
	if metricNamePrefix != "" {
		metricName = metricNamePrefix + "kafka_producer_seconds"
	}
	timer := meter.NewTimer(metricName)
	timer.AddTag(topicKey, topic)
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
