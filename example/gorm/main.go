package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	gometric "github.com/jay-wlj/go-metric"
	"github.com/jay-wlj/go-metric/instrumentation/github.com/gorm.io/gorm/otelgorm"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Age       uint8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	_ = gometric.NewMeter(
		gometric.WithMeterProvider(gometric.PrometheusMeterProvider),
	)
	cfg := &gorm.Config{
		Logger: otelgorm.NewTraceRecorder(logger.Default, "ci-mysql-resource"), // use "" when resource-id doesn't exist
	}

	gormDB, err := gorm.Open(sqlite.Open("gorm.db"), cfg)
	if err != nil {
		panic("failed to open database")
	}

	if err := gormDB.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	gormDB.Create(&User{
		Name:  "san.zhang",
		Email: "san.zhang@dtl.cn",
		Age:   3,
	})
	gormDB.Create(&User{
		Name:  "si.li",
		Email: "si.li@dtl.cn",
		Age:   4,
	})

	var user User
	gormDB.First(&user)

	time.Sleep(time.Minute)
}
