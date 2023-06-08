package component

import (
	"net/http"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/jay-wlj/go-metric/internal/metrics/nop"
)

var _ interfaces.Components = &components{}

type components struct {
	cfg       *config.Config
	baseMeter interfaces.BaseMeter
}

func NewComponents(cfg *config.Config, baseMeter interfaces.BaseMeter) interfaces.Components {
	return &components{
		cfg:       cfg,
		baseMeter: baseMeter}
}

func (c *components) NewHTTPServerTimer(route string, ret string, statusCode int) interfaces.ComponentTimer {
	return newHTTPServerTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		route, ret, statusCode,
	)
}

func (c *components) NewHTTPClientTimer(serverDomain, serverPath string, serverRet string, statusCode int) interfaces.ComponentTimer {
	return newHTTPClientTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		config.AppIdQuerier(c.cfg).GetAppId(serverDomain),
		serverDomain,
		serverPath,
		serverRet,
		statusCode,
	)
}

func (c *components) NewHTTPClientTimerFromResponse(resp *http.Response, serverRet string) interfaces.ComponentTimer {
	if resp == nil {
		return nop.ComponentTimer
	}

	return newHTTPClientTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		config.AppIdQuerier(c.cfg).GetAppId(resp.Request.Host),
		resp.Request.Host,
		resp.Request.URL.Path,
		serverRet,
		resp.StatusCode,
	)
}

func (c *components) NewMysqlTimer(sql string, resource string, hasError bool) interfaces.ComponentTimer {
	return newMysqlTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		sql,
		resource,
		hasError,
	)
}

func (c *components) NewRedisTimer(cmd string, resource string, hasError bool) interfaces.ComponentTimer {
	return newRedisTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		cmd,
		resource,
		hasError,
	)
}

func (c *components) NewKafkaProduceTimer(topic string, resource string, hasError bool) interfaces.ComponentTimer {
	return newKafkaTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		topic,
		resource,
		hasError,
	)
}

func (c *components) NewKafkaConsumeCounter(topic string, resource string) interfaces.ComponentCounter {
	return newKafkaConsumeCounter(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		topic,
		resource,
	)
}

func (c *components) NewESTimer(api, index, resource string, hasError bool) interfaces.ComponentTimer {
	return newESTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		api,
		index,
		resource,
		hasError,
	)
}

func (c *components) NewHBaseTimer(cmd, resource string, hasError bool) interfaces.ComponentTimer {
	return newHBaseTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		cmd,
		resource,
		hasError,
	)
}

func (c *components) NewRMQProduceTimer(exchange, resource string, hasError bool) interfaces.ComponentTimer {
	return newRMQProduceTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		exchange,
		resource,
		hasError,
	)
}

func (c *components) NewRMQConsumeCounter(queue, resource string) interfaces.ComponentCounter {
	return newRMQConsumeCounter(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		queue,
		resource,
	)
}

func (c *components) NewMongoTimer(command, collection, resource string, hasError bool) interfaces.ComponentTimer {
	return newMongoTimer(
		c.baseMeter,
		c.cfg.PrefixMetricName,
		command,
		collection,
		resource,
		hasError,
	)
}
