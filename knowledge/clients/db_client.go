package clients

import (
	"doc-precise-rag/knowledge/config"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func init() {
	mysqlConfig := config.AppConfig.MysqlConfig
	loc, _ := time.LoadLocation("Asia/Shanghai")
	cfg := mysql.Config{
		User:                 *mysqlConfig.User,
		Passwd:               *mysqlConfig.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", *mysqlConfig.Host, *mysqlConfig.Port),
		DBName:               *mysqlConfig.DBName,
		AllowNativePasswords: true,
		Loc:                  loc,
		ParseTime:            true,
		MultiStatements:      false,
	}

	db, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(20)                  // 最大打开连接
	db.SetMaxIdleConns(10)                  // 最大空闲连接
	db.SetConnMaxLifetime(30 * time.Minute) // 连接最长存活
	db.SetConnMaxIdleTime(10 * time.Minute) // 空闲多久回收

	DB = db
}
