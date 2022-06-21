package prom

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/jay-wlj/go-metric/internal/config"
)

func Test_PromMeter_ExceedCount(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randName := func() string {
		rand.Seed(time.Now().UnixNano())
		return strconv.Itoa(rand.Intn(1000))
	}
	m, _ := NewPrometheusMeter(&config.Config{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				m.NewCounter(randName()+"count").AddTag("tag", randName()).IncrOnce()
				m.NewGauge(randName()+"gauge").AddTag("tag", randName()).Update(1)
				m.NewTimer(randName()+"histogram").AddTag("tag", randName()).Update(time.Second)
			}
		}()
	}

	time.Sleep(time.Second * 30)
}
