package utilities

import (
	"database/sql"
	"log"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MySQL *sql.DB

func InitMySQL() {
	var err error
	MySQL, err = sql.Open("mysql", "ovaphlow:ovaph@QH.1123@tcp(82.156.226.151:3306)/crate?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err.Error())
	}
	if err = MySQL.Ping(); err != nil {
		log.Println("连接数据库失败")
		log.Println(err.Error())
	}
	MySQL.SetConnMaxLifetime(time.Second * 30)
	MySQL.SetMaxIdleConns(runtime.NumCPU()*2 + 1)
	MySQL.SetMaxOpenConns(runtime.NumCPU()*2 + 1)
}
