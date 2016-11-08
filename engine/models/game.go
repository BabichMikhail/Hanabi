package enginemodels

import "time"

const (
	GameWait     = 1
	GameActive   = 2
	GameInActive = 4
	GameUnknown  = 8
)

func GameStatusName(status int) string {
	switch status {
	case GameWait:
		return "wait"
	case GameActive:
		return "active"
	case GameInActive:
		return "inactive"
	default:
		return "unknown"
	}
}

type Game struct {
	Id      int
	Owner   string
	Status  string
	Players []Player
}

type GameItem struct {
	Id           int       `orm:"column(id)"`
	OwnerId      int       `orm:"column(owner_id)"`
	Owner        string    `orm:"column(owner)`
	StatusCode   int       `orm:"column(status)"`
	Status       string    ``
	PlayerCount  int       `orm:"column(count)"`
	PlayerPlaces int       `orm:"column(places)"`
	Players      []Player  ``
	UserIn       bool      ``
	URL          string    ``
	Created      time.Time `orm:"column(created)"`
}
