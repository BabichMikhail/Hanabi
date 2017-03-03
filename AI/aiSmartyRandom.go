package ai

import (
	"math/rand"
	"time"

	"github.com/BabichMikhail/Hanabi/game"
)

func (ai *AI) getSmartyRandomActionIdx() int {
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

	rand.Seed(time.Now().UTC().UnixNano())
	if len(usefullActions) > 0 {
		return rand.Intn(len(usefullActions))
	} else {
		return rand.Intn(len(ai.Actions))
	}
}
