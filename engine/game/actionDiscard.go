package game

import "errors"

func (this *Game) NewActionDiscard(playerPosition int, cardPosition int) error {
	state := &this.CurrentState
	if state.BlueTokens == MaxBlueTokens {
		return errors.New("Too many blue tokens")
	}

	for i := 0; i < state.PlayerCount; i++ {
		playerState := &state.PlayerStates[i]
		if i == 0 {
			oldCard := playerState.PlayersCards[playerPosition][cardPosition]
			oldCard.SetKnown(true)
			state.UsedCards = append(state.UsedCards, oldCard)
		}
		playerState.PlayersCards[playerPosition] = append(playerState.PlayersCards[playerPosition][:cardPosition], playerState.PlayersCards[playerPosition][cardPosition+1:]...)
	}

	if len(state.Deck) > 0 {
		card := state.Deck[0]
		state.Deck = state.Deck[1:]
		for i := 0; i < state.PlayerCount; i++ {
			playerState := &state.PlayerStates[i]
			card.SetKnown(i != playerPosition)
			playerState.PlayersCards[playerPosition] = append(playerState.PlayersCards[playerPosition], card)
		}
	}

	state.BlueTokens++
	this.NewAction(TypeActionDiscard, playerPosition, cardPosition)
	return nil
}
