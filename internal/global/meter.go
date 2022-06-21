package global

import (
	"sync/atomic"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/meter/nop"
)

type meterHolder struct {
	meter interfaces.Meter
}

var (
	globalMeter = atomic.Value{}
)

func init() {
	SetNopMeter()
}

func GetMeter() interfaces.Meter {
	return globalMeter.Load().(meterHolder).meter
}

func SetNopMeter() {
	globalMeter.Store(meterHolder{meter: &nop.Meter{}})
}

func SetMeter(m interfaces.Meter) {
	if m == nil {
		return
	}
	globalMeter.Store(meterHolder{
		meter: m,
	})
}
