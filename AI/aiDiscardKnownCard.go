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

func (ai *AIDiscardKnownCard) GetAction() *game.Action {
	ai.setAvailableActions()
	ai.setAvailableInformation()
	info := &ai.PlayerInfo
	pos := info.CurrentPosition
	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[pos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, pos, idx)
			}
		}

		for idx, card := range info.PlayerCards[pos] {
			if card.KnownValue && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, pos, idx)
			}
		}
	}

	return ai.getActionSmartyRandom()
}
