package ai

import (
	"strconv"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI1 struct {
	BaseAI
}

func (ai *AI1) GetAction() *game.Action {
	info := &ai.PlayerInfo
	info.SetProbabilities(false, false)
	myPos := info.CurrentPostion

	for idx, card := range info.PlayerCards[myPos] {
		if card.KnownColor && card.KnownValue && info.TableCards[card.Color].Value+1 == card.Value {
			return game.NewAction(game.TypeActionPlaying, myPos, idx)
		}
	}

	if len(info.PlayerCards) != 5 {
		panic("Not implemented")
	}

	if len(ai.History) > 0 {
		action := ai.History[len(ai.History)-1]
		isInformationAction := action.ActionType == game.TypeActionInformationColor || action.ActionType == game.TypeActionInformationValue
		if isInformationAction {
			cardPos := (action.PlayerPosition - myPos + len(info.PlayerCards)) % len(info.PlayerCards)
			if ai.isCardPlayable(info.PlayerCards[myPos][cardPos]) {
				return game.NewAction(game.TypeActionPlaying, myPos, cardPos)
			}
		}
	}

	nextPos := (myPos + len(info.PlayerCards) + 1) % len(info.PlayerCards)
	if info.BlueTokens > 0 {
		for idx, card := range info.PlayerCards[nextPos] {
			if ai.isCardPlayable(card) {
				infoPos := (myPos + len(info.PlayerCards) + idx + 1) % len(info.PlayerCards)
				cards := info.PlayerCardsInfo[infoPos]
				for j := 0; j < len(cards); j++ {
					if !cards[j].KnownValue && !cards[j].KnownColor {
						return game.NewAction(game.TypeActionInformationValue, infoPos, int(info.PlayerCards[infoPos][j].Value))
					}
				}

				for j := 0; j < len(cards); j++ {
					if !cards[j].KnownValue {
						return game.NewAction(game.TypeActionInformationValue, infoPos, int(info.PlayerCards[infoPos][j].Value))
					}
				}

				for j := 0; j < len(cards); j++ {
					if !cards[j].KnownColor {
						return game.NewAction(game.TypeActionInformationColor, infoPos, int(info.PlayerCards[infoPos][j].Color))
					}
				}

				return game.NewAction(game.TypeActionInformationColor, infoPos, int(info.PlayerCards[infoPos][0].Color))
			}
		}
	}

	if info.BlueTokens > 4 {
		cardsInfo := info.PlayerCardsInfo[nextPos]
		for idx, card := range cardsInfo {
			infoIdx := -1
			if !card.KnownValue {
				copyInfo := info.Copy()
				copyInfo.SetProbabilities(false, false)
				copyInfo.PreviewActionInformationValue(nextPos, info.PlayerCards[nextPos][idx].Value)
				copyInfo.SetProbabilities(false, false)
				if !ai.isCardPlayable(copyInfo.PlayerCardsInfo[nextPos][idx]) {
					infoIdx = idx
				}
			}
			if !card.KnownColor {
				copyInfo := info.Copy()
				copyInfo.SetProbabilities(false, false)
				copyInfo.PreviewActionInformationColor(nextPos, info.PlayerCards[nextPos][idx].Color)
				copyInfo.SetProbabilities(false, false)
				if !ai.isCardPlayable(copyInfo.PlayerCardsInfo[nextPos][idx]) {
					infoIdx = idx
				}
			}

			if infoIdx == -1 {
				continue
			}

			infoPos := (myPos + len(info.PlayerCards) + 1 + infoIdx) % len(info.PlayerCards)
			for cardPos, card := range info.PlayerCardsInfo[infoPos] {
				if !card.KnownValue {
					return game.NewAction(game.TypeActionInformationValue, infoPos, int(info.PlayerCards[infoPos][cardPos].Value))
				} else if !card.KnownColor {
					return game.NewAction(game.TypeActionInformationColor, infoPos, int(info.PlayerCards[infoPos][cardPos].Color))
				}
			}
		}
	}

	if info.BlueTokens <= 4 {
		for idx, card := range info.PlayerCards[myPos] {
			if !ai.isCardMayBeUsefull(card) {
				return game.NewAction(game.TypeActionDiscard, myPos, idx)
			}
		}
		return game.NewAction(game.TypeActionDiscard, myPos, 0)
	}

	if info.BlueTokens != game.MaxBlueTokens {
		panic("Abc")
	}

	panic("Magic " + strconv.Itoa(info.BlueTokens))
}

func NewAI1(baseAI *BaseAI) *AI1 {
	ai := new(AI1)
	ai.BaseAI = *baseAI
	return ai
}
