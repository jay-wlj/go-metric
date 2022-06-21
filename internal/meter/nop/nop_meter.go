package nop

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
	"github.com/jay-wlj/go-metric/internal/metrics/nop"
)

const reportIntervalThreshold = time.Second * 10

var _ interfaces.Meter = &Meter{}

type Meter struct {
	cfg          *config.Config
	cause        error
	lastReportAt int64 // unix seconds
	locker       sync.Mutex
}

func NewNopMeter(cfg *config.Config, cause error) *Meter {
	nm := &Meter{
		cfg:   cfg,
		cause: cause,
	}
	nm.checkOrReportError()
	return nm
}

// 错误抑制，最多 5 秒钟报错一次
func (nm *Meter) checkOrReportError() {
	shouldReport := nm.shouldReport()
	if !shouldReport {
		return
	}

	nm.locker.Lock()
	stillShouldReport := nm.shouldReport()
	if stillShouldReport {
		atomic.StoreInt64(&nm.lastReportAt, time.Now().Unix())
	}
	nm.locker.Unlock()

	if stillShouldReport {
		if nm.cause == nil {
			return
		}
		if nm.cfg == nil {
			return
		}
		nm.cfg.WriteErrorOrNot(fmt.Sprintf(
			"you are using nopMeter now! cause : %s", nm.cause))
	}
}

func (nm *Meter) shouldReport() bool {
	lastReportTime := atomic.LoadInt64(&nm.lastReportAt)
	time.Unix(lastReportTime, 0)
	return time.Now().After(time.Unix(lastReportTime, 0).Add(reportIntervalThreshold))
}

func (nm *Meter) WithRunning(_ bool) {
	nm.checkOrReportError()
}

func (nm *Meter) NewCounter(_ string) interfaces.Counter {
	nm.checkOrReportError()
	return &nop.Counter
}

func (nm *Meter) NewGauge(_ string) interfaces.Gauge {
	nm.checkOrReportError()
	return &nop.Gauge
}

func (nm *Meter) NewTimer(_ string) interfaces.Timer {
	nm.checkOrReportError()
	return &nop.Timer
}

func (nm *Meter) Components() interfaces.Components {
	return &components{}
}
