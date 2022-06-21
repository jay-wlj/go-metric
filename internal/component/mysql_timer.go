package component

import (
	"github.com/jay-wlj/go-metric/interfaces"
	"github.com/jay-wlj/go-metric/internal/labels"
)

const mysqlMetricName = "dtlci_mysql_request_seconds"

func newMysqlTimer(
	meter interfaces.BaseMeter,
	sql string,
	resource string,
	hasError bool,
) interfaces.ComponentTimer {
	timer := meter.NewTimer(mysqlMetricName)
	cmd, sql, ok := labels.Filter.FilterSQL(sql)
	// TODO
	// 由于go 中无法 prepare sql，
	// 导致打了大量的原始sql，为了避免给存储带来压力，不再打原始sql
	if ok {
		timer.AddTag(cmdKey, cmd)
		timer.AddTag(sqlKey, sql)
	} else {
		timer.AddTag(cmdKey, "-")
		timer.AddTag(sqlKey, "-")
	}
	timer.AddTag(resourceKey, labels.Filter.FilterResource(resource))
	if hasError {
		timer.AddTag(errKey, "1")
	} else {
		timer.AddTag(errKey, "0")
	}
	return newComponentTimer(timer)
}
