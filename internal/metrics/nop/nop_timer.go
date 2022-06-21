package nop

import (
	"time"

	"github.com/jay-wlj/go-metric/interfaces"
)

var _ interfaces.Timer = &nopTimer{}
var Timer = nopTimer{}

type nopTimer struct{}

func (nt *nopTimer) Time(f func())                                 { f() }
func (nt *nopTimer) Update(_ time.Duration)                        {}
func (nt *nopTimer) UpdateSince(_ time.Time)                       {}
func (nt *nopTimer) UpdateInMillis(_ float64)                      {}
func (nt *nopTimer) UpdateInSeconds(_ float64)                     {}
func (nt *nopTimer) AddTag(_, _ string) interfaces.Timer           { return nt }
func (nt *nopTimer) WithTags(_ map[string]string) interfaces.Timer { return nt }
