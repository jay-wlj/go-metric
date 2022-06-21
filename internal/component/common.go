package component

import (
	"time"

	"github.com/jay-wlj/go-metric/interfaces"
)

const (
	// http server
	statusKey = "status"
	routeKey  = "route"
	retKey    = "ret"

	// http client
	fromAppIdKey    = "from_appid"
	toAppIdKey      = "to_appid"
	clientIPKey     = "client_ip"
	serverDomainKey = "server_domain"
	serverURLApiKey = "server_url_api"
	errKey          = "error"
	// retKey
	// statusKey

	// mysql
	cmdKey      = "cmd"
	sqlKey      = "sql"
	resourceKey = "resource"
	// errKey

	// redis
	// cmdKey
	// resourceKey
	// errKey

	// kafka
	topicKey = "topic"
	// resourceKey
	// errKey

	// mq
	exchangeKey = "exchange"
	queueKey    = "queue"
	// resourceKey
	// errKey

	// es
	apiKey   = "api"
	indexKey = "index"
	// resourceKey
	// errKey

	// mongo
	commandKey    = "command"
	collectionKey = "collection"
	// resourceKey
	// errKey
)

func newComponentTimer(timer interfaces.Timer) *componentTimer {
	timer.AddTag("dtl_data_type", "base")
	return &componentTimer{timer: timer}
}

type componentTimer struct {
	timer interfaces.Timer
}

func (holder *componentTimer) Update(d time.Duration)      { holder.timer.Update(d) }
func (holder *componentTimer) UpdateSince(start time.Time) { holder.timer.UpdateSince(start) }
func (holder *componentTimer) UpdateInMillis(m float64)    { holder.timer.UpdateInMillis(m) }
func (holder *componentTimer) UpdateInSeconds(s float64)   { holder.timer.UpdateInSeconds(s) }

func newComponentCounter(counter interfaces.Counter) interfaces.ComponentCounter {
	counter.AddTag("dtl_data_type", "base")
	return &componentCounter{counter: counter}
}

type componentCounter struct {
	counter interfaces.Counter
}

func (holder *componentCounter) Incr(delta float64) { holder.counter.Incr(delta) }
func (holder *componentCounter) IncrOnce()          { holder.counter.IncrOnce() }
