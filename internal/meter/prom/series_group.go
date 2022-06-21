package prom

import (
	"sync"
)

type seriesGroup struct {
	series map[uint64]struct{}
	lock   sync.Mutex
}

func newSeriesGroup(seriesID uint64) *seriesGroup {
	return &seriesGroup{
		series: map[uint64]struct{}{seriesID: {}},
	}
}

func (sg *seriesGroup) ExceedThreshold(seriesID uint64) bool {
	sg.lock.Lock()
	defer sg.lock.Unlock()
	_, exist := sg.series[seriesID]
	if exist {
		return false
	}
	if len(sg.series) >= maxSeriesCount {
		return true
	}
	sg.series[seriesID] = struct{}{}
	return false
}
