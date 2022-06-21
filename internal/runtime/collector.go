package runtime

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/config"
)

const defaultRuntimeCollectInterval = time.Second * 10

type Collector interface {
	Start()
	Stop()
}

type collector struct {
	cfg     *config.Config
	meter   interfaces.BaseMeter
	running int32
	closeCh chan struct{}
	// runtime cached info
	msLast          *runtime.MemStats
	msLastTimestamp time.Time
	msMtx           sync.Mutex
	msMaxWait       time.Duration
	msMaxAge        time.Duration
}

func NewCollector(
	cfg *config.Config,
	meter interfaces.BaseMeter,
) Collector {
	return &collector{
		cfg:     cfg,
		meter:   meter,
		running: 0,
		closeCh: make(chan struct{}),
	}
}

func (c *collector) Start() {
	if !c.cfg.ReadRuntimeStats {
		c.cfg.WriteInfoOrNot("WithRuntimeStatsCollector=false, runtime metrics collector is disabled")
		return
	}
	c.cfg.WriteInfoOrNot("WithRuntimeStatsCollector=true, runtime metrics collector is enabled")
	// another goroutine is already running
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		return
	}
	go c.Collector()
}

func (c *collector) Collector() {
	c.cfg.WriteInfoOrNot("runtime metrics collector is running")
	ticker := time.NewTicker(defaultRuntimeCollectInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.closeCh:
			c.msLast = nil
			return
		case <-ticker.C:
			c.collectMemStats()
			c.collectThreads()
		}
	}
}

func (c *collector) Stop() {
	// collector is already stopped
	if !(atomic.CompareAndSwapInt32(&c.running, 1, 0)) {
		return
	}
	c.closeCh <- struct{}{}
	c.cfg.WriteInfoOrNot("runtime metrics collector is stopped")
}

func (c *collector) newGaugeWithTags(metricName string) interfaces.Gauge {
	return c.meter.NewGauge(metricName).AddTag("dtl_data_type", "base")
}

func (c *collector) newCounterWithTags(metricName string) interfaces.Counter {
	return c.meter.NewCounter(metricName).AddTag("dtl_data_type", "base")
}

func (c *collector) collectMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	defer func() {
		c.msLast = &ms
	}()
	if c.msLast == nil {
		return
	}
	// Number of bytes allocated and still in use
	c.newGaugeWithTags(memstatNamespace("alloc_bytes")).
		Update(float64(ms.Alloc))
	// Total number of bytes allocated, even if freed
	c.newCounterWithTags(memstatNamespace("alloc_bytes_total")).
		Incr(float64(ms.TotalAlloc - c.msLast.TotalAlloc))
	// Number of bytes obtained from system
	c.newGaugeWithTags(memstatNamespace("sys_bytes")).
		Update(float64(ms.Sys))
	// Total number of pointer lookups
	c.newCounterWithTags(memstatNamespace("lookups_total")).
		Incr(float64(ms.Lookups - c.msLast.Lookups))
	// Total number of mallocs
	c.newCounterWithTags(memstatNamespace("mallocs_total")).
		Incr(float64(ms.Mallocs - c.msLast.Mallocs))
	// Total number of frees
	c.newCounterWithTags(memstatNamespace("frees_total")).
		Incr(float64(ms.Frees - c.msLast.Frees))
	// Number of heap bytes allocated and still in use
	c.newGaugeWithTags(memstatNamespace("heap_alloc_bytes")).
		Update(float64(ms.HeapAlloc))
	// Number of heap bytes obtained from system
	c.newGaugeWithTags(memstatNamespace("heap_sys_bytes")).
		Update(float64(ms.HeapSys))
	// Number of heap bytes waiting to be used
	c.newGaugeWithTags(memstatNamespace("heap_idle_bytes")).
		Update(float64(ms.HeapIdle))
	// Number of heap bytes that are in use.
	c.newGaugeWithTags(memstatNamespace("heap_inuse_bytes")).
		Update(float64(ms.HeapInuse))
	// Number of heap bytes released to OS.
	c.newGaugeWithTags(memstatNamespace("heap_released_bytes")).
		Update(float64(ms.HeapReleased))
	// Number of allocated objects
	c.newGaugeWithTags(memstatNamespace("heap_objects")).
		Update(float64(ms.HeapObjects))
	// Number of bytes in use by the stack allocator
	c.newGaugeWithTags(memstatNamespace("stack_inuse_bytes")).
		Update(float64(ms.StackInuse))
	// Number of bytes obtained from system for stack allocator.
	c.newGaugeWithTags(memstatNamespace("stack_sys_bytes")).
		Update(float64(ms.StackSys))
	// Number of bytes in use by mspan structures
	c.newGaugeWithTags(memstatNamespace("mspan_inuse_bytes")).
		Update(float64(ms.MSpanInuse))
	// Number of bytes used for mspan structures obtained from system
	c.newGaugeWithTags(memstatNamespace("mspan_sys_bytes")).
		Update(float64(ms.MSpanSys))
	c.newGaugeWithTags(memstatNamespace("mcache_inuse_bytes")).
		Update(float64(ms.MCacheInuse))
	// Number of bytes used for mcache structures obtained from system.
	c.newGaugeWithTags(memstatNamespace("mcache_sys_bytes")).
		Update(float64(ms.MCacheSys))
	// Number of bytes used by the profiling bucket hash table
	c.newGaugeWithTags(memstatNamespace("buck_hash_sys_bytes")).
		Update(float64(ms.BuckHashSys))
	// Number of bytes used for garbage collection system metadata.
	c.newGaugeWithTags(memstatNamespace("gc_sys_bytes")).
		Update(float64(ms.GCSys))
	// Number of bytes used for other system allocations
	c.newGaugeWithTags(memstatNamespace("other_sys_bytes")).
		Update(float64(ms.OtherSys))
	// Number of heap bytes when next garbage collection will take place
	c.newGaugeWithTags(memstatNamespace("next_gc_bytes")).
		Update(float64(ms.NextGC))
	// Number of seconds since 1970 of last garbage collection.
	c.newGaugeWithTags(memstatNamespace("last_gc_time_seconds")).
		Update(float64(ms.LastGC))
	// PauseNs is a circular buffer of recent GC stop-the-world
	// pause times in nanoseconds.
	c.newGaugeWithTags(memstatNamespace("pause_ns")).
		Update(float64(ms.PauseNs[(ms.NumGC+255)%256]))
	// PauseTotalNs is the cumulative nanoseconds in GC
	// stop-the-world pauses since the program started.
	c.newCounterWithTags(memstatNamespace("pause_total_ns")).
		Incr(float64(ms.PauseTotalNs - c.msLast.PauseTotalNs))
	// NumGC is the number of completed GC cycles.
	c.newCounterWithTags(memstatNamespace("num_gc")).
		Incr(float64(ms.NumGC - c.msLast.NumGC))
	// The fraction of this program's available CPU time used by the GC since the program started
	c.newGaugeWithTags(memstatNamespace("gc_cpu_fraction")).
		Update(ms.GCCPUFraction)
}

func (c *collector) collectThreads() {
	// Number of goroutines that currently exist
	c.newGaugeWithTags(dtlSystemNamespace("go_goroutines")).
		Update(float64(runtime.NumGoroutine()))
	// Number of OS threads created.
	n, _ := runtime.ThreadCreateProfile(nil)
	c.newGaugeWithTags(dtlSystemNamespace("go_threads")).
		Update(float64(n))
	// ignore gc duration summary
	// ignore go version, unnecessary information
}
