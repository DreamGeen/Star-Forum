package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"star/settings"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.Conf.MysqlUser,
		settings.Conf.MysqlPassword,
		settings.Conf.MysqlHost,
		settings.Conf.MysqlPort,
		settings.Conf.MysqlDbname,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
		return
	}
	return nil
}
func Close() {
	_ = db.Close()
}
