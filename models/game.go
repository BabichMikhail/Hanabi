package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	engineModels "github.com/BabichMikhail/Hanabi/engine/models"
	"github.com/astaxie/beego/orm"
)

type Game struct {
	Id           int       `orm:"auto"`
	OwnerId      int       `orm:"column(owner_id)"`
	PlayersCount int       `orm:"column(players_count)"`
	Status       int       `orm:"column(status);default(4)"`
	Created      time.Time `orm:"column(created_at);auto_now_add;type(timestamp)"`
}

func (g *Game) TableName() string {
	return "games"
}

func NewGame(userId int, playersCount int, status int) (id int) {
	if playersCount > 5 {
		playersCount = 5
	}
	if playersCount < 2 {
		playersCount = 2
	}
	o := orm.NewOrm()
	o.Begin()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("games", "owner_id", "players_count", "status", "created_at").
		Values("?", "?", "?", "CURRENT_TIMESTAMP")
	sql := qb.String()
	if res, err := o.Raw(sql, userId, playersCount, status).Exec(); err == nil {
		id64, _ := res.LastInsertId()
		id = int(id64)
	}
	if id > 0 {
		if err, _ := JoinGame(id, userId); err == nil {
			o.Commit()
		} else {
			id = 0
			o.Rollback()
		}
	} else {
		id = 0
		o.Rollback()
	}
	return
}

func ActivateGame(gameId int) int {
	o := orm.NewOrm()
	game := new(Game)
	game.Id = gameId
	o.Read(game)
	game.Status = engineModels.GameActive
	o.Update(game)
	return game.Status
}

func CheckGame(gameId int) (status int) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("players_count, status").
		From("games").
		Where("id = ?")
	sql := qb.String()
	var g Game
	o.Raw(sql, gameId).QueryRow(&g)
	currentPlayers := len(GetGamePlayers([]int{gameId})[gameId])
	if g.PlayersCount == currentPlayers {
		status = ActivateGame(gameId)
	} else {
		status = g.Status
	}
	return status
}

func JoinGame(gameId int, userId int) (err error, status string) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("players", "user_id", "game_id").
		Values("?", "?")
	sql := qb.String()
	if _, err = o.Raw(sql, userId, gameId).Exec(); err == nil {
		o.Commit()
		status = engineModels.GameStatusName(CheckGame(gameId))
	} else {
		status = ""
		o.Rollback()
	}
	return err, status
}

// @todo move to utils
func IntSliceToString(slice []int) string {
	values := []string{}
	for _, i := range slice {
		values = append(values, strconv.Itoa(i))
	}
	return strings.Join(values, ", ")
}

func LeaveGame(gameId int, userId int) (string, error) {
	o := orm.NewOrm()

	exist := o.QueryTable("players").
		Filter("game_id", gameId).
		Filter("user_id", userId).
		Exist()
	if !exist {
		return "", errors.New(fmt.Sprintf("User %d not fount in game %d", userId, gameId))
	}

	count, err := o.QueryTable("players").
		Filter("game_id", gameId).
		Count()
	if err != nil {
		return "", err
	}

	o.Begin()
	_, err = o.QueryTable("players").
		Filter("user_id", userId).
		Filter("game_id", gameId).
		Delete()
	if err != nil {
		o.Rollback()
		return "", err
	}
	action := "leave"
	if count == 1 {
		action = "delete"
		num, err := o.QueryTable("games").
			Filter("id", gameId).
			Delete()
		if err != nil {
			o.Rollback()
			return "", err
		}
		if num != 1 {
			o.Rollback()
			return "", errors.New(fmt.Sprintf("Security problem with game %d", gameId))
		}
	}

	o.Commit()
	return action, err
}

func GetGamePlayers(gameIds []int) map[int]([]engineModels.Player) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	inString := IntSliceToString(gameIds)
	qb.Select("p.user_id", "u.nick_name", "p.game_id").
		From("players p").
		InnerJoin("user u").
		On("u.id = p.user_id").
		Where("p.game_id").In(inString)
	sql := qb.String()
	var splayers []struct {
		UserId   int    `orm:"column(user_id)"`
		NickName string `orm:"column(nick_name)"`
		GameId   int    `orm:"column(game_id)"`
	}
	o.Raw(sql).QueryRows(&splayers)
	playersMap := map[int]([]engineModels.Player){}
	for _, v := range splayers {
		playersMap[v.GameId] = append(playersMap[v.GameId], engineModels.Player{v.UserId, v.NickName})
	}
	return playersMap
}

func GetGameList(status []int, userId int) (games []engineModels.GameItem) {
	o := orm.NewOrm()

	args := IntSliceToString(status)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select(
		"g.id", "g.players_count as places", "g.owner_id as owner_id",
		"u.nick_name as owner", "g.status", "g.created_at").
		From("games g").
		InnerJoin("user u").
		On("u.id = g.owner_id").
		Where("status").In(args).
		OrderBy("created_at").Desc()
	sql := qb.String()
	o.Raw(sql).QueryRows(&games)
	for i, game := range games {
		games[i].Status = engineModels.GameStatusName(game.StatusCode)
	}

	ids := []int{}
	for _, gameItem := range games {
		ids = append(ids, gameItem.Id)
	}
	playersCount := getPlayersCount(ids)
	gamesMap := map[int]int{}
	for _, v := range playersCount {
		gamesMap[v.Id] = v.Count
	}
	playersMap := GetGamePlayers(ids)
	userGames := getUserGames(userId)
	userInMap := map[int]bool{}
	for _, v := range userGames {
		userInMap[v.Id] = true
	}
	for i, g := range games {
		games[i].UserIn = userInMap[g.Id] == true
		games[i].Players = playersMap[g.Id]
		games[i].PlayerCount = gamesMap[g.Id]
	}

	return
}

type UserGame struct {
	Id int `orm:"column(game_id)"`
}

type PlayerCount struct {
	Id    int `orm:"column(game_id)"`
	Count int `orm:"column(count)"`
}

func getPlayersCount(ids []int) []PlayerCount {
	o := orm.NewOrm()
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

func getUserGames(userId int) []UserGame {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("game_id").
		From("players").
		Where("user_id = ?")
	sql := qb.String()
	var userGames []UserGame
	o.Raw(sql, userId).QueryRows(&userGames)
	return userGames
}

type GameStatus struct {
	GameId     int    `orm:"column(id)" json:"game_id"`
	StatusCode int    `orm:"column(status)" json:"status_code"`
	StatusName string `json:"status_name"`
	URL        string `json:"URL"`
}

func GetStatuses(userId int) []GameStatus {
	o := orm.NewOrm()
	userGames := func(ug []UserGame) string {
		ans := []string{}
		for _, g := range ug {
			ans = append(ans, strconv.Itoa(g.Id))
		}
		return strings.Join(ans, ", ")
	}(getUserGames(userId))

	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id, status").
		From("games").
		Where("id").In(userGames)
	sql := qb.String()
	var statuses []GameStatus
	o.Raw(sql).QueryRows(&statuses)
	for i, s := range statuses {
		statuses[i].StatusName = engineModels.GameStatusName(s.StatusCode)
	}
	return statuses
}
