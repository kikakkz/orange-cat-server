package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"kkt.com/glog"
)

var db *sql.DB

func DBOpen(cfg *MysqlCfg) error {
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s", cfg.User, cfg.Password, "tcp", cfg.Host, cfg.Db)
	var err error
	db, err = sql.Open("mysql", dsn)
	if nil != err {
		glog.Error(err, dsn)
		return err
	}
	return nil
}

func DBQuery(sqlExec string, scanner func(rows *sql.Rows) error) error {
	glog.Info(sqlExec, "------ Start")

	rows, err := db.Query(sqlExec)
	if nil != err {
		glog.Error(err)
		return err
	}

	glog.Info(sqlExec, "------ End")

	defer rows.Close()
	for rows.Next() {
		err := scanner(rows)
		if nil != err {
			glog.Warning(err)
			continue
		}
	}
	return nil
}

func DBClose() {
	db.Close()
}

func DBExec(sqlExec string) error {
	_, err := db.Exec(sqlExec)
	return err
}
