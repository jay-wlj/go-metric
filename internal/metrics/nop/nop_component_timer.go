package nop

import (
	"time"

	"github.com/jay-wlj/go-metric/interfaces"
)

var _ interfaces.ComponentTimer = &nopComponentTimer{}
var ComponentTimer = &nopComponentTimer{}

type nopComponentTimer struct{}

func (nt *nopComponentTimer) Update(_ time.Duration)    {}
func (nt *nopComponentTimer) UpdateSince(_ time.Time)   {}
func (nt *nopComponentTimer) UpdateInMillis(_ float64)  {}
func (nt *nopComponentTimer) UpdateInSeconds(_ float64) {}
