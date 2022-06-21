package nop

import (
	"github.com/jay-wlj/go-metric/interfaces"
)

var _ interfaces.ComponentCounter = &nopComponentCounter{}
var ComponentCounter = &nopComponentCounter{}

type nopComponentCounter struct{}

func (nt *nopComponentCounter) Incr(_ float64) {}
func (nt *nopComponentCounter) IncrOnce()      {}
