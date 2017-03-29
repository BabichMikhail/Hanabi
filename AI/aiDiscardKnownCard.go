package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type AIDiscardKnownCard struct {
	BaseAI
}

func NewAIDiscardKnownCard(baseAI *BaseAI) *AIDiscardKnownCard {
	ai := new(AIDiscardKnownCard)
	ai.BaseAI = *baseAI
	return ai
}

func (ai *AIDiscardKnownCard) GetAction() game.Action {
	ai.setAvailableInfomation()
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
