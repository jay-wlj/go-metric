package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {

	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
		gometric.WithPush("10.10.21.137:7073", 15*time.Second),
		gometric.WithPrometheusPort(0),
		gometric.WithAppID("watermark_server"),
		gometric.WithPrefixBaseLabelName("dtl_"),
		gometric.WithPrefixMetricName("hll_"),
	)
	r := gin.New()
	r.Use(otelgin.HTTPServerTimerMiddleware())

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		gometric.GetGlobalMeter().Components().
			NewKafkaProduceTimer("abc", "", true).
			UpdateSince(time.Now())

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"ret":  0,
			"name": "unknown",
			"id":   id,
		})
	})
	r.GET("/metrics", func(c *gin.Context) {
		h := gometric.GetGlobalMeter().GetHandler()
		h.ServeHTTP(c.Writer, c.Request)
	})

	go func() {
		services := []string{"config_server"}
		for i := 0; i < 100; i++ {
			services = append(services, "config_server"+strconv.Itoa(i))
		}
		var method []string
		for i := 0; i < 1000; i++ {
			method = append(method, "grpc"+strconv.Itoa(i))
		}
		var rpcCode []string
		for i := 0; i < 100; i++ {
			rpcCode = append(rpcCode, strconv.Itoa(i))
		}
		var codes []string
		for i := 0; i < 10000; i++ {
			codes = append(codes, strconv.Itoa(i))
		}

		ticker := time.NewTicker(10 * time.Second)
		for {

			select {
			case <-ticker.C:
				for _, sev := range services {
					for _, md := range method {
						for _, rc := range rpcCode {
							for _, code := range codes {
								gometric.GetGlobalMeter().NewTimer("dacs_rpc_request_duration_seconds").
									AddTag("dc_sname", "config_server").
									AddTag("dc_sver", "V2.2.10025").
									AddTag("dc_ip", "10.10.24.217").
									AddTag("service", sev).
									AddTag("method", md).
									AddTag("ret", strconv.Itoa(rand.Intn(2))).
									// AddTag("rpc_code", rc).
									// AddTag("code", code).
									Update(time.Duration(rand.Intn(15)) * time.Second)

								gometric.GetGlobalMeter().NewCounter("dacs_rpc_request_total").
									AddTag("dc_sname", "config_server").
									AddTag("dc_sver", "V2.2.10025").
									AddTag("dc_ip", "10.10.24.217").
									AddTag("service", sev).
									AddTag("method", md).
									AddTag("rpc_code", rc).
									AddTag("code", code).
									IncrOnce()
							}
						}
					}
				}
			}
		}

	}()

	_ = r.Run(":8080")
}
