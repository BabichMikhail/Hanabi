package ai

import (
	"math"
	"math/rand"
	"sort"

	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInformationV3 struct {
	BaseAI

	CoefPlayByValue           float64
	CoefPlayByColor           float64
	CoefInfoValue             float64
	CoefInfoColor             float64
	CoefDiscardUselessCard    float64
	CoefDiscardUnknownCard    float64
	CoefDiscardUsefulCard     float64
	CoefDiscardMaybeUsefuCard float64
}

func NewAIUsefulInformationV3(baseAI *BaseAI) *AIUsefulInformationV3 {
	ai := new(AIUsefulInformationV3)
	ai.BaseAI = *baseAI
	ai.CoefPlayByValue = 2.1
	ai.CoefPlayByColor = -0.9
	ai.CoefInfoValue = 1.05
	ai.CoefInfoColor = 1.0
	ai.CoefDiscardUsefulCard = 0.1
	ai.CoefDiscardMaybeUsefuCard = 0.04
	ai.CoefDiscardUselessCard = 0.01
	ai.CoefDiscardUnknownCard = 0.07
	return ai
}

func (ai *AIUsefulInformationV3) SetCoefs(kPlayByValue, kPlayByColor, kInfoValue, kInfoColor, kDiscardUseful, kDiscardMaybeUseful, kDiscardUseless, kDiscardUnknown float64) {
	ai.CoefPlayByValue = kPlayByValue
	ai.CoefPlayByColor = kPlayByColor
	ai.CoefInfoValue = kInfoValue
	ai.CoefInfoColor = kInfoColor
	ai.CoefDiscardUsefulCard = kDiscardUseful
	ai.CoefDiscardMaybeUsefuCard = kDiscardMaybeUseful
	ai.CoefDiscardUselessCard = kDiscardUseless
	ai.CoefDiscardUnknownCard = kDiscardUnknown
}

func (ai *AIUsefulInformationV3) GetAction() game.Action {
	ai.setAvailableInfomation()
	info := &ai.PlayerInfo
	myPos := info.CurrentPostion

	usefulActions := Actions{}

	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[myPos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, myPos, idx)
			}
		}
	}

	subHistory := ai.History[Max(len(ai.History)-len(info.PlayerCards)+1, 0):]
	historyLength := len(subHistory)
	for i, action := range subHistory {
		if action.ActionType == game.TypeActionInformationValue && action.PlayerPosition == myPos {
			count := 0.0
			for _, card := range info.PlayerCards[myPos] {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					count++
				}
			}

			if count == 0 {
				continue
			}

			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: ai.CoefPlayByValue / float64(historyLength-i) / math.Sqrt(count),
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}

		if action.ActionType == game.TypeActionInformationColor && action.PlayerPosition == myPos {
			count := 0.0
			for _, card := range info.PlayerCards[myPos] {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					count++
				}
			}

			if count == 0 {
				continue
			}

			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: ai.CoefPlayByColor / float64(historyLength-i) / math.Sqrt(count),
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}
	}

	if info.BlueTokens > 0 {
		for i := 1; i < len(info.PlayerCards); i++ {
			nextPos := (myPos + i) % len(info.PlayerCards)
			for color, tableCard := range info.TableCards {
				for idx, card := range info.PlayerCards[nextPos] {
					if card.Color == color && card.Value == tableCard.Value+1 {
						cardInfo := &info.PlayerCardsInfo[nextPos][idx]
						if !cardInfo.KnownValue {
							action := UsefulAction{
								Action:     game.NewAction(game.TypeActionInformationValue, nextPos, int(card.Value)),
								Usefulness: ai.CoefInfoValue * (1.0 - float64(i)/float64(len(info.PlayerCards))),
							}
							usefulActions = append(usefulActions, action)
						}

						if !cardInfo.KnownColor {
							action := UsefulAction{
								Action:     game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color)),
								Usefulness: ai.CoefInfoColor * (1.0 - float64(i)/float64(len(info.PlayerCards))),
							}
							usefulActions = append(usefulActions, action)
						}
					}
				}
			}
		}
	}

	if info.BlueTokens < game.MaxBlueTokens {
		for idx, card := range info.PlayerCards[myPos] {
			var coef float64
			if card.KnownColor && card.KnownValue {
				if card.Value > info.TableCards[card.Color].Value {
					coef = ai.CoefDiscardUsefulCard
				} else {
					coef = ai.CoefDiscardUselessCard
				}
			} else if card.KnownValue {
				coef = ai.CoefDiscardUselessCard
				for _, card := range info.TableCards {
					if card.Value+1 == card.Value {
						coef = ai.CoefDiscardMaybeUsefuCard
					}
				}
			} else if card.KnownColor {
				if info.TableCards[card.Color].Value == 5 {
					coef = ai.CoefDiscardUselessCard
				} else {
					coef = ai.CoefDiscardMaybeUsefuCard
				}
			} else {
				coef = ai.CoefDiscardUnknownCard
			}
			action := UsefulAction{
				Action:     game.NewAction(game.TypeActionDiscard, myPos, idx),
				Usefulness: coef,
			}
			usefulActions = append(usefulActions, action)
		}
	}

	if len(usefulActions) > 0 {
		sort.Sort(usefulActions)
		return usefulActions[0].Action
	}

	if info.BlueTokens < game.MaxBlueTokens {
		return game.NewAction(game.TypeActionDiscard, myPos, rand.Intn(len(info.PlayerCards[myPos])))
	}

	return ai.getActionSmartyRandom()
}
