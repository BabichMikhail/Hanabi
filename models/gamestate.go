package models

import (
	"encoding/json"
	"fmt"

	gamePackage "github.com/BabichMikhail/Hanabi/engine/game"
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
	bytes, _ := json.Marshal(gameState)
	state := &GameState{
		GameId: gameId,
		IsInit: isInit,
		Json:   fmt.Sprint(bytes),
	}
	_, err := o.Insert(state)
	return err
}

func ReadCurrentGameState(gameId int) (gameState gamePackage.GameState, err error) {
	o := orm.NewOrm()
	state := new(GameState)
	state.GameId = gameId
	state.IsInit = false
	if err = o.Read(state); err == nil {
		err = json.Unmarshal([]byte(state.Json), &gameState)
	}
	return
}

func ReadInitialGameState(gameId int) (gameState gamePackage.GameState, err error) {
	o := orm.NewOrm()
	state := new(GameState)
	state.GameId = gameId
	state.IsInit = true
	if err = o.Read(state); err == nil {
		err = json.Unmarshal([]byte(state.Json), &gameState)
	}
	return
}

func UpdateGameState(gameId int, gameState gamePackage.GameState) error {
	o := orm.NewOrm()
	state := new(GameState)
	state.GameId = gameId
	state.IsInit = false
	state.Json = fmt.Sprint(json.Marshal(gameState))
	_, err := o.Update(state, "json")
	return err
}
