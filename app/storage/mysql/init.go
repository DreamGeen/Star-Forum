package mysql

import (
	"fmt"
	"star/app/constant/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Client *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.Conf.MysqlUser,
		settings.Conf.MysqlPassword,
		settings.Conf.MysqlHost,
		settings.Conf.MysqlPort,
		settings.Conf.MysqlDbname,
	)
	var err error
	Client, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
}
