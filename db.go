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
	webcfg *ini.Section

	dbname   = "default" // database alias
	dbinited = false
	dbtype   = "mysql"
	dburl    = ""
	maxidle  = 2
	maxconn  = 2
)

// initialize the db driver and config
func InitDB() {
	if dbinited {
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

	dbtype = web.Key("dbtype").String()
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
		maxidle = dbcfg.Key("maxidle").MustInt(2)
		maxconn = dbcfg.Key("maxconn").MustInt(2)
		dburl = username + ":" + password + "@" + url
		orm.RegisterDriver(dbtype, orm.DR_MySQL)
	case "sqlite":
		dburl = dbcfg.Key("url").String()
		dbtype = "sqlite3"
		orm.RegisterDriver(dbtype, orm.DR_Sqlite)
	}

	dbinited = true
}

// Connect to db and try to sync table structs
func SyncDB() {

	// try to connnet to db
	if err := orm.RegisterDataBase(dbname, dbtype, dburl, maxidle, maxconn); err != nil {
		panic(err)
	}

	// create the tables
	force := false                         // drop table then create it
	sqllog := webcfg.Key("sqlon").String() // log the sql string
	verbose := false
	if "on" == sqllog {
		verbose = true
	}

	// try to sync table structs
	err := orm.RunSyncdb(dbname, force, verbose)
	if err != nil {
		log.Error(err.Error())
	}
}
