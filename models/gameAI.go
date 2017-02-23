package models

import (
	"regexp"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/game"
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
}
