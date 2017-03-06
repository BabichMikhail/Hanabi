package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

func (ai *AI) getActionDiscardUsefullCard() game.Action {
	ai.SetAvailableInfomation()
	info := &ai.PlayerInfo
	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[info.Position] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, info.Position, idx)
			}
		}

		for idx, card := range info.PlayerCards[info.Position] {
			if card.KnownValue && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, info.Position, idx)
			}
		}
	}

	return ai.getActionSmartyRandom()
}
