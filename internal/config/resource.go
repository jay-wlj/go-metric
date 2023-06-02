package config

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/jay-wlj/go-metric/internal/labels"
)

// DtlResource 返回 dt 指标的公用tags
// 创建 newGauge、newCounter时 再次指定重名的tag即可覆盖
func DtlResource() *resource.Resource {
	return resource.NewSchemaless(DtlLabels()...)
}

func DtlLabels() labels.Labels {

	baseLabel := GetConfig().BaseLabel

	var ret labels.Labels
	if baseLabel != nil {
		if baseLabel.Appid != "" {
			ret = append(ret, attribute.KeyValue{Key: attribute.Key(baseLabel.Appid), Value: attribute.StringValue(GetConfig().AppId)})
		}
		if baseLabel.Env != "" {
			ret = append(ret, attribute.KeyValue{Key: attribute.Key(baseLabel.Env), Value: attribute.StringValue(GetConfig().GetEnv())})
		}
		if baseLabel.IP != "" {
			ret = append(ret, attribute.KeyValue{Key: attribute.Key(baseLabel.IP), Value: attribute.StringValue(GetConfig().LocalIP)})
		}
		if baseLabel.DataType != "" {
			ret = append(ret, attribute.KeyValue{Key: attribute.Key(baseLabel.DataType), Value: attribute.StringValue("business")})
		}
	} else {
		// the default base label names
		ret = labels.Labels{
			attribute.KeyValue{Key: "dtl_appid", Value: attribute.StringValue(GetConfig().AppId)},
			attribute.KeyValue{Key: "dtl_env", Value: attribute.StringValue(GetConfig().GetEnv())},
			attribute.KeyValue{Key: "dtl_ip", Value: attribute.StringValue(GetConfig().LocalIP)},
			attribute.KeyValue{Key: "dtl_data_type", Value: attribute.StringValue("business")},
		}
	}

	return ret
}
