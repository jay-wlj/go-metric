package gometric

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/global"
	"github.com/jay-wlj/go-metric/internal/labels"
)

// GetGlobalMeter returns the registered global meter.
// If none is registered then a nop meter is returned.
func GetGlobalMeter() interfaces.Meter {
	return global.GetMeter()
}

// Filter extract the common pattern for http route, sql and so on
var (
	Filter = labels.Filter
	_      = Filter // suppress waring
)
