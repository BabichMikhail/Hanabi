package models

import (
	wetalk "github.com/beego/wetalk/modules/models"
)

type User struct {
	wetalk.User
}

type Player struct {
	Id     int `orm:"auto"`
	UserId int `orm:"default(0);column(user_id)"`
	GameId int `orm:"default(0);column(game_id)"`
}

func (p *Player) TableName() string {
	return "players"
}
