package game

import "errors"

func (game *Game) NewActionDiscard(playerPosition int, cardPosition int) (Action, error) {
	return game.AppendAction(game.CurrentState.NewActionDiscard(playerPosition, cardPosition))
}

func (state *GameState) NewActionDiscard(playerPosition int, cardPosition int) (Action, error) {
	if state.BlueTokens == MaxBlueTokens {
		return Action{}, errors.New("Too many blue tokens")
	}

	for i := 0; i < state.PlayerCount; i++ {
		playerState := &state.PlayerStates[i]
		if i == 0 {
			oldCard := playerState.PlayerCards[cardPosition]
			oldCard.SetKnown(true)
			state.UsedCards = append(state.UsedCards, oldCard)
		}
		playerState.PlayerCards = append(playerState.PlayerCards[:cardPosition], playerState.PlayerCards[cardPosition+1:]...)
	}

	if len(state.Deck) > 0 {
		card := state.Deck[0]
		state.Deck = state.Deck[1:]
		for i := 0; i < state.PlayerCount; i++ {
			playerState := &state.PlayerStates[i]
			card.SetKnown(i != playerPosition)
			playerState.PlayerCards = append(playerState.PlayerCards, card)
		}
	}

	state.BlueTokens++
	action := state.NewAction(TypeActionDiscard, playerPosition, cardPosition)
	return action, nil
}
