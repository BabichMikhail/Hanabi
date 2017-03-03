package ai

import (
	"math/rand"

	"github.com/BabichMikhail/Hanabi/game"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (ai *AI) getUsefullInformationAction() *Action {
	ai.SetAvailableInfomation()
	info := &ai.PlayerInfo
	myPos := info.Position

	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[myPos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return NewAction(game.TypeActionPlaying, myPos, idx, 1, 1)
			}
		}
	}

	for _, action := range ai.History[Max(len(ai.History)-len(info.PlayerCards)-1, 0):] {
		if action.ActionType == game.TypeActionInformationValue && action.PlayerPosition == myPos {
			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					return NewAction(game.TypeActionPlaying, myPos, idx, 1, 1)
				}
			}
		}

		if action.ActionType == game.TypeActionInformationColor && action.PlayerPosition == myPos {
			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					return NewAction(game.TypeActionPlaying, myPos, idx, 1, 1)
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
							return NewAction(game.TypeActionInformationValue, nextPos, int(card.Value), 1, 1)
						} else if !cardInfo.KnownColor {
							return NewAction(game.TypeActionInformationColor, nextPos, int(card.Color), 1, 1)
						}
					}
				}
			}
		}
	}

	if info.BlueTokens < game.MaxBlueTokens {
		for idx, card := range info.PlayerCards[myPos] {
			if !card.KnownValue {
				return NewAction(game.TypeActionDiscard, myPos, idx, 1, 1)
			}
		}
		return NewAction(game.TypeActionDiscard, myPos, rand.Intn(len(info.PlayerCards[myPos])), 1, 1)
	}
	return ai.getSmartyRandomAction()
}
