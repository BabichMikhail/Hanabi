package models

import (
	"time"
)

const (
	GameWait      = 1
	GameActive    = 2
	GameNonActive = 4
)

type Game struct {
	Id      int       `orm:"auto"`
	Status  int       `orm:"default(4)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

func (g *Game) TableName() string {
	return "games"
}

type Player struct {
	Id     int `orm:"auto"`
	UserId int `orm:"default(0)"`
	GameId int `orm:"default(0)"`
}

func (p *Player) TableName() string {
	return "players"
}
