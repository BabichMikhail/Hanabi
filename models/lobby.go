package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type LobbyPlayer struct {
	Id       int    `orm:"column(id)" json:"id"`
	NickName string `orm:"column(nick_name)" json:"nick_name"`
}

type LobbyGameItem struct {
	Id          int           `orm:"column(id)"`
	OwnerId     int           `orm:"column(owner_id)"`
	Owner       string        `orm:"column(owner)"`
	StatusCode  int           `orm:"column(status)"`
	Status      string        ``
	PlayerCount int           `orm:"column(count)"`
	Players     []LobbyPlayer ``
	UserIn      bool          ``
	URL         string        ``
	Created     time.Time     `orm:"column(created)"`
}

const (
	StatusWait = 1 << iota
	StatusActive
	StatusFinished
	StatusUnknown
)

func StatusName(status int) string {
	switch status {
	case StatusWait:
		return "wait"
	case StatusActive:
		return "active"
	case StatusFinished:
		return "finished"
	default:
		return "unknown"
	}
}

func GetAllStatuses() []int {
	return []int{StatusWait, StatusActive, StatusFinished}
}

type PlayerCount struct {
	GameId int `orm:"column(game_id)"`
	Count  int `orm:"column(count)"`
}

func getPlayerCount(ids []int) []PlayerCount {
	o := orm.NewOrm()
	if len(ids) == 0 {
		return []PlayerCount{}
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("game_id", "count(user_id) as count").
		From("players").
		Where("game_id").In(IntSliceToString(ids)).
		GroupBy("game_id")
	sql := qb.String()
	var playersCount []PlayerCount
	o.Raw(sql).QueryRows(&playersCount)
	return playersCount
}

type UserGame struct {
	Id int `orm:"column(game_id)"`
}

func getUserGames(userId int) (ids []int) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("game_id").
		From("players").
		Where("user_id = ?")
	sql := qb.String()
	var userGames []struct {
		Id int `orm:"column(game_id)"`
	}
	o.Raw(sql, userId).QueryRows(&userGames)
	for _, userGame := range userGames {
		ids = append(ids, userGame.Id)
	}
	return ids
}

func GetGameList(status []int, userId int) (games []LobbyGameItem) {
	o := orm.NewOrm()

	args := IntSliceToString(status)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select(
		"g.id", "g.player_count as count", "g.owner_id as owner_id",
		"u.nick_name as owner", "g.status", "g.created_at").
		From("games g").
		InnerJoin("user u").
		On("u.id = g.owner_id").
		Where("status").In(args).
		OrderBy("created_at").Desc()
	sql := qb.String()
	o.Raw(sql).QueryRows(&games)
	for i, game := range games {
		games[i].Status = StatusName(game.StatusCode)
	}

	ids := []int{}
	for _, gameItem := range games {
		ids = append(ids, gameItem.Id)
	}
	playersCount := getPlayerCount(ids)
	gamesMap := map[int]int{}
	for _, v := range playersCount {
		gamesMap[v.GameId] = v.Count
	}
	playersMap := GetGamePlayers(ids)
	userGames := getUserGames(userId)
	userInGame := map[int]bool{}
	for _, gameId := range userGames {
		userInGame[gameId] = true
	}
	for i, g := range games {
		games[i].UserIn = userInGame[g.Id] == true
		games[i].Players = playersMap[g.Id]
	}

	return
}
