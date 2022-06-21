package nop

import (
	"net/http"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/metrics/nop"
)

var _ interfaces.Components = &components{}

type components struct{}

func (c *components) NewHTTPServerTimer(_ string, _ string, _ int) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewHTTPClientTimer(_, _ string, _ string, _ int) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewHTTPClientTimerFromResponse(_ *http.Response, _ string) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewMysqlTimer(_, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewRedisTimer(_, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}
func (c *components) NewKafkaProduceTimer(_ string, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewKafkaConsumeCounter(_ string, _ string) interfaces.ComponentCounter {
	return nop.ComponentCounter
}

func (c *components) NewESTimer(_, _, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewHBaseTimer(_, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewRMQProduceTimer(_, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}

func (c *components) NewRMQConsumeCounter(_, _ string) interfaces.ComponentCounter {
	return nop.ComponentCounter
}

func (c *components) NewMongoTimer(_, _, _ string, _ bool) interfaces.ComponentTimer {
	return nop.ComponentTimer
}
