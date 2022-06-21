package middleware

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	gometric "github.com/jay-wlj/go-metric"
	"gorm.io/gorm"
)

type enhancedorm struct {
	orm.Ormer
	sql          string
	resource     string
	hasErrorFunc func(err error) bool
}

func Wraporm(ormI orm.Ormer, sql string, resource string) *enhancedorm {
	return &enhancedorm{
		Ormer:        ormI,
		sql:          sql,
		resource:     resource,
		hasErrorFunc: defaultHasErrorFunc,
	}
}

func (o *enhancedorm) WithHasErrorFunc(hasErrorFunc func(err error) bool) *enhancedorm {
	o.hasErrorFunc = hasErrorFunc
	return o
}

func defaultHasErrorFunc(err error) bool {
	if err == nil {
		return false
	}
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (o *enhancedorm) Read(md interface{}, cols ...string) error {
	start := time.Now()
	err := o.Ormer.Read(md, cols...)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"SELECT -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return err
}

func (o *enhancedorm) ReadForUpdate(md interface{}, cols ...string) error {
	start := time.Now()
	err := o.Ormer.ReadForUpdate(md, cols...)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"SELECT -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return err
}

func (o *enhancedorm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	start := time.Now()
	boolean, id, err := o.Ormer.ReadOrCreate(md, col1, cols...)
	i_u := "SELECT -"
	if boolean {
		i_u = "INSERT -"
	}
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		i_u,
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return boolean, id, err
}

func (o *enhancedorm) Insert(md interface{}) (int64, error) {
	start := time.Now()
	id, err := o.Ormer.Insert(md)
	gometric.GetGlobalMeter().Components().
		NewMysqlTimer("INSERT -", o.resource, o.hasErrorFunc(err)).
		UpdateSince(start)
	return id, err
}

func (o *enhancedorm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	start := time.Now()
	id, err := o.Ormer.InsertOrUpdate(md, colConflitAndArgs...)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"INSERT -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return id, err
}

func (o *enhancedorm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	start := time.Now()
	id, err := o.Ormer.InsertMulti(bulk, mds)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"INSERT -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return id, err
}

func (o *enhancedorm) Update(md interface{}, cols ...string) (int64, error) {
	start := time.Now()
	id, err := o.Ormer.Update(md, cols...)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"UPDATE -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)

	return id, err
}

func (o *enhancedorm) Delete(md interface{}, cols ...string) (int64, error) {
	start := time.Now()
	id, err := o.Ormer.Delete(md, cols...)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"DELETE -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return id, err
}

func (o *enhancedorm) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	start := time.Now()
	id, err := o.Ormer.LoadRelated(md, name, args...)
	gometric.GetGlobalMeter().Components().NewMysqlTimer(
		"SELECT -",
		o.resource,
		o.hasErrorFunc(err),
	).UpdateSince(start)
	return id, err
}

func (o *enhancedorm) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	return o.Ormer.QueryM2M(md, name)
}

func (o *enhancedorm) QueryTable(ptrStructOrTableName interface{}) (qs orm.QuerySeter) {
	return o.Ormer.QueryTable(ptrStructOrTableName)
}

func (o *enhancedorm) Using(name string) error {
	return o.Ormer.Using(name)
}

func (o *enhancedorm) Begin() error {
	return o.Ormer.Begin()
}

func (o *enhancedorm) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	return o.Ormer.BeginTx(ctx, opts)
}

func (o *enhancedorm) Commit() error {
	return o.Ormer.Commit()
}

func (o *enhancedorm) Rollback() error {
	return o.Ormer.Rollback()
}

func (o *enhancedorm) Raw(query string, args ...interface{}) orm.RawSeter {
	return o.Ormer.Raw(query, args...)
}

func (o *enhancedorm) Driver() orm.Driver {
	return o.Ormer.Driver()
}

func (o *enhancedorm) DBStats() *sql.DBStats {
	return o.Ormer.DBStats()
}
