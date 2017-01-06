package models

import (
	"encoding/json"
	"errors"
	"fmt"

	gamePackage "github.com/BabichMikhail/Hanabi/engine/game"
	"github.com/astaxie/beego/orm"
)

func CreateActiveGame(playerIds []int, gameId int) (game gamePackage.Game, err error) {
	game = gamePackage.NewGame(playerIds)
	o := orm.NewOrm()
	o.Begin()
	var ormGame Game
	_, err = o.QueryTable(ormGame).Filter("id", gameId).Update(orm.Params{
		"game": game.SprintGame(),
	})
	if err == nil {
		o.Commit()
		return game, nil
	}
	o.Rollback()
	return
}

func ReadActiveGameById(id int) (game gamePackage.Game, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("game").
		From("games").
		Where("id = ?")
	sql := qb.String()
	var gameModel Game
	err = o.Raw(sql, id).QueryRow(&gameModel)
	if fmt.Sprintf("%s", gameModel.Json) == "" {
		err = errors.New(fmt.Sprintf("Active game #%d not found", id))
		return
	}
	json.Unmarshal([]byte(gameModel.Json), &game)
	return
}

func UpdateCurrentGameById(gameId int, game gamePackage.Game) (err error) {
	o := orm.NewOrm()
	activeGame := new(Game)
	activeGame.Id = gameId
	if err = o.Read(activeGame); err != nil {
		return err
	}
	activeGame.Json = game.SprintGame()
	_, err = o.Update(activeGame, "game")
	return
}
