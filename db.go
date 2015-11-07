package tkits

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbname      = "default" // 数据库别名
	webcfg      *ini.Section
	dbconnected = false
)

func ConnectDB() {
	if dbconnected {
		return
	}
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
	orm.Debug = true

	web, err := cfg.GetSection("web")
	if err != nil {
		panic(err)
	}
	webcfg = web

	dbtype := web.Key("dbtype").String()
	log.Debugf("DB type is %s", dbtype)
	dbcfg, err := cfg.GetSection(dbtype)
	if err != nil {
		panic(err)
	}

	switch dbtype {
	case "mysql":
		var username string = dbcfg.Key("username").String()
		if username, err = GetCrypto().DecryptStr(username); err != nil {
			panic(err)
		}

		var password string = dbcfg.Key("password").String()
		if password, err = GetCrypto().DecryptStr(password); err != nil {
			panic(err)
		}

		url := dbcfg.Key("url").String()
		maxidle := dbcfg.Key("maxidle").MustInt(2)
		maxconn := dbcfg.Key("maxconn").MustInt(2)
		orm.RegisterDriver("mysql", orm.DR_MySQL)
		orm.RegisterDataBase(dbname, "mysql",
			username+":"+password+"@"+url,
			maxidle, maxconn)
	case "sqlite":
		url := dbcfg.Key("url").String()
		orm.RegisterDriver("sqlite3", orm.DR_Sqlite)
		orm.RegisterDataBase(dbname, "sqlite3", url)
	}

	dbconnected = true
}

func SyncDB() {
	force := false                         // drop table 后再建表
	sqllog := webcfg.Key("sqlon").String() // 打印执行过程
	verbose := false
	if "on" == sqllog {
		verbose = true
	}

	// 遇到错误立即返回
	err := orm.RunSyncdb(dbname, force, verbose)
	if err != nil {
		log.Error(err.Error())
	}
}
