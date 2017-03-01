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
	card.SetKnown(true)

	if state.TableCards[card.Color].Value+1 == card.Value {
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
		state.UsedCards = append(state.UsedCards, card)
	}

	playerState := &state.PlayerStates[playerPosition]
	playerState.PlayerCards = append(playerState.PlayerCards[:cardPosition], playerState.PlayerCards[cardPosition+1:]...)
	if len(state.Deck) > 0 {
		newCard := state.Deck[0]
		state.Deck = state.Deck[1:]
		playerState.PlayerCards = append(playerState.PlayerCards, newCard)
	}

	action := state.NewAction(TypeActionPlaying, playerPosition, cardPosition)
	if len(state.Deck) == 0 && state.MaxStep == 0 {
		state.MaxStep = state.Step + len(state.PlayerStates)
	}
	return action, nil
}
