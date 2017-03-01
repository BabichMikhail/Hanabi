package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Game struct {
	Id           int        `orm:"auto" json:"id"`
	OwnerId      int        `orm:"column(owner_id)" json:"owner_id"`
	PlayerCount  int        `orm:"column(player_count)" json:"player_count"`
	Status       int        `orm:"column(status);default(4)" json:"status"`
	Points       int        `orm:"column(points);default(0)"`
	Seed         int64      `orm:"column(seed);null"`
	IsAI         bool       `orm:"column(is_ai_only_game)"`
	InitState    *GameState `orm:"-"`
	CurrentState *GameState `orm:"-"`
	Actions      []*Action  `orm:"-"`
	CreatedAt    time.Time  `orm:"column(created_at);auto_now_add;type(timestamp)" json:"created_at"`
}

func (g *Game) TableName() string {
	return "games"
}

func NewGame(userId int, playersCount int, status int, isAIGame bool) (gameItem LobbyGameItem, err error) {
	if playersCount > 5 {
		playersCount = 5
	}
	if playersCount < 2 {
		playersCount = 2
	}
	o := orm.NewOrm()
	o.Begin()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("games", "owner_id", "player_count", "status", "is_ai_only_game", "created_at").
		Values("?", "?", "?", "?", "CURRENT_TIMESTAMP")
	sql := qb.String()
	var id int
	if res, errExec := o.Raw(sql, userId, playersCount, status, isAIGame).Exec(); errExec == nil {
		id64, _ := res.LastInsertId()
		id = int(id64)
	} else {
		o.Rollback()
		err = errExec
		return
	}

	if err, _ = JoinGame(id, userId); err == nil {
		userNickname := GetUserNickNameById(userId)
		gameItem = LobbyGameItem{
			Id:          id,
			Owner:       userNickname,
			OwnerId:     userId,
			PlayerCount: playersCount,
			Status:      StatusName(status),
			StatusCode:  status,
			UserIn:      true,
		}
		gameItem.Players = append(gameItem.Players, LobbyPlayer{
			Id:       userId,
			NickName: userNickname,
		})
		o.Commit()
	} else {
		o.Rollback()
		err = errors.New("Can't create new active game")
	}
	return
}

func JoinGame(gameId int, userId int) (err error, status string) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	// @todo insert first player right
	//o.Begin()
	qb.InsertInto("players", "user_id", "game_id").
		Values("?", "?")
	sql := qb.String()
	if _, err = o.Raw(sql, userId, gameId).Exec(); err == nil {
		gamePlayers := GetGamePlayers([]int{gameId})[gameId]
		qbGames, _ := orm.NewQueryBuilder("mysql")
		sql := qbGames.Select("player_count").From("games").Where("id = ?").String()
		var game Game
		o.Raw(sql, gameId).QueryRow(&game)
		status = StatusName(StatusWait)
		if game.PlayerCount == len(gamePlayers) {
			playerIds := []int{}
			for i := 0; i < len(gamePlayers); i++ {
				playerIds = append(playerIds, gamePlayers[i].Id)
			}
			_, err = CreateActiveGame(playerIds, gameId)
			status = StatusName(StatusActive)
		}
		//o.Commit()
	} else {
		status = ""
		//o.Rollback()
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

	if game.Status != int(StatusWait) {
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
