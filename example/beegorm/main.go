package main

import (
	middleware "go-metric/example/beegorm/middleware"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	gometric "github.com/jay-wlj/go-metric"
)

type User struct {
	Id   int
	Name string
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(User))
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
}

func main() {
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
	)
	orm.Debug = true
	orm.NewOrm()
	orm.RunSyncdb("default", false, true) //执行操作

	//o := wraporm(orm.NewOrm())
	o := middleware.Wraporm(orm.NewOrm(), "-", "ci-mysql-resource")
	user := User{Name: "user1"}
	o.Insert(&user)
	o.Delete(&user)
	o.Update(&User{Id: 1, Name: "user2"})
	select {}
}
