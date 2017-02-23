package ai

import (
	"math/rand"
	"time"
)

func (ai *AI) getRandomActionIdx() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(len(ai.Actions))
}
