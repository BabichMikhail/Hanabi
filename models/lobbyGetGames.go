package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type LobbyGame struct {
	Id          int           `json:"id" orm:"column(id)"`
	PlayerCount int           `json:"player_count" orm:"column(player_count)"`
	OwnerId     int           `json:"owner_id" orm:"column(owner_id)"`
	OwnerName   string        `json:"owner_name" orm:"column(owner_name)"`
	Status      int           `json:"status" orm:"column(status)"`
	StatusName  string        `json:"status_name"`
	CreatedAt   time.Time     `json:"created_at" orm:"column(created_at)"`
	UserIn      bool          `json:"user_in"`
	Players     []LobbyPlayer `json:"players"`
	URL         string        `json:"URL"`
}

func getGames(gameStatuses []int) (games []LobbyGame) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("g.id, g.status, g.player_count, g.created_at, g.owner_id, u.nick_name as owner_name").
		From("games g").
		InnerJoin("user u").On("u.id = g.owner_id").
		Where("status").In(IntSliceToString(gameStatuses)).
		OrderBy("created_at").Desc()
	sql := qb.String()
	o.Raw(sql).QueryRows(&games)

	gameIds := []int{}
	for _, game := range games {
		gameIds = append(gameIds, game.Id)
	}

	players := GetGamePlayers(gameIds)
	for i, game := range games {
		games[i].Players = players[game.Id]
		games[i].StatusName = StatusName(game.Status)
	}

	return
}

func GetFinishedGames() (games []LobbyGame) {
	return getGames([]int{StatusInactive})
}

func GetMyGames(userId int) (games []LobbyGame) {
	gamesAll := getGames(GetAllStatuses())

	for i, _ := range gamesAll {
		for j, _ := range gamesAll[i].Players {
			if gamesAll[i].Players[j].Id == userId {
				gamesAll[i].UserIn = true
				games = append(games, gamesAll[i])
				break
			}
		}
	}
	return
}

func GetAllGames() (games []LobbyGame) {
	return getGames([]int{StatusActive, StatusWait})
}

func GetActiveGames() (games []LobbyGame) {
	return getGames([]int{StatusActive, StatusWait})
}
