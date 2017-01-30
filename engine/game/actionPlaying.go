package game

import "errors"

func (this *Game) NewActionPlaying(playerPosition int, cardPosition int) error {
	state := &this.CurrentState
	if state.RedTokens == 0 {
		return errors.New("No red tokens")
	}

	card := state.PlayerStates[0].PlayersCards[playerPosition][cardPosition]

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
			oldCard := playerState.PlayersCards[playerPosition][cardPosition]
			oldCard.SetKnown(true)
			if !success {
				state.UsedCards = append(state.UsedCards, oldCard)
			}
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

	this.NewAction(TypeActionPlaying, playerPosition, cardPosition)
	return nil
}
