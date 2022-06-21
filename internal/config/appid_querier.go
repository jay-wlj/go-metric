package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
http://lalaplat2.dtl.work/index.php?_g=api&_m=domain&_a=getAppid&domain=router.my.cn
{
  "ret": 0,
  "msg": "ok",
  "data": "ai-router-api"
}
*/

var (
	once4AppIdQuerier     sync.Once
	appIdQuerierSingleton *appIdQuerier
)

func AppIdQuerier(cfg *Config) *appIdQuerier {
	once4AppIdQuerier.Do(func() {
		appIdQuerierSingleton = &appIdQuerier{
			cfg: cfg,
			client: &http.Client{
				Timeout: time.Millisecond * 100,
			},
		}
		// 每天清除一次
		go appIdQuerierSingleton.cleaner()
	})
	return appIdQuerierSingleton
}

type appIdQuerier struct {
	cfg    *Config
	client *http.Client
	appids sync.Map // host => appid
}

func (q *appIdQuerier) GetAppId(domain string) string {
	domain = strings.TrimSpace(domain)
	if domain == "" || domain == "UNKNOWN" {
		return "-"
	}
	appidI, ok := q.appids.Load(domain)
	if ok {
		return appidI.(string)
	}
	appId := q.query(domain)
	q.appids.Store(domain, appId)
	q.cfg.WriteInfoOrNot(fmt.Sprintf("get appid: %s by host: %s", appId, domain))
	return appId
}

func (q *appIdQuerier) query(host string) string {
	const url = "http://lalaplat2.dtl.work/index.php?_g=api&_m=domain&_a=getAppid&domain="
	req, _ := http.NewRequest(http.MethodGet, url+host, nil)
	resp, err := q.client.Do(req)
	if err != nil {
		q.cfg.WriteErrorOrNot(fmt.Sprintf("failed to query appid of host: %s, %s", host, err.Error()))
		return host
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		q.cfg.WriteErrorOrNot(fmt.Sprintf("failed to query appid of host: %s, %s", host, err.Error()))
	}
	type Resp struct {
		Ret  int    `json:"ret"`
		Data string `json:"data"`
	}
	var r Resp
	if err := json.Unmarshal(content, &r); err != nil {
		q.cfg.WriteErrorOrNot(fmt.Sprintf("failed to query appid of host: %s, %s", host, err.Error()))
		return host
	}
	appid := strings.TrimSpace(r.Data)
	if r.Ret == 0 && len(appid) > 0 {
		return appid
	}
	return host
}

func (q *appIdQuerier) cleaner() {
	q.cfg.WriteInfoOrNot("appIdQuerier is running, cached host->appid mapping will clear in every 24h")
	ticker := time.NewTicker(time.Hour * 24)
	for range ticker.C {
		q.appids.Range(func(key, _ interface{}) bool {
			q.appids.Delete(key)
			return true
		})
	}
}
