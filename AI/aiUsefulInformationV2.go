package ai

import (
	"math/rand"
	"sort"

	"github.com/BabichMikhail/Hanabi/game"
)

type UsefulAction struct {
	Action     game.Action
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
}

func NewAIUsefulInformationV2(baseAI *BaseAI) *AIUsefulInformationV2 {
	ai := new(AIUsefulInformationV2)
	ai.BaseAI = *baseAI
	return ai
}

func (ai *AIUsefulInformationV2) GetAction() game.Action {
	ai.setAvailableInfomation()
	info := &ai.PlayerInfo
	myPos := info.Position
	k_PlayByValue := 1.0
	k_PlayByColor := 1.0
	k_InfoColor := 1.0
	k_InfoValue := 1.0

	usefulActions := Actions{}

	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[myPos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, myPos, idx)
			}
		}
	}

	subHistory := ai.History[Max(len(ai.History)-len(info.PlayerCards)-1, 0):]
	historyLength := float64(len(subHistory))
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
					usefulActions = append(usefulActions, UsefulAction{game.NewAction(game.TypeActionPlaying, myPos, idx), k_PlayByValue * float64(i+1) / historyLength / count})
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
					usefulActions = append(usefulActions, UsefulAction{game.NewAction(game.TypeActionPlaying, myPos, idx), k_PlayByColor * float64(i+1) / historyLength / count})
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
							usefulActions = append(usefulActions, UsefulAction{game.NewAction(game.TypeActionInformationValue, nextPos, int(card.Value)), k_InfoValue * (1.0 - float64((i+1)/len(info.PlayerCards)))})
						}
						if !cardInfo.KnownColor {
							usefulActions = append(usefulActions, UsefulAction{game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color)), k_InfoColor * (1.0 - float64((i+1)/len(info.PlayerCards)))})
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
