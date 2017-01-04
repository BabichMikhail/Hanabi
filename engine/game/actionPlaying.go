package game

import "errors"

func (this Game) NewActionPlaying(playerPosition int, cardPosition int) error {
	state := &this.CurrentState
	if state.RedTokens == 0 {
		return errors.New("No red tokens")
	}

	card := state.PlayerStates[0].PlayersCards[playerPosition][cardPosition]
	card.SetKnown(true)
	if state.TableCards[int(card.Color)].Value+1 == card.Value {
		state.TableCards[int(card.Color)] = Card{
			Color:      card.Color,
			KnownColor: true,
			Value:      card.Value,
			KnownValue: true,
		}
	} else {
		state.RedTokens--
	}

	this.NewAction(TypeActionPlaying, playerPosition, cardPosition)
	return nil
}
