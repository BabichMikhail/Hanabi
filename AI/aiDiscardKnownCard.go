package ai

import (
	"math/rand"
	"time"

	"github.com/BabichMikhail/Hanabi/game"
)

func (ai *AI) getDiscardUsefullCardAction() *Action {
	ai.SetAvailableInfomation()
	info := &ai.PlayerInfo
	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[info.Position] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return NewAction(game.TypeActionDiscard, info.Position, idx, 1, 1)
			}
		}

		for idx, card := range info.PlayerCards[info.Position] {
			if card.KnownValue && card.Value == tableCard.Value+1 {
				return NewAction(game.TypeActionDiscard, info.Position, idx, 1, 1)
			}
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	return ai.Actions[rand.Intn(len(ai.Actions))]
}
