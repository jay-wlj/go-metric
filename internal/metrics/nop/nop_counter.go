package nop

import "github.com/jay-wlj/go-metric/interfaces"

var _ interfaces.Counter = &nopCounter{}
var Counter = nopCounter{}

type nopCounter struct{}

func (nc *nopCounter) Incr(_ float64)                                  {}
func (nc *nopCounter) IncrOnce()                                       {}
func (nc *nopCounter) AddTag(_, _ string) interfaces.Counter           { return nc }
func (nc *nopCounter) WithTags(_ map[string]string) interfaces.Counter { return nc }
