package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type AICheater struct {
	BaseAI
}

func NewAICheater(baseAI *BaseAI) *AICheater {
	ai := new(AICheater)
	ai.BaseAI = *baseAI
	return ai
}

func (ai *AICheater) GetAction() *game.Action {
	info := &ai.PlayerInfo
	myPos := info.CurrentPostion
	myCards := info.PlayerCards[myPos]
	for idx, card := range myCards {
		if card.Value == info.TableCards[card.Color].Value+1 {
			return game.NewAction(game.TypeActionPlaying, myPos, idx)
		}
	}

	magicConst := 0
	if info.BlueTokens < game.MaxBlueTokens && info.DeckSize > magicConst {
		for idx, card := range myCards {
			if !ai.isCardMayBeUsefull(card) {
				return game.NewAction(game.TypeActionDiscard, myPos, idx)
			}
		}
	}

	if info.BlueTokens > 0 {
		pos := (myPos + 1 + len(info.PlayerCards)) % len(info.PlayerCards)
		return game.NewAction(game.TypeActionInformationValue, pos, int(info.PlayerCards[pos][0].Value))
	}

	for idx, card := range myCards {
		if card.Value <= info.TableCards[card.Color].Value {
			return game.NewAction(game.TypeActionDiscard, myPos, idx)
		}
	}

	variants := map[game.ColorValue]int{}
	for idx, card := range myCards {
		colorValue := game.ColorValue{
			Color: card.Color,
			Value: card.Value,
		}
		if _, ok := variants[colorValue]; ok {
			return game.NewAction(game.TypeActionDiscard, myPos, idx)
		}
		variants[colorValue]++
	}

	variants = map[game.ColorValue]int{}
	for pos, cards := range info.PlayerCards {
		if pos == myPos {
			continue
		}

		for _, card := range cards {
			colorValue := game.ColorValue{
				Color: card.Color,
				Value: card.Value,
			}
			variants[colorValue]++
		}
	}

	for _, card := range info.Deck {
		colorValue := game.ColorValue{
			Color: card.Color,
			Value: card.Value,
		}
		variants[colorValue]++
	}

	for idx, card := range myCards {
		colorValue := game.ColorValue{
			Color: card.Color,
			Value: card.Value,
		}
		if _, ok := variants[colorValue]; ok {
			return game.NewAction(game.TypeActionDiscard, myPos, idx)
		}
	}

	return ai.getActionSmartyRandom()
}
