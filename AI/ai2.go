package ai

import (
	"fmt"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI2 struct {
	BaseAI
	IsOriginal bool
	Depth      int
}

func NewAI2(baseAI *BaseAI) *AI2 {
	ai := new(AI2)
	ai.BaseAI = *baseAI
	ai.IsOriginal = true
	ai.Depth = 3
	return ai
}

func (ai *AI2) CheckInformation() {
	info := &ai.PlayerInfo
	myRealPos := info.Position
	if myRealPos == info.CurrentPosition {
		return
	}
	for i := 0; i < len(info.PlayerCards[myRealPos]); i++ {
		card := &info.PlayerCards[myRealPos][i]
		if !card.KnownColor || !card.KnownValue {
			panic("Bad information for my position")
		}
	}
}

func (ai *AI2) GetAction() *game.Action {
	ai.CheckInformation()
	info := &ai.PlayerInfo
	if info.Step != len(ai.History) {
		panic("Bad History or Step")
	}
	baseAIType := Type_AI1
	myPos := info.CurrentPosition

	pdata := ai.Informator.GetCache()
	var data map[int]map[interface{}]interface{}
	if pdata == nil {
		data = map[int]map[interface{}]interface{}{}
		ai.Informator.SetProbabilities(info)
	} else {
		data = pdata.(map[int]map[interface{}]interface{})
		if _, ok := data[myPos]; ok {
			stateHashValue := ai.Informator.PlayerInfoHash(info)
			myData := data[myPos]
			if _, ok := myData[stateHashValue]; ok {
				return myData[stateHashValue].(*game.Action)
			}
		}
	}

	if info.DeckSize > 0 || ai.Depth == 0 {
		hashValue := ai.Informator.PlayerInfoHash(info)
		action := ai.Informator.GetAction(info.Copy(), baseAIType, ai.History)
		data[myPos] = map[interface{}]interface{}{hashValue: action.Copy()}
		ai.Informator.SetCache(data)
		return action
	}

	step := len(ai.History) - len(info.PlayerCards) + 1 + len(info.PlayerCards) - 2
	pinfo := ai.Informator.GetPlayerState(step)
	currentPlayerInfo := &pinfo
	for i := step; i <= len(ai.History)-1; i++ {
		availablePlayerInformation := currentPlayerInfo.AvailablePlayerInformations()
		var originalPlayerInfoIdx int
		g := func() {
			defer func() {
				if r := recover(); r != nil && ai.IsOriginal {
					panic(fmt.Sprint("Recover:", r))
				}
			}()
			/* Use this function for debug ONLY! */
			originalPlayerInfoIdx = ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformation, i)
		}
		g()

		probSum := 0.0
		histAction := ai.History[i]
		newAvailablePlayerInformation := []*game.AvailablePlayerGameInfo{}
		for idx, information := range availablePlayerInformation {
			playerInfo := information.PlayerInfo
			curPos := playerInfo.CurrentPosition
			playerCards := playerInfo.PlayerCards[curPos]
			playerInfo.PlayerCards[curPos] = playerInfo.PlayerCardsInfo[curPos]

			f := func() *game.Action {
				var action *game.Action
				defer func(action *game.Action) *game.Action {
					if r := recover(); r != nil {
						return nil
					}
					return action
				}(action)
				newPlayerInfo := playerInfo.Copy()
				newAI := NewAI(*newPlayerInfo, ai.History[:i], Type_AI2, ai.Informator).(*AI2)
				newAI.Depth = ai.Depth - 1
				newAI.IsOriginal = false
				action = newAI.GetAction()
				return action
			}
			action := f()

			if action == nil && originalPlayerInfoIdx == idx {
				panic("Action is nil on original state")
			}

			if action == nil {
				fmt.Println("Action is nil")
				continue
			}

			playerInfo.PlayerCards[curPos] = playerCards

			if histAction.Equal(action) {
				probSum += information.Probability
				newAvailablePlayerInformation = append(newAvailablePlayerInformation, information)
			}

			if !histAction.Equal(action) && originalPlayerInfoIdx == idx {
				fmt.Println(ai.Informator.PlayerInfoHash(playerInfo))
			}
		}

		for j := 0; j < len(newAvailablePlayerInformation); j++ {
			newAvailablePlayerInformation[j].Probability /= probSum
		}

		if !ai.IsOriginal && len(newAvailablePlayerInformation) == 0 {
			return nil
		}

		if ai.IsOriginal && len(newAvailablePlayerInformation) == 0 {
			panic("I have no available information after filtering")
		}

		nextPlayerInformation := ai.Informator.GetPlayerState(i + 1)
		currentPlayerInfo = &nextPlayerInformation
	}

	resultAction := ai.Informator.GetAction(info, baseAIType, ai.History)
	data[myPos] = map[interface{}]interface{}{info.HashKey: resultAction.Copy()}
	ai.Informator.SetCache(data)
	return resultAction
}
