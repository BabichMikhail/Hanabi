package ai

import (
	"fmt"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI3 struct {
	BaseAI
	IsOriginal bool
	Depth      int
}

func NewAI3(baseAI *BaseAI) *AI3 {
	ai := new(AI3)
	ai.BaseAI = *baseAI
	ai.IsOriginal = true
	ai.Depth = 3
	return ai
}

func (ai *AI3) CheckInformation() {
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

func (ai *AI3) CheckPlayerInfo(playerInfo *game.PlayerGameInfo) {
	for pos, cards := range playerInfo.PlayerCards {
		for i := 0; i < len(cards); i++ {
			card := &cards[i]
			if !card.KnownColor || !card.KnownValue {
				fmt.Println(pos, i)
				panic("Bad player inforamtion")
			}
		}
	}
}

func (ai *AI3) CheckPreview(preview *game.ResultPreviewPlayerInformations) {
	for _, result := range preview.Results {
		ai.CheckPlayerInfo(result.Info)
	}
}

func (ai *AI3) FilterAndPreviewAvailablePlayerInformations(availablePlayerInformations []*game.AvailablePlayerGameInfo, step int, histAction *game.Action, idx1 int) []*game.AvailablePlayerGameInfo {
	nextPlayerInformation := ai.Informator.GetPlayerState(step + 1)
	nextAvailablePlayerInformations := map[string]*game.AvailablePlayerGameInfo{}
	fmt.Println(histAction)
	fmt.Println("Index:", idx1)
	for idx, information := range availablePlayerInformations {
		if idx1 == idx {
			fmt.Println("Attention:", idx)
		}

		previewResult, err := information.PlayerInfo.PreviewAction(histAction)
		if err != nil {
			panic("Bad preview. " + err.Error())
		}

		ai.CheckPreview(previewResult)
		for _, result := range previewResult.Results {
			resultPInfo := result.Info
			playerPos := histAction.PlayerPosition
			resultOK := true
			if histAction.ActionType == game.TypeActionInformationColor {
				cards := nextPlayerInformation.PlayerCards[playerPos]
				for j := 0; j < len(cards); j++ {
					if cards[j].KnownColor && cards[j].Color != resultPInfo.PlayerCards[playerPos][j].Color {
						resultOK = false
						break
					}
				}
				if !resultOK {
					continue
				}
				/*resultOK = false
				color := game.CardColor(histAction.Value)
				for j := 0; j < len(cards); j++ {
					_, ok1 := cards[j].ProbabilityColors[color]
					_, ok2 := resultPInfo.PlayerCards[playerPos][j].ProbabilityColors[color]
					resultOK = resultOK || !ok1 && ok2
				}*/
			} else if histAction.ActionType == game.TypeActionInformationValue {
				cards := nextPlayerInformation.PlayerCards[playerPos]
				for j := 0; j < len(cards); j++ {
					if cards[j].KnownValue && cards[j].Value != resultPInfo.PlayerCards[playerPos][j].Value {
						resultOK = false
						break
					}
				}
				if !resultOK {
					continue
				}
				/*resultOK = false
				value := game.CardValue(histAction.Value)
				for j := 0; j < len(cards); j++ {
					_, ok1 := cards[j].ProbabilityValues[value]
					_, ok2 := resultPInfo.PlayerCards[playerPos][j].ProbabilityValues[value]
					resultOK = resultOK || !ok1 && ok2
				}*/
			} else if histAction.ActionType == game.TypeActionDiscard {
				cards := nextPlayerInformation.UsedCards
				card1 := cards[len(cards)-1]
				card2 := result.Info.UsedCards[len(cards)-1]
				if card1.Color != card2.Color || card1.Value != card2.Value {
					continue
				}
			} else if histAction.ActionType == game.TypeActionPlaying {
				if nextPlayerInformation.BlueTokens != resultPInfo.BlueTokens {
					//fmt.Println("Cont1")
					continue
				}
				if nextPlayerInformation.RedTokens != resultPInfo.RedTokens {
					//fmt.Println("Cont2", nextPlayerInformation.RedTokens, resultPInfo.RedTokens, nextPlayerInformation.CurrentPosition)
					continue
				}

				cards := nextPlayerInformation.UsedCards
				if len(cards) > 0 {
					card1 := cards[len(cards)-1]
					card2 := result.Info.UsedCards[len(cards)-1]
					if card1.Color != card2.Color || card1.Value != card2.Value {
						//fmt.Println("Cont3")
						continue
					}
				}

				for color, card := range resultPInfo.TableCards {
					if nextPlayerInformation.TableCards[color].Value != card.Value {
						resultOK = false
						//fmt.Println("Cont4")
						break
					}
				}
			}

			if !resultOK {
				continue
			}

			hashKey := ai.Informator.PlayerInfoHash(result.Info)
			if APInfo, ok := nextAvailablePlayerInformations[hashKey]; ok {
				APInfo.Probability += result.Probability
			} else {
				newAPInfo := &game.AvailablePlayerGameInfo{
					PlayerInfo:  result.Info,
					Probability: result.Probability,
				}
				nextAvailablePlayerInformations[hashKey] = newAPInfo
			}

			if idx == idx1 {
				fmt.Println("etalon OK")
			}
		}

	}

	fmt.Println("Delta:", len(availablePlayerInformations), len(nextAvailablePlayerInformations), histAction.String())

	newAvailablePlayerInformations := []*game.AvailablePlayerGameInfo{}
	for _, APInfo := range nextAvailablePlayerInformations {
		newAvailablePlayerInformations = append(newAvailablePlayerInformations, APInfo)
	}
	return newAvailablePlayerInformations
}

func (ai *AI3) GetAction() *game.Action {
	ai.CheckInformation()
	info := &ai.PlayerInfo
	baseTime := 9 * len(info.PlayerCards)
	if info.Step != len(ai.History) {
		fmt.Println("History/Step:", info.Step, len(ai.History))
		panic("Bad History or Step")
	}
	baseAIType := Type_AIUsefulInformationV2 // Type_AI1
	myPos := info.CurrentPosition

	pdata := ai.Informator.GetCache()
	var data map[int]map[interface{}]interface{}
	if pdata == nil {
		data = map[int]map[interface{}]interface{}{}
		//ai.Informator.SetProbabilities(info)
	} else {
		data = pdata.(map[int]map[interface{}]interface{})
		if _, ok := data[myPos]; ok {
			stateHashValue := ai.Informator.PlayerInfoHash(info)
			myCache := data[myPos]
			if _, ok := myCache[stateHashValue]; ok {
				return myCache[stateHashValue].(*game.Action)
			}
		}
	}

	myCache := data[myPos]
	if myCache == nil {
		if info.Step > 10 && ai.IsOriginal {
			panic("MAGIC")
		}
		myCache = map[interface{}]interface{}{}
	}

	if ai.IsOriginal {
		fmt.Println("CurrentStep:", info.Step)
	}

	if info.Step < baseTime {
		copyInfo := info.Copy()
		hashValue := ai.Informator.PlayerInfoHash(info)
		ai.Informator.SetProbabilities(info)
		//fmt.Println("STep and HistLen", len(ai.History), info.Step)
		action := ai.Informator.GetAction(info.Copy(), baseAIType, ai.History)
		myCache[hashValue] = action.Copy()
		ai.Informator.SetProbabilities(copyInfo)
		myCache[info.Step] = copyInfo

		data[myPos] = myCache
		//fmt.Println(data)
		if ai.IsOriginal {
			fmt.Println("Save:", info.Step, myPos)
			ai.Informator.SetCache(data)
		}
		return action
	}

	fmt.Println("Line 230", info.Step, myPos, info.Step-len(info.PlayerCards))
	//if myCache[info.Step-len(info.PlayerCards)] == nil {
	//fmt.Println(myCache)
	//}
	/*prevStep = info.Step - len(info.PlayerCards)
	if prevStep < baseTime+1-len(info.PlayerCards) {
		prevStep = baseTime + 1 - len(info.PlayerCards)
	}*/
	myLastInfo := myCache[info.Step-len(info.PlayerCards)].(*game.PlayerGameInfo)

	//step := len(ai.History) - len(info.PlayerCards) + 1 + len(info.PlayerCards) - 2
	//playerInfo := ai.Informator.GetPlayerState(step)
	//ai.Informator.SetProbabilities(&playerInfo)
	fmt.Println("Generate information")
	availablePlayerInformations := myLastInfo.AvailablePlayerInformations()
	fmt.Println("Step:", myLastInfo.Step)
	//originalPlayerInfoIdx := ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformations, myLastInfo.Step)
	//fmt.Println(originalPlayerInfoIdx)
	if ai.IsOriginal {
		fmt.Println("len(availablePlayerInformation)1:", len(availablePlayerInformations))
	}

	availablePlayerInformations = ai.FilterAndPreviewAvailablePlayerInformations(availablePlayerInformations, myLastInfo.Step, &ai.History[myLastInfo.Step], 0)
	//originalPlayerInfoIdx = ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformations, myLastInfo.Step+1)
	//fmt.Println(originalPlayerInfoIdx)
	if ai.IsOriginal {
		fmt.Println("len(availablePlayerInformation)2:", len(availablePlayerInformations))
	}

	for i := myLastInfo.Step + 1; i < len(ai.History); i++ {
		probSum := 0.0
		histAction := &ai.History[i]

		newAvailablePlayerInformations := []*game.AvailablePlayerGameInfo{}
		for _, information := range availablePlayerInformations {
			playerInfo := information.PlayerInfo
			ai.CheckPlayerInfo(playerInfo)
			curPos := playerInfo.CurrentPosition
			playerCards := playerInfo.PlayerCards[curPos]
			playerInfo.PlayerCards[curPos] = playerInfo.PlayerCardsInfo[curPos]

			f := func() *game.Action {
				var action *game.Action
				defer func(action *game.Action) *game.Action {
					if r := recover(); r != nil {
						fmt.Println("Recover:", r)
						return nil
					}
					return action
				}(action)
				newPlayerInfo := playerInfo.Copy()
				//fmt.Println("AQWELJWQFIEW", newPlayerInfo.Step, len(ai.History[:i]), i)
				newAI := NewAI(*newPlayerInfo, ai.History[:i], Type_AI3, ai.Informator).(*AI3)
				//newAI.Depth = ai.Depth - 1
				newAI.IsOriginal = false
				action = newAI.GetAction()
				return action
			}
			action := f()

			if action == nil {
				//panic("ABC")
				//fmt.Println("Action is nil")
				continue
			}

			playerInfo.PlayerCards[curPos] = playerCards

			if histAction.Equal(action) {
				probSum += information.Probability
				newAvailablePlayerInformations = append(newAvailablePlayerInformations, information)
			}
		}

		availablePlayerInformations = newAvailablePlayerInformations
		fmt.Println("1234qwe")
		//originalPlayerInfoIdx = ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformations, i)

		//fmt.Println(originalPlayerInfoIdx)
		if ai.IsOriginal {
			fmt.Println("len(availablePlayerInformation)3:", len(availablePlayerInformations))
		}
		for j := 0; j < len(availablePlayerInformations); j++ {
			availablePlayerInformations[j].Probability /= probSum
		}

		if !ai.IsOriginal && len(availablePlayerInformations) == 0 {
			return nil
		}

		if ai.IsOriginal && len(availablePlayerInformations) == 0 {
			panic("I have no available information after filtering")
		}
		if ai.IsOriginal {
			fmt.Println("LENLEN1", len(availablePlayerInformations))
		}

		availablePlayerInformations = ai.FilterAndPreviewAvailablePlayerInformations(availablePlayerInformations, i, &ai.History[i], 0)
		//originalPlayerInfoIdx = ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformations, i)
		if ai.IsOriginal {
			fmt.Println("LENLEN2", len(availablePlayerInformations))
		}
	}

	resultAction := ai.Informator.GetAction(info, baseAIType, ai.History)
	data[myPos] = map[interface{}]interface{}{info.HashKey: resultAction.Copy()}
	if ai.IsOriginal {
		ai.Informator.SetCache(data)
	}
	return resultAction
}
