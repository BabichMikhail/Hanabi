package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInfoAndMinMax struct {
	AIUsefulInfoAndMMEnd
}

func NewAIUsefulInfoAndMinMax(baseAI *BaseAI) *AIUsefulInfoAndMinMax {
	ai := new(AIUsefulInfoAndMinMax)
	ai.BaseAI = *baseAI
	ai.Depth = 4
	ai.AIUsefulInfoAndMMEnd.resultIsBetterThan = ai.resultIsBetterThan
	return ai
}

func (ai *AIUsefulInfoAndMinMax) resultIsBetterThan(result1, result2 *game.ResultPreviewPlayerInformations) bool {
	return result2 == nil || result1.Min > result2.Min ||
		result1.Min == result2.Min && result1.Med > result2.Med ||
		result1.Min == result2.Min && result1.Med == result2.Med && result1.Max > result2.Max
}
