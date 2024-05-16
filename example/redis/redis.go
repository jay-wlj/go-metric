package main

import (
	"context"
	"time"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redismock/v9"
	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/interfaces"
)

// Reference
// https://redis.uptrace.dev/zh/guide/go-redis-hook.html

type hook struct {
	resource     string
	hasErrorFunc func(err error) bool
}

var _ redis.Hook = &hook{}

func NewHook(resource string) *hook {
	return &hook{
		resource: resource,
		hasErrorFunc: func(err error) bool {
			if err == redis.Nil {
				return false
			}
			if err == context.Canceled || err == context.DeadlineExceeded {
				return false
			}

			return err != nil
		},
	}
}

// DialHook: 当创建网络连接时调用的hook
func (h *hook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook: 执行命令时调用的hook
func (h *hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		begin := time.Now()
		var err = next(ctx, cmd)

		ret := "0"
		if h.hasErrorFunc(cmd.Err()) {
			ret = "1"
			fmt.Sprintf("ProcessHook cmd:%v cmdErr:%v err:%v\n", cmd.Name(), cmd.Err(), err)
		}
		gometric.GetGlobalMeter().NewTimer("hll_redis_request_duration_seconds").
			AddTag("cmd", cmd.Name()).
			AddTag("resource", h.resource).
			AddTag("error", ret).
			UpdateSince(begin)
		return err
	}
}

// ProcessPipelineHook: 执行管道命令时调用的hook
func (h *hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		begin := time.Now()
		var err = next(ctx, cmds)

		for _, cmd := range cmds {
			ret := "0"
			if h.hasErrorFunc(cmd.Err()) {
				ret = "1"
				fmt.Sprintf("ProcessPipelineHook cmd:%v cmdErr:%v err:%v\n", cmd.Name(), cmd.Err(), err)
			}

			gometric.GetGlobalMeter().NewTimer("hll_redis_request_duration_seconds").
				AddTag("cmd", cmd.Name()).
				AddTag("resource", h.resource).
				AddTag("error", ret).
				UpdateSince(begin)
		}

		return err
	}
}



func main() {
	ops := []interfaces.Option{
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
		gometric.WithPrometheusPort(16671),
		gometric.WithAppID("ci-hty"),
		gometric.WithAppVer("v1.21.0"),
		gometric.WithHistogramBoundaries([]float64{0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 10, 15}), // histogram默认分桶序列,单位:秒, 不要随便变动
		gometric.WithBaseLabelNames("dc_sname", "", "dc_ip", "", "dc_sver"),
	}

	// if pushCfg != nil && pushCfg.Endpoint != "" {
	// 	var period = time.Duration(pushCfg.Period) * time.Second
	// 	if period == 0 {
	// 		period = time.Duration(15 * time.Second)
	// 	}
	// 	ops = append(ops, gometric.WithPush(pushCfg.Endpoint, period))
	// }
	_ = gometric.NewMeter(
		ops...,
	)
	

	client, mock := redismock.NewClientMock()
	client.AddHook(
		// use "" when resource-id doesn't exist
		NewHook("ci-redis-resource1"),
	)
	mock.ExpectGet("1").SetVal("xxx")
	client.Get(context.TODO(), "1")
	mock.ExpectDel("2", "3").SetVal(2)
	client.Del(context.TODO(), "2", "3")

	// pipeline mode
	client2, mock2 := redismock.NewClientMock()
	client2.AddHook(
		NewHook("ci-redis-resource2"),
	)
	pipeline := client2.Pipeline()
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("%d", i)
		mock2.ExpectSet(key, key, time.Minute).SetVal(key)
	}
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("%d", i)
		pipeline.Set(context.TODO(), key, key, time.Minute)
	}
	_, _ = pipeline.Exec(context.TODO())
	time.Sleep(time.Minute)
}
