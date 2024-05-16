package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redismock/v9"

	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/instrumentation/github.com/go-redis/redis/v8/otelredis"
)

// https://github.com/go-redis/redismock

func main() {
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
	)

	client, mock := redismock.NewClientMock()
	client.AddHook(
		// use "" when resource-id doesn't exist
		otelredis.NewHook("ci-redis-resource1"),
	)
	mock.ExpectGet("1").SetVal("xxx")
	client.Get(context.TODO(), "1")
	mock.ExpectDel("2", "3").SetVal(2)
	client.Del(context.TODO(), "2", "3")

	// pipeline mode
	client2, mock2 := redismock.NewClientMock()
	client2.AddHook(
		otelredis.NewHook("ci-redis-resource2"),
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
