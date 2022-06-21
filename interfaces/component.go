package interfaces

import (
	"net/http"
	"time"
)

// ComponentTimer : 中间件的 timer 函数，没有拓展tags的功能
type ComponentTimer interface {
	Update(d time.Duration)      // 记录一段耗时
	UpdateSince(start time.Time) // 记录从起始时间的耗时
	UpdateInMillis(m float64)    // 记录一段毫秒单位的耗时
	UpdateInSeconds(s float64)   // 记录一段秒单位的耗时
}

// ComponentCounter : 中间件的 counter 函数，没有拓展tags的功能
type ComponentCounter interface {
	Incr(delta float64) // Incr(1)
	IncrOnce()          // +1
}

// Components 用于封装中间件的埋点
type Components interface {
	// NewHTTPServerTimer HTTP请求服务端埋点 （当前服务被调的埋点）
	// dtlci_api_request_seconds
	//        path  (*http.Request).URL.Path
	//        ret    响应 body 的ret码, 非数字时统一为 '-'
	//        status (*http.Response).StatusCode
	NewHTTPServerTimer(path string, ret string, statusCode int) ComponentTimer

	// NewHTTPClientTimer : HTTP请求客户端埋点 （当前服务调用下游服务时埋点）
	// dtlci_service_http_call_seconds
	//        from_appid，环境中获取
	//        to_appid, 根据serverDomain 从lalaplat2获取
	//        client_ip，环境中获取
	//        server_domain，(*http.Request).Host
	//        server_url_api, (*http.Request).URL.Path
	//        error，自动判定 状态码 >= 400, <=600 为 "1"，否则为 "0"
	//        ret
	//        status (*http.Response).StatusCode
	// ---------------------------
	// serverDomain 被调服务的 domain, 领域域名
	// serverPath 即 被调服务的 path
	// serverRet 即响应b body 的 ret,  非数字时统一为 '-'
	// status 响应的http 状态
	NewHTTPClientTimer(serverDomain, serverPath string, serverRet string, statusCode int) ComponentTimer
	// NewHTTPClientTimerFromResponse 与 NewHTTPClientTimer 类似，只是简化了参数
	// 为避免在SDK内使用 ioutil.ReadAll，用户需要自行解析 resp.Body 显式的传入 serverRet
	NewHTTPClientTimerFromResponse(resp *http.Response, serverRet string) ComponentTimer

	// NewMysqlTimer : mysql埋点
	// dtlci_mysql_request_seconds
	//        cmd，select/insert/delete/update
	//        resource, 资源
	//        sql, 小写的sql串，变量被占位符填充
	//        error, "1" / "0"
	// ---------------------------
	// sql： 原始sql，会自动解析出 cmd
	// resource: 资源名
	// hasError: 是否有错误
	NewMysqlTimer(sql string, resource string, hasError bool) ComponentTimer

	// NewRedisTimer : redis 埋点
	// dtlci_redis_request_seconds
	//        cmd，get / set/ del
	//        resource, 资源
	//        error, "1" / "0"
	// ---------------------------
	// cmd： redis command
	// resource: 资源名
	// hasError: 是否有错误
	NewRedisTimer(cmd string, resource string, hasError bool) ComponentTimer

	// NewESTimer elastic-search 请求埋点
	// dtlci_es_request_seconds
	//        api, es api
	//        index, 索引
	//        resource, 资源
	//        error, "1" / "0"
	// ---------------------------
	// api: es api
	// index: 索引
	// resource: 资源名
	// hasError: 是否有错误
	NewESTimer(api, index, resource string, hasError bool) ComponentTimer

	// NewHBaseTimer hbase 请求埋点
	// dtlci_hbase_request_seconds
	//        cmd, hbase 命令
	//        resource, 资源
	//        error, "1" / "0"
	// ---------------------------
	// cmd: hbase 命令
	// resource: 资源名
	// hasError: 是否有错误
	NewHBaseTimer(cmd, resource string, hasError bool) ComponentTimer

	// NewRMQProduceTimer rmq 生产埋点
	// dtlci_rabbit_producer_seconds histogram
	//        exchange, 交换机名称
	//        resource, 资源
	//        error, "1" / "0"
	// ---------------------------
	// exchange: 交换机
	// resource: 资源名
	// hasError: 是否有错误
	NewRMQProduceTimer(exchange, resource string, hasError bool) ComponentTimer

	// NewRMQConsumeCounter rmq 消费埋点
	// dtlci_rabbit_consumer_total counter
	//        queue, 队列
	//        resource, 资源
	// ---------------------------
	// queue: 队列
	// resource: 资源名
	NewRMQConsumeCounter(queue, resource string) ComponentCounter

	// NewMongoTimer mongo 请求埋点
	// dtlci_mongo_request_seconds
	//        command, mongo 命令, find/insert/update/delete等
	//        collection, 集合
	//        resource，资源
	//        error, "1" / "0"
	// ---------------------------
	// command: mongo command
	// collection: 集合
	// resource: 资源名
	// hasError: 是否有错误
	NewMongoTimer(command, collection, resource string, hasError bool) ComponentTimer

	// NewKafkaProduceTimer : kafka 生产埋点
	// dtlci_kafka_producer_seconds
	//        topic,
	//        resource, 资源
	//        error, "1" / "0"
	// ---------------------------
	// topic： kafka topic
	// resource: 资源名
	// hasError: 是否有错误
	NewKafkaProduceTimer(topic string, resource string, hasError bool) ComponentTimer

	// NewKafkaConsumeCounter : kafka 消费埋点
	// dtlci_kafka_consumer_total
	//        topic,
	//        resource, 资源
	// ---------------------------
	// topic： kafka topic
	// resource: 资源名
	NewKafkaConsumeCounter(topic string, resource string) ComponentCounter
}
