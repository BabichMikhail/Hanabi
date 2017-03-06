package ai

import (
	"math/rand"
	"time"

	"github.com/BabichMikhail/Hanabi/game"
)

func (ai *AI) getActionSmartyRandom() game.Action {
	ai.SetAvailableInfomation()
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
