package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ngaut/log"
)

var Mysql *sql.DB

var (
	host     = "localhost"
	database = "db_muser"
	user     = "root"
	passwd   = ""
	port     = ""
)

func init() {
	// host = Conf.String("mysql::host")
	// user = Conf.String("mysql::user")
	// passwd = Conf.String("mysql::pass")
	// database = Conf.String("mysql::db")
	// port = Conf.String("mysql::port")
	// ODB = &SConnDB{conn()}
	Mysql = conn()
}

// 连接数据库
func conn() (db *sql.DB) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", "root", "123456", "127.0.0.1", "33060", "guess_plane"))
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(100)
	err = db.Ping()
	if err == nil {
		log.Infof("Mysql Success:Host=\"%s\";DB=\"%s\";", host, database)
	} else {
		log.Error(err, 111)
		os.Exit(0)
	}
	return
}
