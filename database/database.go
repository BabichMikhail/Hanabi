package database

import (
	"github.com/BabichMikhail/Hanabi/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:root@/hanabi?charset=utf8")
	if err != nil {
		panic(err)
	}
	models.Registration()
	forse := false
	verbose := true
	orm.RunSyncdb("default", forse, verbose)
}
