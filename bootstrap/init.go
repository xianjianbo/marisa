package bootstrap

import (
	"github.com/xianjianbo/marisa/library/resource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() {
	if err := initMysql(); err != nil {
		panic(err)
	}
}

func initMysql() (err error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root123456@tcp(127.0.0.1:3306)/marisa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	resource.MysqlClientGorm = db
	return
}
