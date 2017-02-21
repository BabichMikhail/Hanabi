package models

import (
	"encoding/json"

	gamePackage "github.com/BabichMikhail/Hanabi/game"
	"github.com/astaxie/beego/orm"
)

type GameState struct {
	Id     int    `orm:"auto"`
	GameId int    `orm:"column(game_id)"`
	Game   *Game  `orm:"-"`
	IsInit bool   `orm:"column(is_init_state)"`
	Json   string `orm:"column(json);type(text)"`
}

func (gs *GameState) TableName() string {
	return "game_states"
}

func (gs *GameState) TableIndex() [][]string {
	return [][]string{
		[]string{"game_id"},
	}
}

func (gs *GameState) TableUnique() [][]string {
	return [][]string{
		[]string{"is_init_state", "game_id"},
	}
}

func NewGameState(gameId int, gameState *gamePackage.GameState, isInit bool) error {
	o := orm.NewOrm()
	state := &GameState{
		GameId: gameId,
		IsInit: isInit,
		Json:   gameState.Sprint(),
	}
	_, err := o.Insert(state)
	return err
}

func ReadCurrentGameState(gameId int) (gameState gamePackage.GameState, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("json").From("game_states").Where("game_id = ?").And("is_init_state = ?")
	var state GameState

	if err = o.Raw(qb.String(), gameId, false).QueryRow(&state); err == nil {
		err = json.Unmarshal([]byte(state.Json), &gameState)
	}
	return
}

func ReadInitialGameState(gameId int) (gameState gamePackage.GameState, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("json").From("game_states").Where("game_id = ?").And("is_init_state = ?")
	var state GameState
	if err = o.Raw(qb.String(), gameId, true).QueryRow(&state); err == nil {
		err = json.Unmarshal([]byte(state.Json), &gameState)
	}
	return
}

func UpdateGameState(gameId int, gameState gamePackage.GameState) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("game_states").Filter("game_id", gameId).Filter("is_init_state", false).Update(orm.Params{
		"json": gameState.Sprint(),
	})
	return err

}
