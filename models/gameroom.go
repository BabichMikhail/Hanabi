package models

import (
	gamePackage "github.com/BabichMikhail/Hanabi/engine/game"
	lobby "github.com/BabichMikhail/Hanabi/engine/lobby"
	"github.com/astaxie/beego/orm"
)

func CreateActiveGame(playerIds []int, gameId int) (game *gamePackage.Game, err error) {
	game = gamePackage.NewGame(playerIds)
	o := orm.NewOrm()
	o.Begin()
	var ormGame Game
	_, err = o.QueryTable(ormGame).Filter("id", gameId).Update(orm.Params{
		"seed": game.Seed,
	})
	if err != nil {
		o.Rollback()
		return game, err
	}

	if err = NewGameState(gameId, &game.InitState, true); err != nil {
		o.Rollback()
		return game, err
	}

	if err = NewGameState(gameId, &game.CurrentState, false); err != nil {
		o.Rollback()
		return
	}

	o.Commit()
	return game, err
}

func SetGameInactiveStatus(gameId int) {
	o := orm.NewOrm()
	var ormGame Game
	o.QueryTable(ormGame).Filter("id", gameId).Update(orm.Params{
		"status": lobby.GameInactive,
	})
}
