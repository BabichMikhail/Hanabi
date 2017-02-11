package models

import (
	"github.com/astaxie/beego/orm"
	wetalk "github.com/beego/wetalk/modules/models"
)

type User struct {
	wetalk.User
}

func GetUserNickNameById(id int) string {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("nick_name").
		From("user").
		Where("id = ?")
	sql := qb.String()
	var nickName string
	if err := o.Raw(sql, id).QueryRow(&nickName); err != nil {
		return ""
	} else {
		return nickName
	}
}

type Player struct {
	Id     int `orm:"auto"`
	UserId int `orm:"default(0);column(user_id)"`
	GameId int `orm:"default(0);column(game_id)"`
}

func (p *Player) TableName() string {
	return "players"
}

func (p *Player) TableUnique() [][]string {
	return [][]string{
		[]string{"user_id", "game_id"},
	}
}

func (p *Player) TableIndex() [][]string {
	return [][]string{
		[]string{"game_id"},
		[]string{"user_id"},
	}
}
