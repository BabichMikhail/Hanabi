package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInfoAndMedMax struct {
	AIUsefulInfoAndMMEnd
}

func NewAIUsefulInfoAndMedMax(baseAI *BaseAI) *AIUsefulInfoAndMedMax {
	ai := new(AIUsefulInfoAndMedMax)
	ai.BaseAI = *baseAI
	ai.Depth = Min(4, ai.PlayerInfo.PlayerCount)
	ai.AIUsefulInfoAndMMEnd.resultIsBetterThan = ai.resultIsBetterThan
	return ai
}

func (ai *AIUsefulInfoAndMedMax) resultIsBetterThan(result1, result2 *game.ResultPreviewPlayerInformations) bool {
	return result2 == nil || result1.Med > result2.Med ||
		result1.Med == result2.Med && result1.Max > result2.Max ||
		result1.Med == result2.Med && result1.Max == result2.Max && result1.Min > result2.Min
}
