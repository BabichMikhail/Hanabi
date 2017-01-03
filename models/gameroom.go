package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BabichMikhail/Hanabi/engine"
	"github.com/astaxie/beego/orm"
)

type ActiveGame struct {
	Id      int       `orm:"auto"`
	GameID  int       `orm:"column(game_id)"`
	Json    string    `orm:"column(game);type(text)"`
	Created time.Time `orm:"column(created_at);auto_now_add;type(timestamp)"`
}

func (g *ActiveGame) TableName() string {
	return "active_games"
}

func ReadActiveGameById(id int) (game engine.Game, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("game").
		From("active_games").
		Where("id = ?")
	sql := qb.String()
	var jsonString string
	err = o.Raw(sql, id).QueryRow(&jsonString)
	json.Unmarshal([]byte(jsonString), &game)
	fmt.Print(err)
	return game, err
}

func ReadActiveGameByGameId(gameId int) (game engine.Game, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("game").
		From("active_games").
		Where("game_id = ?")
	sql := qb.String()
	var jsonString string
	err = o.Raw(sql, gameId).QueryRow(&jsonString)
	json.Unmarshal([]byte(jsonString), &game)
	return game, err
}

func CreateActiveGame(playerIds []int, gameId int) (game engine.Game, err error) {
	game = engine.NewGame(playerIds)
	o := orm.NewOrm()
	o.Begin()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("active_games", "game_id", "game", "created_at").
		Values("?", "?", "CURRENT_TIMESTAMP")
	sql := qb.String()
	if res, err := o.Raw(sql, gameId, game.SprintGame()).Exec(); err == nil {
		id64, _ := res.LastInsertId()
		id := int(id64)
		o.Commit()
		return ReadActiveGameById(id)
	}
	o.Rollback()
	return
}
