package enginelobby

import (
	"errors"
	"fmt"
	"time"

	wetalk "github.com/beego/wetalk/modules/models"
)

const (
	GameWait = 1 << iota
	GameActive
	GameInactive
	GameUnknown
)

func GameStatusName(status int) string {
	switch status {
	case GameWait:
		return "wait"
	case GameActive:
		return "active"
	case GameInactive:
		return "inactive"
	default:
		return "unknown"
	}
}

type Game struct {
	Id      int      `json:"id"`
	Owner   string   `json:"owner"`
	Status  string   `json:"status"`
	Players []Player `json:"players"`
}

type GameItem struct {
	Id           int       `orm:"column(id)"`
	OwnerId      int       `orm:"column(owner_id)"`
	Owner        string    `orm:"column(owner)"`
	StatusCode   int       `orm:"column(status)"`
	Status       string    ``
	PlayerCount  int       `orm:"column(count)"`
	PlayerPlaces int       `orm:"column(places)"`
	Players      []Player  ``
	UserIn       bool      ``
	URL          string    ``
	Created      time.Time `orm:"column(created)"`
}

func RevertGameItems(items []GameItem) []GameItem {
	for i := 0; i < len(items)/2; i++ {
		items[i], items[len(items)-i-1] = items[len(items)-i-1], items[i]
	}
	return items
}

func CopyGameItems(items []GameItem) []GameItem {
	copyItems := make([]GameItem, len(items))
	copy(copyItems, items)
	return copyItems
}

type Player struct {
	Id       int    `orm:"column(id)" json:"id"`
	NickName string `orm:"column(nick_name)" json:"nick_name"`
}

func GetAllStatuses() []int {
	return []int{GameWait, GameActive, GameInactive}
}

func MakeGame(id int, user wetalk.User) (game Game, err error) {
	if id <= 0 {
		err = errors.New(fmt.Sprintf("Can't make game with Id = %d", id))
		return
	}
	game.Id = id
	game.Status = GameStatusName(GameWait)
	game.Owner = user.NickName
	player := Player{user.Id, user.NickName}
	game.Players = []Player{player}
	return game, nil
}
