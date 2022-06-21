package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	gometric "github.com/jay-wlj/go-metric"
)

func main() {
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
	)
	startTime := time.Now()

	req, _ := http.NewRequest(http.MethodPost, "http://lalaplat2-dev.dtl.work/index.php", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	type DtlResp struct {
		Ret  int         `json:"ret"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	var dtResp DtlResp
	if err = json.Unmarshal(data, &dtResp); err != nil {
		panic(err)
	}
	gometric.GetGlobalMeter().
		Components().
		NewHTTPClientTimerFromResponse(resp, strconv.Itoa(dtResp.Ret)).
		UpdateSince(startTime)
	time.Sleep(time.Second * 30)
}
