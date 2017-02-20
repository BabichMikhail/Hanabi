package models

import (
	"errors"
	"fmt"
	"time"

	lobby "github.com/BabichMikhail/Hanabi/engine/lobby"
	"github.com/astaxie/beego/orm"
)

type Game struct {
	Id           int        `orm:"auto" json:"id"`
	OwnerId      int        `orm:"column(owner_id)" json:"owner_id"`
	PlayerCount  int        `orm:"column(player_count)" json:"player_count"`
	Status       int        `orm:"column(status);default(4)" json:"status"`
	Points       int        `orm:"column(points);default(0)"`
	Seed         int64      `orm:"column(seed);null"`
	InitState    *GameState `orm:"-"`
	CurrentState *GameState `orm:"-"`
	Actions      []*Action  `orm:"-"`
	CreatedAt    time.Time  `orm:"column(created_at);auto_now_add;type(timestamp)" json:"created_at"`
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
	qb.InsertInto("games", "owner_id", "player_count", "status", "created_at").
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
	game.Status = lobby.GameActive
	o.Update(game)
	return game.Status
}

func CheckGame(gameId int) (status int) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("player_count, status").
		From("games").
		Where("id = ?")
	sql := qb.String()
	var g Game
	o.Raw(sql, gameId).QueryRow(&g)
	currentPlayers := len(GetGamePlayers([]int{gameId})[gameId])
	if g.PlayerCount == currentPlayers {
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
		gamePlayers := GetGamePlayers([]int{gameId})[gameId]
		qbGames, _ := orm.NewQueryBuilder("mysql")
		sql := qbGames.Select("player_count").From("games").Where("id = ?").String()
		var game Game
		o.Raw(sql, gameId).QueryRow(&game)
		if game.PlayerCount == len(gamePlayers) {
			playerIds := []int{}
			for i := 0; i < len(gamePlayers); i++ {
				playerIds = append(playerIds, gamePlayers[i].Id)
			}
			_, err = CreateActiveGame(playerIds, gameId)
		}
	}

	if err == nil {
		o.Commit()
		status = lobby.GameStatusName(CheckGame(gameId))
	} else {
		status = ""
		o.Rollback()
	}
	return err, status
}

func LeaveGame(gameId int, userId int) (string, error) {
	o := orm.NewOrm()

	var players []Player
	qbPlayers, _ := orm.NewQueryBuilder("mysql")
	qbPlayers.Select("user_id").From("players").Where("game_id = ?")
	_, err := o.Raw(qbPlayers.String(), gameId).QueryRows(&players)
	if err != nil {
		return "", err
	}

	var game Game
	qbGame, _ := orm.NewQueryBuilder("mysql")
	qbGame.Select("status, player_count").From("games").Where("id = ?")
	err = o.Raw(qbGame.String(), gameId).QueryRow(&game)
	if err != nil {
		return "", err
	}

	if game.Status != int(lobby.GameWait) {
		return "", errors.New(fmt.Sprintf("User #%d can't leave game #%d", userId, gameId))
	}

	exist := false
	for _, player := range players {
		exist = exist || player.UserId == userId
	}

	if !exist {
		return "", errors.New(fmt.Sprintf("User #%d not found in game #%d", userId, gameId))
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
	if len(players) == 1 {
		action = "delete"
		num, errDelete := o.QueryTable("games").
			Filter("id", gameId).
			Delete()
		if errDelete != nil {
			o.Rollback()
			return "", errDelete
		}

		if num != 1 {
			o.Rollback()
			return "", errors.New(fmt.Sprintf("Security problem with game %d", gameId))
		}
	}

	o.Commit()
	return action, err
}

func GetGamePlayers(gameIds []int) map[int]([]lobby.Player) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	playersMap := map[int]([]lobby.Player){}
	if len(gameIds) > 0 {
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

		for _, v := range splayers {
			playersMap[v.GameId] = append(playersMap[v.GameId], lobby.Player{
				Id:       v.UserId,
				NickName: v.NickName,
			})
		}
	}
	return playersMap
}

func GetGameList(status []int, userId int) (games []lobby.GameItem) {
	o := orm.NewOrm()

	args := IntSliceToString(status)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select(
		"g.id", "g.player_count as places", "g.owner_id as owner_id",
		"u.nick_name as owner", "g.status", "g.created_at").
		From("games g").
		InnerJoin("user u").
		On("u.id = g.owner_id").
		Where("status").In(args).
		OrderBy("created_at").Desc()
	sql := qb.String()
	o.Raw(sql).QueryRows(&games)
	for i, game := range games {
		games[i].Status = lobby.GameStatusName(game.StatusCode)
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
	userInGame := map[int]bool{}
	for _, v := range userGames {
		userInGame[v.Id] = true
	}
	for i, g := range games {
		games[i].UserIn = userInGame[g.Id] == true
		games[i].Players = playersMap[g.Id]
		games[i].PlayerCount = gamesMap[g.Id]
	}

	return
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

type UserGame struct {
	Id int `orm:"column(game_id)"`
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

type LobbyGame struct {
	Id          int            `json:"id" orm:"column(id)"`
	PlayerCount int            `json:"player_count" orm:"column(player_count)"`
	OwnerId     int            `json:"owner_id" orm:"column(owner_id)"`
	OwnerName   string         `json:"owner_name" orm:"column(owner_name)"`
	Status      int            `json:"status" orm:"column(status)"`
	StatusName  string         `json:"status_name"`
	CreatedAt   time.Time      `json:"created_at" orm:"column(created_at)"`
	UserIn      bool           `json:"user_in"`
	Players     []lobby.Player `json:"players"`
	URL         string         `json:"URL"`
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
		games[i].StatusName = lobby.GameStatusName(game.Status)
	}

	return
}

func GetFinishedGames() (games []LobbyGame) {
	return getGames([]int{lobby.GameInactive})
}

func GetMyGames(userId int) (games []LobbyGame) {
	gamesAll := getGames(lobby.GetAllStatuses())

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
	return getGames([]int{lobby.GameActive, lobby.GameWait})
}

func GetActiveGames() (games []LobbyGame) {
	return getGames([]int{lobby.GameActive, lobby.GameWait})
}
