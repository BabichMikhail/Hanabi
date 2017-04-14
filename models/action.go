package models

import (
	gamePackage "github.com/BabichMikhail/Hanabi/game"
	"github.com/astaxie/beego/orm"
)

type Action struct {
	Id       int   `orm:"auto"`
	GameId   int   `orm:"column(game_id)"`
	Game     *Game `orm:"-"`
	Type     int   `orm:"column(action_type)"`
	Position int   `orm:"column(player_position)"`
	Value    int   `orm:"column(value)"`
}

func (a *Action) TableIndex() [][]string {
	return [][]string{
		[]string{"game_id"},
	}
}

func (a *Action) TableName() string {
	return "actions"
}

func NewAction(gameId int, action *gamePackage.Action) (err error) {
	o := orm.NewOrm()
	dbAction := &Action{
		Position: action.PlayerPosition,
		Type:     int(action.ActionType),
		Value:    action.Value,
		GameId:   gameId,
	}
	_, err = o.Insert(dbAction)
	return
}

func ReadActions(gameId int) ([]gamePackage.Action, error) {
	o := orm.NewOrm()
	var actions []Action
	count, err := o.QueryTable(Action{}).Filter("game_id", gameId).OrderBy("id").All(&actions)
	gameActions := make([]gamePackage.Action, count)
	for i, action := range actions {
		gameActions[i] = gamePackage.Action{
			ActionType:     gamePackage.ActionType(action.Type),
			PlayerPosition: action.Position,
			Value:          action.Value,
		}
	}
	return gameActions, err
}

func GetActionCount(gameId int) (int, error) {
	o := orm.NewOrm()
	count64, err := o.QueryTable(Action{}).Filter("game_id", gameId).OrderBy("id").Count()
	count := int(count64)
	return count, err
}
