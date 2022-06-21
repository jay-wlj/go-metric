package main

import (
	"fmt"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/interfaces"
)

var count int64

func main() {
	runtime.GOMAXPROCS(2)
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
		gometric.WithRuntimeStatsCollector(),
		gometric.WithConsul("consul-dev.my.cn:8500", "d21eb730-b3a0-57af-2b21-efe123f0fd53"),
		gometric.WithAppID("demo"),
		gometric.WithEnv("dev"),
	)

	meter := gometric.GetGlobalMeter()

	go daemonCounter(meter)
	go daemonGauge(meter)
	go daemonHistogram(meter)

	go speedPrint()

	time.AfterFunc(time.Second*30, func() {
		meter.WithRunning(false)
	})
	time.AfterFunc(time.Second*60, func() {
		meter.WithRunning(true)
	})

	select {}
}

func randFloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Intn(10000))
}

func randIP() string {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(10)
	return fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)
}

func speedPrint() {
	ticker := time.NewTicker(time.Second * 10)
	var lastCount int64
	for {
		lastCount = atomic.LoadInt64(&count)
		select {
		case <-ticker.C:
			fmt.Println("[go-metric][info]:", (atomic.LoadInt64(&count)-lastCount)/10, "qps")
		}
	}
}

func daemonCounter(m interfaces.Meter) {
	hostname, _ := os.Hostname()

	for x := 0; x < 10; x++ {
		go func(x int) {
			for {
				m.NewCounter("counter"+strconv.Itoa(x)).
					AddTag("host", hostname).
					AddTag("dtl_data_type", "base").
					AddTag("zone", "sh").
					AddTag("ip", randIP()).
					IncrOnce()
				atomic.AddInt64(&count, 1)
			}
		}(x)
	}
}

func daemonGauge(m interfaces.Meter) {
	hostname, _ := os.Hostname()
	for x := 0; x < 10; x++ {
		go func(x int) {
			for {
				m.NewGauge("gauge"+strconv.Itoa(x)).
					AddTag("host", hostname).
					AddTag("zone", "sh").
					AddTag("ip", randIP()).
					Update(randFloat())
				atomic.AddInt64(&count, 1)
			}
		}(x)
	}
}

func daemonHistogram(m interfaces.Meter) {
	hostname, _ := os.Hostname()

	for x := 0; x < 10; x++ {
		go func(x int) {
			for {
				m.NewTimer("timer"+strconv.Itoa(x)).
					AddTag("host", hostname).
					AddTag("zone", "sh").
					AddTag("ip", randIP()).
					Update(time.Second)
				atomic.AddInt64(&count, 1)
			}
		}(x)
	}
}
