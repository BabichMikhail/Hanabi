package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInfoAndMaxMax struct {
	AIUsefulInfoAndMMEnd
}

func NewAIUsefulInfoAndMaxMax(baseAI *BaseAI) *AIUsefulInfoAndMaxMax {
	ai := new(AIUsefulInfoAndMaxMax)
	ai.BaseAI = *baseAI
	ai.Depth = 3
	ai.AIUsefulInfoAndMMEnd.resultIsBetterThan = ai.resultIsBetterThan
	return ai
}

func (ai *AIUsefulInfoAndMaxMax) resultIsBetterThan(result1, result2 *game.ResultPreviewPlayerInformations) bool {
	return result2 == nil || result1.Max > result2.Max ||
		result1.Max == result2.Max && result1.Med > result2.Med ||
		result1.Max == result2.Max && result1.Med == result2.Med && result1.Min > result2.Min
}
