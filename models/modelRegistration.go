package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func RegisterDatabase(driverName string, dataSource string, maxIdle int, maxOpen int) {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		beego.Error(err)
	}
	registerModels()
	forse := false
	verbose := true
	orm.RunSyncdb("default", forse, verbose)
}

func registerModels() {
	orm.RegisterModel(new(Game))
	orm.RegisterModel(new(Player))
}
