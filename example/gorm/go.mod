module github.com/jay-wlj/go-metric/example/gorm

replace github.com/jay-wlj/go-metric => ../../

go 1.15

require (
	github.com/jay-wlj/go-metric v1.0.3
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)
