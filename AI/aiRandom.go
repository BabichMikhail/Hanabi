package ai

import (
	"math/rand"
	"time"
)

func (ai *AI) getRandomAction() *Action {
	rand.Seed(time.Now().UTC().UnixNano())
	return ai.Actions[rand.Intn(len(ai.Actions))]
}
