package game

import "errors"

func (this *Game) NewActionInformation(playerPosition int, value int, actionType ActionType, actionFunc func(card *Card, value int)) error {
	state := &this.CurrentState
	if state.BlueTokens == 0 {
		return errors.New("No blue tokens")
	}

	for i := 0; i < len(state.PlayerStates[playerPosition].PlayersCards[playerPosition]); i++ {
		actionFunc(&state.PlayerStates[playerPosition].PlayersCards[playerPosition][i], value)
	}

	state.BlueTokens--
	this.NewAction(TypeActionInformationColor, playerPosition, int(value))
	return nil
}

func (this *Game) NewActionInformationColor(playerPosition int, cardColor CardColor) error {
	return this.NewActionInformation(playerPosition, int(cardColor), TypeActionInformationColor, func(card *Card, value int) {
		if card.Color == CardColor(value) {
			card.KnownColor = true
		}
	})
}

func (this *Game) NewActionInformationValue(playerPosition int, cardValue CardValue) error {
	return this.NewActionInformation(playerPosition, int(cardValue), TypeActionInformationValue, func(card *Card, value int) {
		if (*card).Value == CardValue(value) {
			(*card).KnownValue = true
		}
	})
}