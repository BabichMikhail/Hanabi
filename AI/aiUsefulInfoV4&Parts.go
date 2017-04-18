package ai

import (
	"math"

	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInfoV4AndPartsCoefs struct {
	CoefPlayByInfoA float64
	CoefPlayByInfoB float64
	CoefInfoA       float64
	CoefInfoB       float64
	CoefDiscardA    float64
	CoefDiscardB    float64
	CoefPlayA       float64
	CoefPlayB       float64
}

type AIUsefulInfoV4AndParts struct {
	BaseAI

	Coefs []AIUsefulInfoV4AndPartsCoefs
}

func NewAIUsefulInfoV4AndParts(baseAI *BaseAI) *AIUsefulInfoV4AndParts {
	ai := new(AIUsefulInfoV4AndParts)
	ai.BaseAI = *baseAI
	ai.Coefs = []AIUsefulInfoV4AndPartsCoefs{
		/*
			// Universal coefs
			AIUsefulInfoV4AndPartsCoefs{
				CoefPlayByInfoA: 2.47,
				CoefPlayByInfoB: -0.07,
				CoefInfoA:       2.7,
				CoefInfoB:       -0.61,
				CoefDiscardA:    0.7,
				CoefDiscardB:    0.7,
				CoefPlayA:       2.5,
				CoefPlayB:       -0.6,
			},
		*/
		AIUsefulInfoV4AndPartsCoefs{
			CoefPlayByInfoA: 2.47,
			CoefPlayByInfoB: -0.07,
			CoefInfoA:       2.7,
			CoefInfoB:       -0.61,
			CoefDiscardA:    0.7,
			CoefDiscardB:    0.7,
			CoefPlayA:       2.5,
			CoefPlayB:       -0.6,
		},
		AIUsefulInfoV4AndPartsCoefs{
			CoefPlayByInfoA: 2.64,
			CoefPlayByInfoB: 0.23,
			CoefInfoA:       2.6,
			CoefInfoB:       -0.06,
			CoefDiscardA:    0.7,
			CoefDiscardB:    0.7,
			CoefPlayA:       2.6,
			CoefPlayB:       -0.3,
		},
		AIUsefulInfoV4AndPartsCoefs{
			CoefPlayByInfoA: 1.87,
			CoefPlayByInfoB: 0.83,
			CoefInfoA:       3.2,
			CoefInfoB:       -0.61,
			CoefDiscardA:    0.7,
			CoefDiscardB:    0.7,
			CoefPlayA:       4.5,
			CoefPlayB:       0.7,
		},
	}
	return ai
}

func (ai *AIUsefulInfoV4AndParts) GetCoefs(part int) []float64 {
	if part >= len(ai.Coefs) || part < 0 {
		panic("Bad part for ai.GetCoefs()")
	}
	coefs := ai.Coefs[part]
	return []float64{
		coefs.CoefPlayByInfoA,
		coefs.CoefPlayByInfoB,
		coefs.CoefInfoA,
		coefs.CoefInfoB,
		coefs.CoefDiscardA,
		coefs.CoefDiscardB,
		coefs.CoefPlayA,
		coefs.CoefPlayB,
	}
}

func (ai *AIUsefulInfoV4AndParts) SetCoefs(part int, coefs ...float64) {
	if part >= len(ai.Coefs) || part < 0 {
		panic("Bad part for ai.SetCoefs()")
	}
	ai.Coefs[part] = AIUsefulInfoV4AndPartsCoefs{
		CoefPlayByInfoA: coefs[0],
		CoefPlayByInfoB: coefs[1],
		CoefInfoA:       coefs[2],
		CoefInfoB:       coefs[3],
		CoefDiscardA:    coefs[4],
		CoefDiscardB:    coefs[5],
		CoefPlayA:       coefs[6],
		CoefPlayB:       coefs[7],
	}
}

func (ai *AIUsefulInfoV4AndParts) GetPartOfGame() int {
	info := &ai.PlayerInfo
	if info.Step <= 16 {
		return 0
	}
	if info.DeckSize > 0 {
		return 1
	}
	return 2
}

func (ai *AIUsefulInfoV4AndParts) GetAction() *game.Action {
	info := &ai.PlayerInfo
	myPos := info.CurrentPostion
	oldPlayerCards := info.PlayerCards[myPos]
	info.PlayerCards[myPos] = info.PlayerCardsInfo[myPos]
	defer func() {
		info.PlayerCards[myPos] = oldPlayerCards
	}()
	ai.setAvailableInformation()
	usefulActions := Actions{}
	coefs := ai.Coefs[ai.GetPartOfGame()]

	subHistory := ai.History[Max(len(ai.History)-len(info.PlayerCards)+1, 0):]
	for i, action := range subHistory {
		if action.ActionType == game.TypeActionInformationValue && action.PlayerPosition == myPos {
			newCardIdxs := map[int]int{}
			step := len(ai.History) - len(subHistory) + i
			oldInfo := ai.Informator.GetPlayerState(step)
			for idx, card := range info.PlayerCards[myPos] {
				if !oldInfo.PlayerCards[myPos][idx].KnownValue && card.KnownValue && card.Value == game.CardValue(action.Value) {
					newCardIdxs[idx]++
				}
			}

			if len(newCardIdxs) == 0 {
				continue
			}

			valueIsValid := true
			for j := i + 1; j < len(subHistory); j++ {
				actionOld := &subHistory[j]
				if actionOld.ActionType == game.TypeActionPlaying {
					step := len(ai.History) - len(subHistory) + j
					oldInfo := ai.Informator.GetPlayerState(step)
					if actionOld.PlayerPosition != oldInfo.CurrentPostion {
						panic("Bad CurrentPostion")
					}
					card := &oldInfo.PlayerCards[actionOld.PlayerPosition][actionOld.Value]
					if game.CardValue(action.Value) == card.Value {
						valueIsValid = false
						break
					}
				}
			}

			if !valueIsValid {
				continue
			}

			for idx, _ := range info.PlayerCards[myPos] {
				if _, ok := newCardIdxs[idx]; ok {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByInfoA*float64(i)/float64(len(subHistory)) + coefs.CoefPlayByInfoB,
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}

		if action.ActionType == game.TypeActionInformationColor && action.PlayerPosition == myPos {
			newCardIdxs := map[int]int{}
			step := len(ai.History) - len(subHistory) + i
			oldInfo := ai.Informator.GetPlayerState(step)
			for idx, card := range info.PlayerCards[myPos] {
				if !oldInfo.PlayerCards[myPos][idx].KnownValue && card.KnownColor && card.Color == game.CardColor(action.Value) {
					newCardIdxs[idx]++
				}
			}

			if len(newCardIdxs) == 0 {
				continue
			}

			colorIsValid := true
			for j := i + 1; j < len(subHistory); j++ {
				actionOld := &subHistory[j]
				if actionOld.ActionType == game.TypeActionPlaying {
					step := len(ai.History) - len(subHistory) + j
					oldInfo := ai.Informator.GetPlayerState(step)
					if actionOld.PlayerPosition != oldInfo.CurrentPostion {
						panic("Bad CurrentPostion")
					}
					card := &oldInfo.PlayerCards[actionOld.PlayerPosition][actionOld.Value]
					if game.CardColor(action.Value) == card.Color {
						colorIsValid = false
						break
					}
				}
			}

			if !colorIsValid {
				continue
			}

			for idx, _ := range info.PlayerCards[myPos] {
				if _, ok := newCardIdxs[idx]; ok {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByInfoA*float64(i)/float64(len(subHistory)) + coefs.CoefPlayByInfoB,
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}
	}

	if info.BlueTokens > 0 {
		for i := 1; i < len(info.PlayerCards); i++ {
			nextPos := (myPos + i) % len(info.PlayerCards)
			cards := info.PlayerCards[nextPos]
			cardsInfo := info.PlayerCardsInfo[nextPos]

			usefulCards := []int{}
			for idx, card := range cards {
				tableCard := info.TableCards[card.Color]
				if card.Value == tableCard.Value+1 {
					usefulCards = append(usefulCards, idx)
				}
			}

			if len(usefulCards) == 0 {
				continue
			}

			maxUsefulByValue := 0
			maxUsefulByColor := 0
			maxNeutralByValue := 0
			maxNeutralByColor := 0
			minAdverseByValue := 5
			minAdverseByColor := 5
			var minCardValueByValue game.CardValue
			var minCardValueByColor game.CardValue
			var idxValue int
			var idxColor int
			subHistory := ai.History[Max(len(ai.History)-len(info.PlayerCards)+1-i, 0):]
			infoColors := map[game.CardColor]struct{}{}
			infoValues := map[game.CardValue]struct{}{}
			for _, action := range subHistory {
				if action.PlayerPosition != nextPos {
					continue
				}
				if action.ActionType == game.TypeActionInformationValue {
					infoValues[game.CardValue(action.Value)] = struct{}{}
				}

				if action.ActionType == game.TypeActionInformationColor {
					infoColors[game.CardColor(action.Value)] = struct{}{}
				}
			}

			for _, idx := range usefulCards {
				card := &cards[idx]
				cardInfo := &cardsInfo[idx]
				usefulByValue := 0
				usefulByColor := 0
				neutralByValue := 0
				neutralByColor := 0
				adverseByValue := 0
				adverseByColor := 0
				for j := 0; j < len(cards); j++ {
					if _, ok := infoValues[card.Value]; !ok && cards[j].Value == card.Value {
						if info.TableCards[cards[j].Color].Value+1 == cards[j].Value {
							if !cardInfo.KnownValue {
								usefulByValue++
							}
						} else if cardInfo.KnownColor {
							neutralByValue++
						} else {
							adverseByValue++
						}
					}

					if _, ok := infoColors[card.Color]; !ok && cards[j].Color == card.Color {
						if info.TableCards[cards[j].Color].Value+1 == cards[j].Value {
							if !cardInfo.KnownColor {
								usefulByColor++
							}
						} else if cardInfo.KnownValue {
							neutralByColor++
						} else {
							adverseByColor++
						}
					}

				}

				isBetter := adverseByValue < minAdverseByValue ||
					adverseByValue == minAdverseByValue && card.Value < minCardValueByValue ||
					adverseByValue == minAdverseByValue && card.Value == minCardValueByValue && usefulByValue > maxUsefulByValue ||
					adverseByValue == minAdverseByValue && card.Value == minCardValueByValue && usefulByValue == maxUsefulByValue && neutralByValue > maxNeutralByValue
				if isBetter {
					maxUsefulByValue = usefulByValue
					maxNeutralByValue = neutralByValue
					minAdverseByValue = adverseByValue
					minCardValueByValue = card.Value
					idxValue = idx
				}

				isBetter = adverseByColor < minAdverseByColor ||
					adverseByColor == minAdverseByColor && card.Value < minCardValueByColor ||
					adverseByColor == minAdverseByColor && card.Value == minCardValueByColor && usefulByColor > maxUsefulByColor ||
					adverseByColor == minAdverseByColor && card.Value == minCardValueByColor && usefulByColor == maxUsefulByColor && neutralByColor > maxNeutralByColor
				if isBetter {
					maxUsefulByColor = usefulByColor
					maxNeutralByColor = neutralByColor
					minAdverseByColor = adverseByColor
					minCardValueByColor = card.Value
					idxColor = idx
				}
			}

			if maxUsefulByValue == 0 && maxUsefulByColor == 0 && maxNeutralByValue == 0 && maxNeutralByColor == 0 {
				continue
			}

			isBetterInfoByValue := minAdverseByValue < minAdverseByColor ||
				minAdverseByValue == minAdverseByColor && minCardValueByValue < minCardValueByColor ||
				minAdverseByValue == minAdverseByColor && minCardValueByValue == minCardValueByColor && maxUsefulByValue > maxUsefulByColor ||
				minAdverseByValue == minAdverseByColor && minCardValueByValue == minCardValueByColor && maxUsefulByValue == maxUsefulByColor && maxNeutralByValue > maxNeutralByColor ||
				minAdverseByValue == minAdverseByColor && minCardValueByValue == minCardValueByColor && maxUsefulByValue == maxUsefulByColor && maxNeutralByValue == maxNeutralByColor

			if isBetterInfoByValue {
				card := &cards[idxValue]
				bonus := -float64(i)/1000 - float64(minAdverseByValue)/1001 + float64(maxUsefulByValue)/100 + float64(maxNeutralByValue)/1002
				action := UsefulAction{
					Action:     game.NewAction(game.TypeActionInformationValue, nextPos, int(card.Value)),
					Usefulness: coefs.CoefInfoA*(1.0-float64(i)/float64(len(info.PlayerCards))+bonus) + coefs.CoefInfoB,
				}
				usefulActions = append(usefulActions, action)
			} else {
				card := &cards[idxColor]
				bonus := -float64(i)/1000 - float64(minAdverseByColor)/1001 + float64(maxUsefulByColor)/100 + float64(maxNeutralByColor)/1002
				action := UsefulAction{
					Action:     game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color)),
					Usefulness: coefs.CoefInfoA*(1.0-float64(i)/float64(len(info.PlayerCards))+bonus) + coefs.CoefInfoB,
				}
				usefulActions = append(usefulActions, action)
			}
		}
	}

	for idx, card := range info.PlayerCards[myPos] {
		usefulnessNow := 0.0
		usefulnessPotential := 0.0
		if card.KnownColor && card.KnownValue {
			tableValue := info.TableCards[card.Color].Value
			if tableValue+1 == card.Value {
				usefulnessNow = 1.0
			} else {
				usefulnessNow = 0.0
			}

			if tableValue < card.Value {
				usefulnessPotential = 1.0
			} else {
				usefulnessPotential = 0.0
			}
		} else {
			for hashValue, probability := range card.ProbabilityCard {
				color, value := game.ColorValueByHashColorValue(hashValue)
				if info.TableCards[color].Value+1 == value {
					usefulnessNow += probability
				}

				colorValue := game.ColorValue{
					Color: color,
					Value: value,
				}
				if info.TableCards[color].Value < value && info.VariantsCount[colorValue] == 1 {
					usefulnessPotential += probability
				}
			}
		}

		if info.BlueTokens < game.MaxBlueTokens {
			actionDiscard := UsefulAction{
				Action:     game.NewAction(game.TypeActionDiscard, myPos, idx),
				Usefulness: coefs.CoefDiscardA*math.Sin(math.Pow(1-usefulnessPotential, 2)*math.Pi/2) + coefs.CoefDiscardB - float64(idx)/1000,
			}
			usefulActions = append(usefulActions, actionDiscard)
		}

		actionPlay := UsefulAction{
			Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
			Usefulness: coefs.CoefPlayA*math.Sin(math.Pow(usefulnessNow, 2)*math.Pi/2) + coefs.CoefPlayB - float64(idx)/1000,
		}
		usefulActions = append(usefulActions, actionPlay)
	}

	if len(usefulActions) > 0 {
		bestActionIdx := 0
		for i := 1; i < len(usefulActions); i++ {
			if usefulActions.Less(i, bestActionIdx) {
				bestActionIdx = i
			}
		}
		return usefulActions[bestActionIdx].Action
	}
	return ai.getActionSmartyRandom()
}
