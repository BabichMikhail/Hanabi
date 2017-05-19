package ai

import (
	"fmt"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI5 struct {
	BaseAI
	IsOriginal bool
	//Depth      int
}

func NewAI5(baseAI *BaseAI) *AI5 {
	ai := new(AI5)
	ai.BaseAI = *baseAI
	ai.IsOriginal = true
	//ai.Depth = 3
	return ai
}

func (ai *AI5) CheckInformation() {
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

func (ai *AI5) CheckPlayerInfo(playerInfo *game.PlayerGameInfo) {
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

func (ai *AI5) CheckPreview(preview *game.ResultPreviewPlayerInformations) {
	for _, result := range preview.Results {
		ai.CheckPlayerInfo(result.Info)
	}
}

func (ai *AI5) FilterAndPreviewAvailablePlayerInformations(availablePlayerInformations []*game.AvailablePlayerGameInfo, step int, histAction *game.Action) []*game.AvailablePlayerGameInfo {
	nextPlayerInformation := ai.Informator.GetPlayerState(step + 1)
	nextAvailablePlayerInformations := map[string]*game.AvailablePlayerGameInfo{}
	fmt.Println(nextPlayerInformation)
	fmt.Println(histAction)
	for _, information := range availablePlayerInformations {
		fmt.Println(information.PlayerInfo)
		previewResult, err := information.PlayerInfo.PreviewAction(histAction)
		if err != nil {
			panic("Bad preview. " + err.Error())
		}

		fmt.Println("RESULTS:", len(availablePlayerInformations))
		//ai.CheckPreview(previewResult)
		for _, result := range previewResult.Results {
			resultPInfo := result.Info
			resultPInfo.SetProbabilities(false, false)
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
			} else if histAction.ActionType == game.TypeActionInformationValue {
				cards := nextPlayerInformation.PlayerCards[playerPos]
				for j := 0; j < len(cards); j++ {
					if cards[j].KnownValue && cards[j].Value != resultPInfo.PlayerCards[playerPos][j].Value {
						resultOK = false
						break
					}
				}
			} else if histAction.ActionType == game.TypeActionDiscard {
				cards := nextPlayerInformation.UsedCards
				card1 := cards[len(cards)-1]
				card2 := resultPInfo.UsedCards[len(cards)-1]
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
					card2 := resultPInfo.UsedCards[len(cards)-1]
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

			hashKey := ai.Informator.PlayerInfoHash(resultPInfo)
			fmt.Println(hashKey)
			if APInfo, ok := nextAvailablePlayerInformations[hashKey]; ok {
				APInfo.Probability += result.Probability
			} else {
				newAPInfo := &game.AvailablePlayerGameInfo{
					PlayerInfo:  resultPInfo,
					Probability: result.Probability,
				}
				nextAvailablePlayerInformations[hashKey] = newAPInfo
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

func (ai *AI5) GetAction() *game.Action {
	ai.CheckInformation()
	info := &ai.PlayerInfo
	baseTime := 0 * len(info.PlayerCards)
	if info.Step != len(ai.History) {
		fmt.Println("History/Step:", info.Step, len(ai.History))
		panic("Bad History or Step")
	}
	baseAIType := Type_AI1
	myPos := info.CurrentPosition

	pdata := ai.Informator.GetCache()
	var data map[int]map[interface{}]interface{}
	if pdata == nil {
		data = map[int]map[interface{}]interface{}{}
	} else {
		data = pdata.(map[int]map[interface{}]interface{})
		/*if _, ok := data[myPos]; ok {
			stateHashValue := ai.Informator.PlayerInfoHash(info)
			myCache := data[myPos]
			if _, ok := myCache[stateHashValue]; ok {
				return myCache[stateHashValue].(*game.Action)
			}
		}*/
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

	if info.Step == 0 || info.Step < baseTime {
		copyInfo := info.Copy()
		//hashValue := ai.Informator.PlayerInfoHash(info)
		ai.Informator.SetProbabilities(info)
		//fmt.Println("STep and HistLen", len(ai.History), info.Step)
		action := ai.Informator.GetAction(info.Copy(), baseAIType, ai.History)
		//myCache[hashValue] = action.Copy()
		ai.Informator.SetProbabilities(copyInfo)
		myCache[info.Step] = copyInfo
		data[myPos] = myCache
		//fmt.Println(data)
		//if ai.IsOriginal {
		//	fmt.Println("Save:", info.Step, myPos)
		ai.Informator.SetCache(data)
		//}
		return action
	}

	//if myCache[info.Step-len(info.PlayerCards)] == nil {
	//fmt.Println(myCache)
	//}
	/*prevStep = info.Step - len(info.PlayerCards)
	if prevStep < baseTime+1-len(info.PlayerCards) {
		prevStep = baseTime + 1 - len(info.PlayerCards)
	}*/
	prevState := myCache[info.Step-len(info.PlayerCards)]
	if prevState == nil {
		if info.Step >= len(info.PlayerCards) {
			panic("MAGIC")
		}
		if info.Step == 0 {
			prevState = info
		} else {
			firstPlayerInfo := ai.Informator.GetPlayerState(0)
			prevState = &firstPlayerInfo
		}
	}
	myLastInfo := prevState.(*game.PlayerGameInfo)
	ai.Informator.SetProbabilities(myLastInfo)

	//step := len(ai.History) - len(info.PlayerCards) + 1 + len(info.PlayerCards) - 2
	//playerInfo := ai.Informator.GetPlayerState(step)
	//ai.Informator.SetProbabilities(&playerInfo)
	fmt.Println("Generate information")
	//availablePlayerInformations := myLastInfo.AvailablePlayerInformations()
	availablePlayerInformations := game.AvailablePlayerGameInfos{
		&game.AvailablePlayerGameInfo{
			PlayerInfo:  myLastInfo,
			Probability: 1.0,
		},
	}
	fmt.Println("Step:", myLastInfo.Step)
	//originalPlayerInfoIdx := ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformations, myLastInfo.Step)
	//fmt.Println(originalPlayerInfoIdx)

	for i := myLastInfo.Step; i < len(ai.History); i++ {
		probSum := 0.0

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

		availablePlayerInformations = ai.FilterAndPreviewAvailablePlayerInformations(availablePlayerInformations, i, &ai.History[i])
		//originalPlayerInfoIdx = ai.Informator.CheckAvailablePlayerInformation(availablePlayerInformations, i)
		if ai.IsOriginal {
			fmt.Println("LENLEN2", len(availablePlayerInformations))
		}
	}

	if len(availablePlayerInformations) == 0 {
		panic("Len = 0")
	}

	if availablePlayerInformations[0].PlayerInfo.Step != info.Step {
		panic("Bad step")
	}

	/*for _, cards := range info.PlayerCards {
		for i := 0; i < len(cards); i++ {
			card := &cards[i]
			if card.KnownColor && card.KnownValue {
				card.ProbabilityColors = map[game.CardColor]float64{card.Color: 1.0}
				card.ProbabilityValues = map[game.CardValue]float64{card.Value: 1.0}
				card.ProbabilityCard = map[game.HashValue]float64{
					game.HashColorValue(card.Color, card.Value): 1.0,
				}
			}
		}
	}*/

	availablePlayerInformations.UpdatePlayerInformation(info)
	myCache[info.Step] = info.Copy()
	//resultAction := ai.Informator.GetAction(info, baseAIType, ai.History)
	newAI := NewAI(*info, ai.History, baseAIType, ai.Informator).(*AI1)
	resultAction := newAI.GetAction()

	/*data[myPos] = map[interface{}]interface{}{info.HashKey: resultAction.Copy()}
	if ai.IsOriginal {
		ai.Informator.SetCache(data)
	}*/
	data[myPos] = myCache
	ai.Informator.SetCache(data)
	return resultAction
}
