package game

import "errors"

func (game *Game) NewActionPlaying(playerPosition int, cardPosition int) (Action, error) {
	return game.AppendAction(game.CurrentState.NewActionPlaying(playerPosition, cardPosition))
}

func (state *GameState) NewActionPlaying(playerPosition int, cardPosition int) (Action, error) {
	if state.RedTokens == 0 {
		return Action{}, errors.New("No red tokens")
	}

	card := state.PlayerStates[playerPosition].PlayerCards[cardPosition]

	var success bool
	if success = state.TableCards[card.Color].Value+1 == card.Value; success {
		card.SetKnown(true)
		state.TableCards[card.Color] = Card{
			Color:      card.Color,
			KnownColor: true,
			Value:      card.Value,
			KnownValue: true,
		}
		if card.Value == Five && state.BlueTokens < MaxBlueTokens {
			state.BlueTokens++
		}
	} else {
		state.RedTokens--
	}

	for i := 0; i < state.PlayerCount; i++ {
		playerState := &state.PlayerStates[i]
		if i == playerPosition {
			oldCard := playerState.PlayerCards[cardPosition]
			oldCard.SetKnown(true)
			if !success {
				state.UsedCards = append(state.UsedCards, oldCard)
			}
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

	action := state.NewAction(TypeActionPlaying, playerPosition, cardPosition)
	return action, nil
}
