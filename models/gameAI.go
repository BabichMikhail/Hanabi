package models

import (
	"errors"
	"regexp"
	"strconv"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/game"
	"github.com/astaxie/beego/orm"
	wetalk "github.com/beego/wetalk/modules/models"
)

func ApplyAction(gameId int, actionType game.ActionType, playerPosition int, actionValue int) (err error) {
	state, err := ReadCurrentGameState(gameId)
	if err != nil {
		return
	}

	var action game.Action
	switch actionType {
	case game.TypeActionDiscard:
		action, _ = state.NewActionDiscard(playerPosition, actionValue)
	case game.TypeActionInformationColor:
		action, _ = state.NewActionInformationColor(playerPosition, game.CardColor(actionValue))
	case game.TypeActionInformationValue:
		action, _ = state.NewActionInformationValue(playerPosition, game.CardValue(actionValue))
	case game.TypeActionPlaying:
		action, _ = state.NewActionPlaying(playerPosition, actionValue)
	}

	NewAction(gameId, action)
	UpdateGameState(gameId, state)
	return
}

func CheckAI(gameId int) {
	state, err := ReadCurrentGameState(gameId)
	if state.IsGameOver() {
		SetGameFinishedStatus(gameId)
		return
	}

	if err != nil {
		return
	}
	pos := state.CurrentPosition
	playerId := state.PlayerStates[pos].PlayerId
	playerInfo := state.GetPlayerGameInfo(playerId)
	if ok, _ := regexp.MatchString("AI_.*", GetUserNickNameById(playerId)); !ok {
		return
	}
	AI := ai.NewAI(playerInfo, ai.AI_RandomAction)
	action := AI.GetAction()
	ApplyAction(gameId, action.ActionType, action.PlayerPosition, action.Value)

	CheckAI(gameId)
}

type AIUser struct {
	Id     int `orm:"auto"`
	Type   int `orm:"column(type)"`
	UserId int `orm:"column(user_id)"`
}

func (user *AIUser) TableName() string {
	return "ai_users"
}

func GetAIUserIds(AIType, playerCount int) (ids []int, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("u.id").From("user u").
		InnerJoin("ai_users aiu").On("u.id = aiu.user_id").
		Where("u.user_name LIKE 'AI\\_%'").
		And("aiu.type = ?").
		Limit(playerCount)
	var users []wetalk.User
	_, err = o.Raw(qb.String(), AIType).QueryRows(&users)
	if len(users) == 0 {
		return []int{}, errors.New("AI Users not found")
	}

	for _, user := range users {
		ids = append(ids, user.Id)
	}
	return ids, err
}

func CreateAIUsers(AIType int) (ids []int, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.InsertInto("user", "user_name", "nick_name", "email", "created", "updated").
		Values("?", "?", "?", "CURRENT_TIMESTAMP", "CURRENT_TIMESTAMP")
	username := ai.DefaultUsernamePrefix(AIType)
	o.Begin()
	var id int
	for i := 0; i < 5; i++ {
		newUsername := username + "_" + strconv.Itoa(i)
		if res, errExec := o.Raw(qb.String(), newUsername, newUsername, newUsername+"@notmail").Exec(); errExec != nil {
			o.Rollback()
			return []int{}, errExec
		} else {
			id64, _ := res.LastInsertId()
			id = int(id64)
			ids = append(ids, id)
		}

		qbAI, _ := orm.NewQueryBuilder("mysql")
		qbAI.InsertInto("ai_users", "type", "user_id").Values("?", "?")
		if _, errExec := o.Raw(qbAI.String(), AIType, id).Exec(); errExec != nil {
			o.Rollback()
			return []int{}, errExec
		}
	}
	o.Commit()
	return ids, nil

}
