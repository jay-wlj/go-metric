package nop

import "github.com/jay-wlj/go-metric/interfaces"

var _ interfaces.Gauge = &nopGauge{}
var Gauge = nopGauge{}

type nopGauge struct{}

func (ng *nopGauge) Update(_ float64)                              {}
func (ng *nopGauge) AddTag(_, _ string) interfaces.Gauge           { return ng }
func (ng *nopGauge) WithTags(_ map[string]string) interfaces.Gauge { return ng }
