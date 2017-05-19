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
	ai.Depth = 1
	ai.AIUsefulInfoAndMMEnd.resultIsBetterThan = ai.resultIsBetterThan
	return ai
}

func (ai *AIUsefulInfoAndMinMax) resultIsBetterThan(result1, result2 *game.ResultPreviewPlayerInformations) bool {
	if result2 == nil {
		return true
	}

	res1 := 0.0
	for i := 0; i < len(result1.Results); i++ {
		result := result1.Results[i]
		res1 += result.Probability * ai.Informator.GetQualitativeAssessmentOfState(result.Info)
	}

	res2 := 0.0
	for i := 0; i < len(result2.Results); i++ {
		result := result2.Results[i]
		res2 += result.Probability * ai.Informator.GetQualitativeAssessmentOfState(result.Info)
	}
	return res1 > res2
	/*return result2 == nil || result1.Min > result2.Min ||
	result1.Min == result2.Min && result1.Med > result2.Med ||
	result1.Min == result2.Min && result1.Med == result2.Med && result1.Max > result2.Max*/
}
