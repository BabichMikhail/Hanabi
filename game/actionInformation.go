package game

import (
	"errors"
)

func (state *GameState) NewActionInformation(playerPosition int, value int, actionType ActionType, actionFunc func(card *Card, value int)) (Action, error) {
	if state.BlueTokens == 0 {
		return Action{}, errors.New("No blue tokens")
	}

	for i := 0; i < len(state.PlayerStates[playerPosition].PlayerCards); i++ {
		actionFunc(&state.PlayerStates[playerPosition].PlayerCards[i], value)
	}

	state.BlueTokens--
	action := state.NewAction(actionType, playerPosition, int(value))
	return action, nil
}

func (state *GameState) NewActionInformationColor(playerPosition int, cardColor CardColor) (Action, error) {
	return state.NewActionInformation(playerPosition, int(cardColor), TypeActionInformationColor, func(card *Card, color int) {
		if card.Color == CardColor(color) {
			card.KnownColor = true
			if _, ok := card.ProbabilityColors[CardColor(color)]; ok {
				delete(card.ProbabilityColors, CardColor(color))
			}
		}
	})
}

func (state *GameState) NewActionInformationValue(playerPosition int, cardValue CardValue) (Action, error) {
	return state.NewActionInformation(playerPosition, int(cardValue), TypeActionInformationValue, func(card *Card, value int) {
		if card.Value == CardValue(value) {
			card.KnownValue = true
			if _, ok := card.ProbabilityValues[CardValue(value)]; ok {
				delete(card.ProbabilityValues, CardValue(value))
			}
		}
	})
}

func (game *Game) NewActionInformationColor(playerPosition int, cardColor CardColor) (Action, error) {
	return game.AppendAction(game.CurrentState.NewActionInformationColor(playerPosition, cardColor))
}

func (game *Game) NewActionInformationValue(playerPosition int, cardValue CardValue) (Action, error) {
	return game.AppendAction(game.CurrentState.NewActionInformationValue(playerPosition, cardValue))
}
