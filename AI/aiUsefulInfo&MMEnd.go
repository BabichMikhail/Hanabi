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

func (ai *AIUsefulInfoAndMMEnd) GetNewHistory(newAction game.Action) []game.Action {
	newActions := make([]game.Action, len(ai.History)+1, len(ai.History)+1)
	copy(newActions, ai.History)
	newActions[len(newActions)-1] = newAction
	return newActions
}

func (ai *AIUsefulInfoAndMMEnd) getBestResultWithDepth() *game.ResultPreviewPlayerInformations {
	info := &ai.PlayerInfo
	pos := info.CurrentPostion
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

	for i := 0; i < len(info.PlayerCards); i++ {
		if i == pos {
			continue
		}
		cardColors := map[game.CardColor]struct{}{}
		cardValues := map[game.CardValue]struct{}{}
		for k := 0; k < len(info.PlayerCards[i]); k++ {
			cardColors[info.PlayerCards[i][k].Color] = struct{}{}
			cardValues[info.PlayerCards[i][k].Value] = struct{}{}
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
	pos := info.CurrentPostion
	var bestResult *game.ResultPreviewPlayerInformations

	for i := 0; i < len(info.PlayerCards[pos]); i++ {
		resultDiscard, err := info.PreviewActionDiscard(i)
		if err == nil && ai.resultIsBetterThan(resultDiscard, bestResult) {
			bestResult = resultDiscard
		}

		resultPlaying, err := info.PreviewActionPlaying(i)
		if err == nil && ai.resultIsBetterThan(resultPlaying, bestResult) {
			bestResult = resultPlaying
		}
	}

	for i := 0; i < len(info.PlayerCards); i++ {
		if i == pos {
			continue
		}
		cardColors := map[game.CardColor]struct{}{}
		cardValues := map[game.CardValue]struct{}{}
		for k := 0; k < len(info.PlayerCards[i]); k++ {
			cardColors[info.PlayerCards[i][k].Color] = struct{}{}
			cardValues[info.PlayerCards[i][k].Value] = struct{}{}
		}

		for cardColor, _ := range cardColors {
			resultInfo, err := info.PreviewActionInformationColor(i, cardColor)
			if err == nil && ai.resultIsBetterThan(resultInfo, bestResult) {
				bestResult = resultInfo
			}
		}

		for cardValue, _ := range cardValues {
			resultInfo, err := info.PreviewActionInformationValue(i, cardValue)
			if err == nil && ai.resultIsBetterThan(resultInfo, bestResult) {
				bestResult = resultInfo
			}
		}
	}

	return bestResult
}

func (ai *AIUsefulInfoAndMMEnd) GetBestResult() *game.ResultPreviewPlayerInformations {
	info := &ai.PlayerInfo
	if info.IsGameOver() {
		return nil
	}

	if ai.Depth > 1 && info.Step <= info.MaxStep {
		return ai.getBestResultWithDepth()
	}

	return ai.getBestResultWithoutDepth()
}

func (ai *AIUsefulInfoAndMMEnd) GetAction() game.Action {
	ai.setAvailableInfomation()
	info := &ai.PlayerInfo
	if info.DeckSize > 0 {
		ai := NewAI(*info, ai.History, Type_AIUsefulInformationV3, ai.Informator)
		return ai.GetAction()
	}

	return ai.GetBestResult().Action
}
