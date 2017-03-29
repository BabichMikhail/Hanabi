package ai

import (
	"math/rand"

	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInformation struct {
	BaseAI
}

func NewAIUsefulInformation(baseAI *BaseAI) *AIUsefulInformation {
	ai := new(AIUsefulInformation)
	ai.BaseAI = *baseAI
	return ai
}

func (ai *AIUsefulInformation) GetAction() game.Action {
	ai.setAvailableInfomation()
	info := &ai.PlayerInfo
	myPos := info.Position

	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[myPos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, myPos, idx)
			}
		}
	}

	for _, action := range ai.History[Max(len(ai.History)-len(info.PlayerCards)-1, 0):] {
		if action.ActionType == game.TypeActionInformationValue && action.PlayerPosition == myPos {
			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					return game.NewAction(game.TypeActionPlaying, myPos, idx)
				}
			}
		}

		if action.ActionType == game.TypeActionInformationColor && action.PlayerPosition == myPos {
			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					return game.NewAction(game.TypeActionPlaying, myPos, idx)
				}
			}
		}
	}

	if info.BlueTokens > 0 {
		for i := 0; i < len(info.PlayerCards)-1; i++ {
			nextPos := (myPos + i) % len(info.PlayerCards)
			for color, tableCard := range info.TableCards {
				for idx, card := range info.PlayerCards[nextPos] {
					if card.Color == color && card.Value == tableCard.Value+1 {
						cardInfo := &info.PlayerCardsInfo[nextPos][idx]
						if !cardInfo.KnownValue {
							return game.NewAction(game.TypeActionInformationValue, nextPos, int(card.Value))
						} else if !cardInfo.KnownColor {
							return game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color))
						}
					}
				}
			}
		}
	}

	if info.BlueTokens < game.MaxBlueTokens {
		for idx, card := range info.PlayerCards[myPos] {
			if !card.KnownValue {
				return game.NewAction(game.TypeActionDiscard, myPos, idx)
			}
		}
		return game.NewAction(game.TypeActionDiscard, myPos, rand.Intn(len(info.PlayerCards[myPos])))
	}
	return ai.getActionSmartyRandom()
}
