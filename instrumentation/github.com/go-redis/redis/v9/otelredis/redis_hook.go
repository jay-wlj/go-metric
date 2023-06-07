package otelredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/jay-wlj/go-metric/internal/global"
)

// Reference
// https://github.com/go-redis/redis/blob/v8.0.0-beta.5/redisext/otel.go

type hook struct {
	resource     string
	hasErrorFunc func(err error) bool
}

var _ redis.Hook = &hook{}

const startTimeKey = "startTime"

func NewHook(resource string) *hook {
	return &hook{
		resource:     resource,
		hasErrorFunc: defaultHasErrorFunc,
	}
}

func (h *hook) WithHasErrorFunc(hasErrorFunc func(err error) bool) redis.Hook {
	h.hasErrorFunc = hasErrorFunc
	return h
}

func (h *hook) BeforeProcess(ctx context.Context, _ redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTimeKey, time.Now()), nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	st, ok := ctx.Value(startTimeKey).(time.Time)
	if !ok {
		return nil
	}
	global.GetMeter().
		Components().
		NewRedisTimer(cmd.Name(), h.resource, h.hasErrorFunc(cmd.Err())).
		UpdateSince(st)
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTimeKey, time.Now()), nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	st, ok := ctx.Value(startTimeKey).(time.Time)
	if !ok {
		return nil
	}
	const numCmdLimit = 100
	for i, cmd := range cmds {
		if i >= numCmdLimit {
			break
		}
		global.GetMeter().
			Components().
			NewRedisTimer(cmd.Name(), h.resource, h.hasErrorFunc(cmd.Err())).
			UpdateSince(st)
	}

	return nil
}

func defaultHasErrorFunc(err error) bool {
	if err == redis.Nil {
		return false
	}
	return err != nil
}
