package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type FCompareResults func(result1, result2 *game.ResultPreviewPlayerInformations) bool

type AIUsefulInfoAndMMEnd struct {
	BaseAI
	Depth              int
	resultIsBetterThan FCompareResults
}

type AIUsefulAndMMEndInterface interface {
	SetDepth(depth int)
	GetBestResult() *game.ResultPreviewPlayerInformations
}

func (ai *AIUsefulInfoAndMMEnd) SetDepth(depth int) {
	ai.Depth = depth
}

func (ai *AIUsefulInfoAndMMEnd) GetNewHistory(newAction *game.Action) []game.Action {
	newActions := make([]game.Action, len(ai.History)+1, len(ai.History)+1)
	copy(newActions, ai.History)
	newActions[len(newActions)-1] = *newAction
	return newActions
}

func (ai *BaseAI) checkCardUsefulByValues(card game.Card, f func(tableValue, cardValue game.CardValue) bool) bool {
	info := &ai.PlayerInfo
	if card.KnownColor && card.KnownValue {
		return f(info.TableCards[card.Color].Value, card.Value)
	}

	for colorValue, _ := range card.ProbabilityCard {
		color, value := game.ColorValueByHashColorValue(colorValue)
		if f(info.TableCards[color].Value, value) {
			return true
		}
	}
	return false
}

func (ai *BaseAI) isCardPlayable(card game.Card) bool {
	f := func(tableValue, cardValue game.CardValue) bool {
		return tableValue+1 == cardValue
	}
	return ai.checkCardUsefulByValues(card, f)
}

func (ai *BaseAI) isCardMayBeUsefull(card game.Card) bool {
	f := func(tableValue, cardValue game.CardValue) bool {
		return tableValue < cardValue
	}
	return ai.checkCardUsefulByValues(card, f)
}

func (ai *AIUsefulInfoAndMMEnd) getBestResultWithDepth() *game.ResultPreviewPlayerInformations {
	info := &ai.PlayerInfo
	pos := info.CurrentPosition
	var bestResult *game.ResultPreviewPlayerInformations

	for i := 0; i < len(info.PlayerCards[pos]); i++ {
		resultDiscard, err := info.PreviewActionDiscard(i)
		if err != nil {
			continue
		}
		newHistory := ai.GetNewHistory(resultDiscard.Action)
		newMax := -1
		newMin := 26
		newMed := 0.0
		for j := 0; j < len(resultDiscard.Results); j++ {
			newAI := NewAI(*resultDiscard.Results[j].Info, newHistory, ai.Type, ai.Informator).(AIUsefulAndMMEndInterface)
			newAI.SetDepth(ai.Depth - 1)
			newResult := newAI.GetBestResult()
			if newResult == nil {
				break
			}
			if newResult.Max > newMax {
				newMax = newResult.Max
			}
			if newResult.Min < newMin {
				newMin = newResult.Min
			}
			newMed += newResult.Med * resultDiscard.Results[j].Probability
		}

		if newMax != -1 {
			resultDiscard.Max = newMax
			resultDiscard.Min = newMin
			resultDiscard.Med = newMed
		}

		if ai.resultIsBetterThan(resultDiscard, bestResult) {
			bestResult = resultDiscard
		}
	}

	for i := 0; i < len(info.PlayerCards[pos]); i++ {
		resultPlaying, err := info.PreviewActionPlaying(i)
		if err != nil {
			continue
		}
		newHistory := ai.GetNewHistory(resultPlaying.Action)
		newMax := -1
		newMin := 26
		newMed := 0.0
		for j := 0; j < len(resultPlaying.Results); j++ {
			newAI := NewAI(*resultPlaying.Results[j].Info, newHistory, ai.Type, ai.Informator).(AIUsefulAndMMEndInterface)
			newAI.SetDepth(ai.Depth - 1)
			newResult := newAI.GetBestResult()
			if newResult == nil {
				break
			}
			if newResult.Max > newMax {
				newMax = newResult.Max
			}
			if newResult.Min < newMin {
				newMin = newResult.Min
			}
			newMed += newResult.Med * resultPlaying.Results[j].Probability
		}

		if newMax != -1 {
			resultPlaying.Max = newMax
			resultPlaying.Min = newMin
			resultPlaying.Med = newMed
		}

		if ai.resultIsBetterThan(resultPlaying, bestResult) {
			bestResult = resultPlaying
		}
	}

	isActionInfo := false
	for i := 0; i < len(info.PlayerCards); i++ {
		if i == pos || isActionInfo && info.MaxStep > info.Step && (i-info.CurrentPosition+len(info.PlayerCards))%len(info.PlayerCards) > info.MaxStep-info.Step {
			continue
		}

		cardColors := map[game.CardColor]struct{}{}
		cardValues := map[game.CardValue]struct{}{}
		for k := 0; k < len(info.PlayerCards[i]); k++ {
			if !ai.isCardMayBeUsefull(info.PlayerCards[i][k]) {
				continue
			}
			isActionInfo = true
			cardColors[info.PlayerCards[i][k].Color] = struct{}{}
			cardValues[info.PlayerCards[i][k].Value] = struct{}{}
		}

		if i == len(info.PlayerCards)-1 && !isActionInfo {
			k := len(info.PlayerCards[i]) - 1
			cardColors[info.PlayerCards[i][k].Color] = struct{}{}
		}

		for cardColor, _ := range cardColors {
			resultInfo, err := info.PreviewActionInformationColor(i, cardColor)
			if err != nil {
				continue
			}
			newHistory := ai.GetNewHistory(resultInfo.Action)
			newMax := -1
			newMin := 26
			newMed := 0.0
			for j := 0; j < len(resultInfo.Results); j++ {
				newAI := NewAI(*resultInfo.Results[j].Info, newHistory, ai.Type, ai.Informator).(AIUsefulAndMMEndInterface)
				newAI.SetDepth(ai.Depth - 1)
				newResult := newAI.GetBestResult()
				if newResult == nil {
					break
				}
				if newResult.Max > newMax {
					newMax = newResult.Max
				}
				if newResult.Min < newMin {
					newMin = newResult.Min
				}
				newMed += newResult.Med * resultInfo.Results[j].Probability
			}

			if newMax != -1 {
				resultInfo.Max = newMax
				resultInfo.Min = newMin
				resultInfo.Med = newMed
			}

			if ai.resultIsBetterThan(resultInfo, bestResult) {
				bestResult = resultInfo
			}
		}

		for cardValue, _ := range cardValues {
			resultInfo, err := info.PreviewActionInformationValue(i, cardValue)
			if err != nil {
				continue
			}
			newHistory := ai.GetNewHistory(resultInfo.Action)
			newMax := -1
			newMin := 26
			newMed := 0.0
			for j := 0; j < len(resultInfo.Results); j++ {
				newAI := NewAI(*resultInfo.Results[j].Info, newHistory, ai.Type, ai.Informator).(AIUsefulAndMMEndInterface)
				newAI.SetDepth(ai.Depth - 1)
				newResult := newAI.GetBestResult()
				if newResult == nil {
					break
				}
				if newResult.Max > newMax {
					newMax = newResult.Max
				}
				if newResult.Min < newMin {
					newMin = newResult.Min
				}
				newMed += newResult.Med * resultInfo.Results[j].Probability
			}

			if newMax != -1 {
				resultInfo.Max = newMax
				resultInfo.Min = newMin
				resultInfo.Med = newMed
			}

			if ai.resultIsBetterThan(resultInfo, bestResult) {
				bestResult = resultInfo
			}
		}
	}
	return bestResult
}

func (ai *AIUsefulInfoAndMMEnd) getBestResultWithoutDepth() *game.ResultPreviewPlayerInformations {
	info := &ai.PlayerInfo
	var bestResult *game.ResultPreviewPlayerInformations

	myPos := info.CurrentPosition
	for i := 0; i < len(info.PlayerCards[myPos]); i++ {
		newResult, err := info.PreviewActionPlaying(i)
		if err == nil && ai.resultIsBetterThan(newResult, bestResult) {
			bestResult = newResult
		}
	}

	for i := 0; i < len(info.PlayerCards[myPos]); i++ {
		newResult, err := info.PreviewActionDiscard(i)
		if err == nil && ai.resultIsBetterThan(newResult, bestResult) {
			bestResult = newResult
		}
	}

	for i := 0; i < len(info.PlayerCards); i++ {
		if i == myPos {
			continue
		}

		for j := 0; j < len(info.PlayerCards[i]); j++ {
			newResult, err := info.PreviewActionInformationColor(i, info.PlayerCards[i][j].Color)
			if err == nil && ai.resultIsBetterThan(newResult, bestResult) {
				bestResult = newResult
			}
		}

		for j := 0; j < len(info.PlayerCards[i]); j++ {
			newResult, err := info.PreviewActionInformationValue(i, info.PlayerCards[i][j].Value)
			if err == nil && ai.resultIsBetterThan(newResult, bestResult) {
				bestResult = newResult
			}
		}
	}

	/*newAI := NewAI(*info, ai.History, Type_AIUsefulInformationV3, ai.Informator).(*AIUsefulInfoV3AndParts)
	newAI.setAvailableInformation()
	action := newAI.GetAction()
	switch action.ActionType {
	case game.TypeActionDiscard:
		resultDiscard, err := info.PreviewActionDiscard(action.Value)
		if err == nil && ai.resultIsBetterThan(resultDiscard, bestResult) {
			bestResult = resultDiscard
		}
	case game.TypeActionPlaying:
		resultPlaying, err := info.PreviewActionPlaying(action.Value)
		if err == nil && ai.resultIsBetterThan(resultPlaying, bestResult) {
			bestResult = resultPlaying
		}
	case game.TypeActionInformationColor:
		resultInfoColor, err := info.PreviewActionInformationColor(action.PlayerPosition, game.CardColor(action.Value))
		if err == nil && ai.resultIsBetterThan(resultInfoColor, bestResult) {
			bestResult = resultInfoColor
		}
	case game.TypeActionInformationValue:
		resultInfoValue, err := info.PreviewActionInformationValue(action.PlayerPosition, game.CardValue(action.Value))
		if err == nil && ai.resultIsBetterThan(resultInfoValue, bestResult) {
			bestResult = resultInfoValue
		}
	}*/
	return bestResult
}

func (ai *AIUsefulInfoAndMMEnd) GetBestResult() *game.ResultPreviewPlayerInformations {
	info := &ai.PlayerInfo
	if info.IsGameOver() {
		return nil
	}

	/*if ai.Depth > 1 && info.Step <= info.MaxStep {
		return ai.getBestResultWithDepth()
	}*/
	return ai.getBestResultWithoutDepth()
}

func (ai *AIUsefulInfoAndMMEnd) GetAction() *game.Action {
	ai.setAvailableInformation()
	//info := &ai.PlayerInfo
	return ai.GetBestResult().Action
	/*if info.DeckSize > 0 {
		ai := NewAI(*info, ai.History, Type_AIUsefulInformationV3, ai.Informator)
		return ai.GetAction()
	}*/

	//return ai.GetBestResult().Action
}
