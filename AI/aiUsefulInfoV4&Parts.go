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
	isUniversal bool
	SafeMode    bool
	Coefs       map[int][]AIUsefulInfoV4AndPartsCoefs
}

func NewAIUsefulInfoV4AndParts(baseAI *BaseAI, isUniversal bool) *AIUsefulInfoV4AndParts {
	ai := new(AIUsefulInfoV4AndParts)
	ai.BaseAI = *baseAI
	ai.isUniversal = isUniversal
	ai.SafeMode = true
	if ai.isUniversal {
		ai.Coefs = map[int][]AIUsefulInfoV4AndPartsCoefs{
			2: []AIUsefulInfoV4AndPartsCoefs{
				AIUsefulInfoV4AndPartsCoefs{
					CoefPlayByInfoA: 2.46,
					CoefPlayByInfoB: 1.23,
					CoefInfoA:       2.3,
					CoefInfoB:       0.89,
					CoefDiscardA:    0.49,
					CoefDiscardB:    0.7,
					CoefPlayA:       2,
					CoefPlayB:       -0.6,
				},
			},
			5: []AIUsefulInfoV4AndPartsCoefs{
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
			},
		}
	} else {
		ai.Coefs = map[int][]AIUsefulInfoV4AndPartsCoefs{
			5: []AIUsefulInfoV4AndPartsCoefs{
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
			},
		}
	}

	return ai
}

func (ai *AIUsefulInfoV4AndParts) GetCoefs(part int) []float64 {
	if part >= len(ai.Coefs) || part < 0 {
		panic("Bad part for ai.GetCoefs()")
	}
	coefs := ai.Coefs[len(ai.PlayerInfo.PlayerCards)][part]
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
	ai.Coefs[len(ai.PlayerInfo.PlayerCards)][part] = AIUsefulInfoV4AndPartsCoefs{
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
	if ai.isUniversal {
		return 0
	}
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
	coefsByPlayerCount, ok := ai.Coefs[len(ai.PlayerInfo.PlayerCards)]
	if !ok {
		panic("Coefs are undefined")
	}
	coefs := coefsByPlayerCount[ai.GetPartOfGame()]

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

			for idx, card := range info.PlayerCards[myPos] {
				if _, ok := newCardIdxs[idx]; ai.isCardPlayable(card) && ok {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByInfoA*float64(i)/float64(len(subHistory)) + coefs.CoefPlayByInfoB,
					}
					usefulActions = append(usefulActions, action)

					/*usefulnessNow := 0.0
					if card.KnownColor && card.KnownValue {
						tableValue := info.TableCards[card.Color].Value
						if tableValue+1 == card.Value {
							usefulnessNow = 1.0
						} else {
							usefulnessNow = 0.0
						}
					} else {
						for hashValue, probability := range card.ProbabilityCard {
							color, value := game.ColorValueByHashColorValue(hashValue)
							if info.TableCards[color].Value+1 == value {
								usefulnessNow += probability
							}
						}
					}

					fmt.Println("USEFULNESS NOW: ", usefulnessNow)*/
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

			for idx, card := range info.PlayerCards[myPos] {
				if _, ok := newCardIdxs[idx]; ai.isCardPlayable(card) && ok {
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
			newInfo := info.Copy()
			newInfo.PlayerCards[nextPos] = newInfo.PlayerCardsInfo[nextPos]
			newInfo.SetProbabilities(false, false)
			cards := info.PlayerCards[nextPos]
			cardsInfo := newInfo.PlayerCardsInfo[nextPos]

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

			testIsBetter := func(adverse1, adverse2 int, value1, value2 game.CardValue, useful1, useful2, neutral1, neutral2 int) bool {
				return adverse1 < adverse2 ||
					adverse1 == adverse2 && value1 < value2 ||
					adverse1 == adverse2 && value1 == value2 && useful1 > useful2 ||
					adverse1 == adverse2 && value1 == value2 && useful1 == useful2 && neutral1 > neutral2
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
				threshold := 0.3
				for j := 0; j < len(cards); j++ {
					if _, ok := infoValues[card.Value]; !ok && cards[j].Value == card.Value {
						if info.TableCards[cards[j].Color].Value+1 == cards[j].Value {
							if !cardInfo.KnownValue {
								isUseful := true
								for _, prob := range cards[j].ProbabilityCard {
									if prob > threshold {
										isUseful = false
										break
									}
								}
								if isUseful {
									usefulByValue++
								} else {
									neutralByValue++
								}

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
								isUseful := true
								for _, prob := range cards[j].ProbabilityCard {
									if prob > threshold {
										isUseful = false
										break
									}
								}
								if isUseful {
									usefulByColor++
								} else {
									neutralByColor++
								}
							}
						} else if cardInfo.KnownValue {
							neutralByColor++
						} else {
							adverseByColor++
						}
					}

				}

				isBetter := testIsBetter(
					adverseByValue, minAdverseByValue, card.Value, minCardValueByValue,
					usefulByValue, maxUsefulByValue, neutralByValue, maxNeutralByValue,
				)

				if isBetter {
					maxUsefulByValue = usefulByValue
					maxNeutralByValue = neutralByValue
					minAdverseByValue = adverseByValue
					minCardValueByValue = card.Value
					idxValue = idx
				}

				isBetter = testIsBetter(
					adverseByColor, minAdverseByColor, card.Value, minCardValueByColor,
					usefulByColor, maxUsefulByColor, neutralByColor, maxNeutralByColor,
				)

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
			isBetterInfoByValue := testIsBetter(
				minAdverseByValue, minAdverseByColor, minCardValueByValue, minCardValueByColor,
				maxUsefulByValue, maxUsefulByColor, maxNeutralByValue, maxNeutralByColor,
			)

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
		bestAction := usefulActions[bestActionIdx].Action
		if !ai.SafeMode || info.BlueTokens == 0 {
			return bestAction
		}

		newInfo := info.Copy()
		newInfo.PreviewAction(bestAction)
		newAI := NewAI(*newInfo, append(ai.History, *bestAction), ai.Type, ai.Informator)
		newAI.(*AIUsefulInfoV4AndParts).SafeMode = false
		nextAction := newAI.GetAction()
		if nextAction.ActionType == game.TypeActionDiscard || nextAction.ActionType == game.TypeActionPlaying {
			card := &info.PlayerCards[nextAction.PlayerPosition][nextAction.Value]
			if info.VariantsCount[game.ColorValue{Color: card.Color, Value: card.Value}] == 1 {
				nextPos := nextAction.PlayerPosition
				cardPos := nextAction.Value
				if info.TableCards[card.Color].Value+1 == card.Value {
					if nextAction.ActionType == game.TypeActionDiscard {
						return game.NewAction(game.TypeActionInformationValue, nextPos, cardPos)
					} else {
						return bestAction
					}
				} else {
					card := newInfo.PlayerCardsInfo[nextPos][cardPos]
					if ai.isCardPlayable(card) {
						copyCard := card.Copy()
						if !card.KnownValue {
							card.Value = newInfo.PlayerCards[nextPos][cardPos].Value
							card.KnownValue = true
							for _, color := range game.ColorsTable {
								hashValue := game.HashColorValue(color, card.Value)
								if _, ok := card.ProbabilityCard[hashValue]; ok {
									delete(card.ProbabilityCard, hashValue)
								}
							}

							if !ai.isCardPlayable(card) {
								return game.NewAction(game.TypeActionInformationValue, nextPos, cardPos)
							}
						}

						card = copyCard
						if !card.KnownColor {
							card.Color = newInfo.PlayerCards[nextPos][cardPos].Color
							card.KnownColor = true
							for _, value := range game.Values {
								hashValue := game.HashColorValue(card.Color, value)
								if _, ok := card.ProbabilityCard[hashValue]; ok {
									delete(card.ProbabilityCard, hashValue)
								}
							}

							if !ai.isCardPlayable(card) {
								return game.NewAction(game.TypeActionInformationColor, nextPos, cardPos)
							}
						}
					}
				}
			}
		}

		return bestAction
	}
	return ai.getActionSmartyRandom()
}
