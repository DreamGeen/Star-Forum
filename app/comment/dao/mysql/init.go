package mysql

import (
	"fmt"
	"star/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.Conf.MysqlConfig.MysqlUser,
		settings.Conf.MysqlConfig.MysqlPassword,
		settings.Conf.MysqlConfig.MysqlHost,
		settings.Conf.MysqlConfig.MysqlPort,
		settings.Conf.MysqlConfig.MysqlDbname)
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return nil
}

func Close() {
	db.Close()
}
