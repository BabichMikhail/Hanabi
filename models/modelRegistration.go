package models

import (
	"github.com/astaxie/beego/orm"
)

func Registration() {
	orm.RegisterModel(new(Game))
	orm.RegisterModel(new(Player))
}
