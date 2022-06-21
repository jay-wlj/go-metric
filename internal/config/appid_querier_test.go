package config

import "testing"

func Test_Querier_GetAppId(t *testing.T) {
	q := AppIdQuerier(&Config{})
	q.GetAppId("router.my.cn")
	q.GetAppId("rcsapi-stg.my.cn")
	q.GetAppId("geosvr-stg.my.cn")
	q.GetAppId("rcsapi-stg.my.cn")
	q.GetAppId("dsinfo2-stg.my.cn")
}
