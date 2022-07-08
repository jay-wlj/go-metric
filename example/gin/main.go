package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
		gometric.WithPush("10.10.21.37:7072", 15*time.Second),
		gometric.WithPrometheusPort(0),
	)
	r := gin.New()
	r.Use(otelgin.HTTPServerTimerMiddleware())

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"ret":  0,
			"name": "unknown",
			"id":   id,
		})
	})
	_ = r.Run(":8080")
}
