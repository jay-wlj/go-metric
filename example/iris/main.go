package main

import (
	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/example/iris/middleware"
	"github.com/kataras/iris/v12"
)

func main() {
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
	)
	app := iris.New()

	app.Use(middleware.UseMiddleware)

	app.Get("/users/{id:int64}", func(ctx iris.Context) {
		userID, _ := ctx.Params().GetUint64("id")
		ctx.JSON(iris.Map{"ret": 1,
			"id":   userID,
			"name": "unknown"})
	})

	app.Listen(":8080")
}
