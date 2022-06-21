package main

import (
	"encoding/json"
	"go-metric/example/beego/middleware"
	"net/http"

	"github.com/astaxie/beego"
	gometric "github.com/jay-wlj/go-metric"
)

type ExampleController struct {
	beego.Controller
}

func (c *ExampleController) Get() {
	type Response struct {
		Ret  int         `json:"ret"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	var resp = Response{
		Ret: 1,
		Msg: "ok",
	}
	data, _ := json.Marshal(resp)
	c.Ctx.WriteString(string(data))
}

func (c *ExampleController) Post() {
	c.CustomAbort(http.StatusNoContent, "not-allowed")
}

func main() {
	//  Init the trace and meter provider
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
	)

	// Disable autorender
	beego.BConfig.WebConfig.AutoRender = true

	// Create routes
	beego.Router("/users/:id", &ExampleController{}, "get:Get")
	beego.Router("/users", &ExampleController{}, "post:Post")

	// Create the middleware
	mware := middleware.NewOTelBeegoMiddleWare()

	// Start the server using the OTel middleware
	beego.RunWithMiddleWares(":8000", mware)
}
