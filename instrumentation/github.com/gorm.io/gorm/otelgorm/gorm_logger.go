package otelgorm

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/jay-wlj/go-metric/internal/global"
)

type traceRecorder struct {
	logger.Interface
	resource     string
	hasErrorFunc func(err error) bool
}

func NewTraceRecorder(loggerI logger.Interface, resource string) *traceRecorder {
	return &traceRecorder{
		Interface:    loggerI,
		resource:     resource,
		hasErrorFunc: defaultHasErrorFunc,
	}
}

func (l *traceRecorder) WithHasErrorFunc(hasErrorFunc func(err error) bool) *traceRecorder {
	l.hasErrorFunc = hasErrorFunc
	return l
}

// LogMode log mode
func (l *traceRecorder) LogMode(level logger.LogLevel) logger.Interface {
	l.Interface = l.Interface.LogMode(level)
	return l
}

// Trace print sql message
func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.Interface.Trace(ctx, begin, fc, err)

	sql, _ := fc()

	defer func() {
		global.GetMeter().Components().
			NewMysqlTimer(sql, l.resource, l.hasErrorFunc(err)).
			UpdateSince(begin)
	}()
}

func defaultHasErrorFunc(err error) bool {
	if err == nil {
		return false
	}
	return !errors.Is(err, gorm.ErrRecordNotFound)
}
