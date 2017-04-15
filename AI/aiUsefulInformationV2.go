package ai

import (
	"math"
	"math/rand"
	"sort"

	"github.com/BabichMikhail/Hanabi/game"
)

type UsefulAction struct {
	Action     *game.Action
	Usefulness float64
}

type Actions []UsefulAction

func (a Actions) Len() int {
	return len(a)
}

func (a Actions) Less(i, j int) bool {
	return a[i].Usefulness >= a[j].Usefulness
}

func (a Actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type AIUsefulInformationV2 struct {
	BaseAI

	CoefPlayByValue float64
	CoefPlayByColor float64
	CoefInfoValue   float64
	CoefInfoColor   float64
}

func NewAIUsefulInformationV2(baseAI *BaseAI) *AIUsefulInformationV2 {
	ai := new(AIUsefulInformationV2)
	ai.BaseAI = *baseAI
	ai.CoefPlayByValue = 1.1
	ai.CoefPlayByColor = -0.9
	ai.CoefInfoValue = 1.05
	ai.CoefInfoColor = 1.0
	return ai
}

func (ai *AIUsefulInformationV2) SetCoefs(kPlayByValue, kPlayByColor, kInfoValue, kInfoColor float64) {
	ai.CoefPlayByValue = kPlayByValue
	ai.CoefPlayByColor = kPlayByColor
	ai.CoefInfoValue = kInfoValue
	ai.CoefInfoColor = kInfoColor
}

func (ai *AIUsefulInformationV2) GetAction() *game.Action {
	ai.setAvailableInformation()
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
								Usefulness: ai.CoefInfoValue * (1.0 - float64(i)/float64(len(info.PlayerCards)-1)),
							}
							usefulActions = append(usefulActions, action)
						}

						if !cardInfo.KnownColor {
							action := UsefulAction{
								Action:     game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color)),
								Usefulness: ai.CoefInfoColor * (1.0 - float64(i)/float64(len(info.PlayerCards)-1)),
							}
							usefulActions = append(usefulActions, action)
						}
					}
				}
			}
		}
	}

	if len(usefulActions) > 0 {
		sort.Sort(usefulActions)
		return usefulActions[0].Action
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
