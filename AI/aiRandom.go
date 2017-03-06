package ai

import (
	"math/rand"
	"time"

	"github.com/BabichMikhail/Hanabi/game"
)

func (ai *AI) getActionRandom() game.Action {
	rand.Seed(time.Now().UTC().UnixNano())
	return ai.Actions[rand.Intn(len(ai.Actions))].Action
}
