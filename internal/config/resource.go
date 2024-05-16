package config

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"github.com/jay-wlj/go-metric/internal/labels"
)

// DtlResource 返回 dt 指标的公用tags
// 创建 newGauge、newCounter时 再次指定重名的tag即可覆盖
func DtlResource() *resource.Resource {
	// return resource.NewSchemaless(DtlLabels()...)
	rs, err := resource.New(context.Background(),
		// 流量考虑这些都屏蔽了
		// resource.WithProcess(),		// 是否带上进程标签 process_command_args="[./redis]",process_executable_name="redis",process_executable_path="/opt/src/go-metric/example/redis/redis",process_owner="jayden",process_pid="4858",process_runtime_description="go version go1.22.1 darwin/arm64",process_runtime_name="go",process_runtime_version="go1.22.1"
		// resource.WithOS(),			// 是否带上os标签 os_description="macOS 14.2.1 (23C71) (Darwin jaydendeMacBook-Pro.local 23.2.0 Darwin Kernel Version 23.2.0: Wed Nov 15 21:53:34 PST 2023; root:xnu-10002.61.3~2/RELEASE_ARM64_T8103 arm64)"
		// resource.WithHost(), 		// 是否带上host标签 host_name="jaydendeMacBook-Pro.local"
		// resource.WithHostID(),		// 是否带上host_id标签 host_id="63F3E324-ABBC-5EF4-B583-270E43258C00"
		// resource.WithFromEnv(),
		// resource.WithContainer(), 
		// resource.WithOSType(),		// 是否带上系统标签，如 os_type="darwin"
		// resource.WithTelemetrySDK(),	// 指标label是否包括otel的版本信息，telemetry_sdk_language="go",telemetry_sdk_name="opentelemetry",telemetry_sdk_version="1.19.0"
		resource.WithAttributes(DtlLabels()...))

	if err != nil {
		GetConfig().WriteErrorOrNot("failed to initialize Resource:  " + err.Error())
	}
	return rs
}

func DtlLabels() labels.Labels {

	baseLabel := GetConfig().BaseLabel

	var ret labels.Labels
	if baseLabel != nil {
		if baseLabel.Appid != "" {
			ret = append(ret, attribute.KeyValue{Key: attribute.Key(baseLabel.Appid), Value: attribute.StringValue(GetConfig().AppId)})
		}
		if baseLabel.AppVer != "" {
			ret = append(ret, attribute.KeyValue{Key: attribute.Key(baseLabel.AppVer), Value: attribute.StringValue(GetConfig().AppVer)})
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
			attribute.KeyValue{Key: "dc_sname", Value: attribute.StringValue(GetConfig().AppId)},
			attribute.KeyValue{Key: "dc_sver", Value: attribute.StringValue(GetConfig().AppVer)},
			attribute.KeyValue{Key: "dc_ip", Value: attribute.StringValue(GetConfig().LocalIP)},
			// attribute.KeyValue{Key: "dc_data_type", Value: attribute.StringValue("business")},
		}
	}

	return ret
}
