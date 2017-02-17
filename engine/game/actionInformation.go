package game

import "errors"

func (game *Game) NewActionInformation(playerPosition int, value int, actionType ActionType, actionFunc func(card *Card, value int)) (Action, error) {
	return game.CurrentState.NewActionInformation(playerPosition, value, actionType, actionFunc)
}

func (state *GameState) NewActionInformation(playerPosition int, value int, actionType ActionType, actionFunc func(card *Card, value int)) (Action, error) {
	if state.BlueTokens == 0 {
		return Action{}, errors.New("No blue tokens")
	}

	for i := 0; i < len(state.PlayerStates[playerPosition].PlayersCards[playerPosition]); i++ {
		actionFunc(&state.PlayerStates[playerPosition].PlayersCards[playerPosition][i], value)
	}

	state.BlueTokens--
	action := state.NewAction(TypeActionInformationColor, playerPosition, int(value))
	return action, nil
}

func (game *Game) NewActionInformationColor(playerPosition int, cardColor CardColor) (Action, error) {
	return game.NewActionInformation(playerPosition, int(cardColor), TypeActionInformationColor, func(card *Card, value int) {
		if card.Color == CardColor(value) {
			card.KnownColor = true
		}
	})
}

func (game *Game) NewActionInformationValue(playerPosition int, cardValue CardValue) (Action, error) {
	return game.NewActionInformation(playerPosition, int(cardValue), TypeActionInformationValue, func(card *Card, value int) {
		if (*card).Value == CardValue(value) {
			(*card).KnownValue = true
		}
	})
}
