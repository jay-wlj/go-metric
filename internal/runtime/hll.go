package runtime

import "github.com/jay-wlj/go-metric/internal/config"

func dtlSystemNamespace(s string) string {
	return config.GetConfig().PrefixMetricName + s
}

func memstatNamespace(s string) string {
	return dtlSystemNamespace("go_memstats_" + s)
}
