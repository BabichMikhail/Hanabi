package models

import (
	"regexp"

	ai "github.com/BabichMikhail/Hanabi/AI"
)

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
	ai.NewAI(playerInfo, ai.AI_RandomAction)
}
