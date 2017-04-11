package ai

import (
	"math/rand"
	"time"

	"github.com/BabichMikhail/Hanabi/game"
)

type AISmartyRandom struct {
	BaseAI
}

func NewAISmartyRandom(baseAI *BaseAI) *AIRandom {
	ai := new(AIRandom)
	ai.BaseAI = *baseAI
	return ai
}

func (ai *AISmartyRandom) GetAction() game.Action {
	return ai.getActionSmartyRandom()
}

func (ai *BaseAI) getActionSmartyRandom() game.Action {
	ai.setAvailableActions()
	ai.setAvailableInfomation()
	info := &ai.PlayerInfo
	var usefullActions []*Action
	for _, action := range ai.InfoValueActions {
		for _, card := range info.PlayerCards[action.PlayerPosition] {
			if game.CardValue(action.Value) == card.Value && !card.KnownValue {
				usefullActions = append(usefullActions, action)
			}
		}
	}

	for _, action := range ai.InfoColorAcions {
		for _, card := range info.PlayerCards[action.PlayerPosition] {
			if game.CardColor(action.Value) == card.Color && !card.KnownColor {
				usefullActions = append(usefullActions, action)
			}
		}
	}

	for _, action := range ai.DiscardActions {
		usefullActions = append(usefullActions, action)
	}

	for _, action := range ai.PlayActions {
		usefullActions = append(usefullActions, action)
	}

	if len(usefullActions) > 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		return usefullActions[rand.Intn(len(usefullActions))].Action
	} else {
		return ai.getActionRandom()
	}
}
