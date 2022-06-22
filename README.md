# README

golang 监控 SDK, 基于 open-telemetry SDK

(目前仅支持指标，trace功能根据用户需求决定是否开发)

详细文档参考



```golang

go get github.com/jay-wlj/go-metric


```



## 初始化
### case1，使用 prometheus 进行初始化, 默认使用标准输出打印 info 和 error 日志
不使用 Option 显示初始化时，配置会尝试从以下环境变量中获取
- dtl.consul.host: consul host;
- dtl.consul.port: consul 端口；
- dtl.consul.token: consul token；
- dtl.monitor.port: prometheus 端口
- dtl.app.id: appid

```golang
_ = gometric.NewMeter(
   gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
   gometric.WithConsul("consul-dev.xxx.cn:8500", "d21eb730-*********"),
   gometric.WithAppID("xx-xx-demo"),
   gometric.WithEnv("dev"),
   gometric.WithPrometheusPort(12345),
)
meter := gometric.GetGlobalMeter()
meter.NewCounter("counter").
    AddTag("host", hostname).
    AddTag("zone", "sh").
    AddTag("ip", "1.1.1.1").
    IncrOnce()
```
### case2, 使用 zap、logrus 等
```golang
m := gometric.NewMeter(
    gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
    gometric.WithInfoLogWrite(func(infoLine string) {
    	logger.info(infoLine)
    }),
    gometric.WithErrorLogWrite(func(errorLine string) {
        logger.error(infoLine)
    }),
    gometric.WithRuntimeStatsCollector(), // 采集runtime指标
)
m.NewGauge("gauge1").
    WithTags(map[string]string{
        "host": hostname,
        "zone": "sh",
        "ip": "1.1.1.1",
    }).Update(3.0)
```
## 框架埋点
#### http server
gin参考 example/gin/main.go

通过 otelgin.HTTPServerTimerMiddleware() 创建中间件

#### http client

参考 example/httpclient/main.go

(需要在代码内主动埋点，替换 原生http 库的方式成本较高，待引入trace后再考虑)

#### go-redis
参考 example/go-redis/v8/main.go

通过 otelredis.NewHook("ci-redis-resource1") 创建钩子函数

没有resource则留空, 异常判定可用 WithHasErrorFunc 更改

#### gorm
参考 example/gorm/main.go

通过 otelgorm.NewTraceRecorder(logger.Default, "ci-mysql-resource") 创建 gorm logger，

没有resource则留空, 异常判定可用 WithHasErrorFunc 更改

## 其他框架埋点
以上未登记的框架，可以联系监控组迭代支持，也可自行埋点。
暂时开放的简化 api 见 interface.Components 文件

## 限制
- tags key 不能以 __ 双下划线开头
必须是 "^[a-zA-Z_][a-zA-Z0-9_]*$"
- metric 数量限制为 150 条
- 单个metric的最大组合数为 300 
