package ai

import (
	"math/rand"
	"time"

	"github.com/BabichMikhail/Hanabi/game"
)

type AIRandom struct {
	BaseAI
}

func NewAIRandom(baseAI *BaseAI) *AIRandom {
	ai := new(AIRandom)
	ai.BaseAI = *baseAI
	return ai
}

func (ai *AIRandom) GetAction() game.Action {
	return ai.getActionRandom()
}

func (ai *BaseAI) getActionRandom() game.Action {
	ai.setAvailableActions()
	rand.Seed(time.Now().UTC().UnixNano())
	return ai.Actions[rand.Intn(len(ai.Actions))].Action
}
