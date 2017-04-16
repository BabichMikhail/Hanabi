package ai

import (
	"math"

	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInfoV4AndPartsCoefs struct {
	CoefPlayByValueA float64
	CoefPlayByValueB float64
	CoefPlayByColorA float64
	CoefPlayByColorB float64
	CoefInfoValueA   float64
	CoefInfoValueB   float64
	CoefInfoColorA   float64
	CoefInfoColorB   float64
	CoefDiscardA     float64
	CoefDiscardB     float64
	CoefPlayA        float64
	CoefPlayB        float64
}

type AIUsefulInfoV4AndParts struct {
	BaseAI

	Coefs []AIUsefulInfoV4AndPartsCoefs
}

func NewAIUsefulInfoV4AndParts(baseAI *BaseAI) *AIUsefulInfoV4AndParts {
	ai := new(AIUsefulInfoV4AndParts)
	ai.BaseAI = *baseAI
	ai.Coefs = []AIUsefulInfoV4AndPartsCoefs{
		AIUsefulInfoV4AndPartsCoefs{
			CoefPlayByValueA: 2.6,
			CoefPlayByValueB: 1.0,
			CoefPlayByColorA: -1.0,
			CoefPlayByColorB: 0.0,
			CoefInfoValueA:   1.05,
			CoefInfoValueB:   0.0,
			CoefInfoColorA:   0.5,
			CoefInfoColorB:   0.0,
			CoefDiscardA:     0.0,
			CoefDiscardB:     0.0,
			CoefPlayA:        1.1,
			CoefPlayB:        0.0,
		},
		AIUsefulInfoV4AndPartsCoefs{
			CoefPlayByValueA: 1.1,
			CoefPlayByValueB: 0.0,
			CoefPlayByColorA: -1.0,
			CoefPlayByColorB: 0.0,
			CoefInfoValueA:   1.05,
			CoefInfoValueB:   0.0,
			CoefInfoColorA:   1.0,
			CoefInfoColorB:   0.0,
			CoefDiscardA:     1.0,
			CoefDiscardB:     0.0,
			CoefPlayA:        0.1,
			CoefPlayB:        0.0,
		},
		AIUsefulInfoV4AndPartsCoefs{
			CoefPlayByValueA: 1.1,
			CoefPlayByValueB: 0.0,
			CoefPlayByColorA: -1.0,
			CoefPlayByColorB: 0.0,
			CoefInfoValueA:   1.05,
			CoefInfoValueB:   0.0,
			CoefInfoColorA:   1.0,
			CoefInfoColorB:   0.0,
			CoefDiscardA:     1.0,
			CoefDiscardB:     0.0,
			CoefPlayA:        0.1,
			CoefPlayB:        0.0,
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
		coefs.CoefPlayByValueA,
		coefs.CoefPlayByValueB,
		coefs.CoefPlayByColorA,
		coefs.CoefPlayByColorB,
		coefs.CoefInfoValueA,
		coefs.CoefInfoValueB,
		coefs.CoefInfoColorA,
		coefs.CoefInfoColorB,
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
		CoefPlayByValueA: coefs[0],
		CoefPlayByValueB: coefs[1],
		CoefPlayByColorA: coefs[2],
		CoefPlayByColorB: coefs[3],
		CoefInfoValueA:   coefs[4],
		CoefInfoValueB:   coefs[5],
		CoefInfoColorA:   coefs[6],
		CoefInfoColorB:   coefs[7],
		CoefDiscardA:     coefs[8],
		CoefDiscardB:     coefs[9],
		CoefPlayA:        coefs[10],
		CoefPlayB:        coefs[11],
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
	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[myPos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, myPos, idx)
			}
		}
	}

	subHistory := ai.History[Max(len(ai.History)-len(info.PlayerCards)+1, 0):]
	for i, action := range subHistory {
		if action.ActionType == game.TypeActionInformationValue && action.PlayerPosition == myPos {
			newCardIdxs := []int{}
			step := len(ai.History) - len(subHistory) + i
			oldInfo := ai.Informator.GetPlayerState(step)
			for idx, card := range info.PlayerCards[myPos] {
				if !oldInfo.PlayerCards[myPos][idx].KnownValue && card.KnownValue && card.Value == game.CardValue(action.Value) {
					newCardIdxs = append(newCardIdxs, idx)
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
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByValueA/math.Sqrt(float64(len(newCardIdxs))) + coefs.CoefPlayByValueB,
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}

		if action.ActionType == game.TypeActionInformationColor && action.PlayerPosition == myPos {
			newCardIdxs := []int{}
			step := len(ai.History) - len(subHistory) + i
			oldInfo := ai.Informator.GetPlayerState(step)
			for idx, card := range info.PlayerCards[myPos] {
				if !oldInfo.PlayerCards[myPos][idx].KnownValue && card.KnownColor && card.Color == game.CardColor(action.Value) {
					newCardIdxs = append(newCardIdxs, idx)
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
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByColorA/math.Sqrt(float64(len(newCardIdxs))) + coefs.CoefPlayByColorB,
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}
	}

	if info.BlueTokens > 0 {
		for i := 1; i < len(info.PlayerCards); i++ {
			nextPos := (myPos + i) % len(info.PlayerCards)
			for color, tableCard := range info.TableCards {
				for idx, card := range info.PlayerCards[nextPos] {
					if card.Color == color && card.Value == tableCard.Value+1 {
						cardInfo := &info.PlayerCardsInfo[nextPos][idx]
						if !cardInfo.KnownValue {
							action := UsefulAction{
								Action:     game.NewAction(game.TypeActionInformationValue, nextPos, int(card.Value)),
								Usefulness: coefs.CoefInfoValueA*(1.0-float64(i)/float64(len(info.PlayerCards))) + coefs.CoefInfoValueB - float64(i)/10,
							}
							usefulActions = append(usefulActions, action)
						}

						if !cardInfo.KnownColor {
							action := UsefulAction{
								Action:     game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color)),
								Usefulness: coefs.CoefInfoColorA*(1.0-float64(i)/float64(len(info.PlayerCards))) + coefs.CoefInfoColorB - float64(i)/10,
							}
							usefulActions = append(usefulActions, action)
						}
					}
				}
			}
		}
	}

	for idx, card := range info.PlayerCards[myPos] {
		usefulnessNow := 0.0
		usefulnessPotential := 0.0
		if card.KnownColor && card.KnownValue {
			usefulnessNow = 1.0
			usefulnessPotential = 1.0
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
				Usefulness: coefs.CoefDiscardA*math.Pow(1-usefulnessPotential, 1) + coefs.CoefDiscardB,
			}
			usefulActions = append(usefulActions, actionDiscard)
		}

		actionPlay := UsefulAction{
			Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
			Usefulness: coefs.CoefPlayA*math.Pow(usefulnessNow, 1) + coefs.CoefPlayB,
		}
		usefulActions = append(usefulActions, actionPlay)
	}

	if len(usefulActions) > 0 {
		var bestAction *game.Action
		var topUsefulness float64
		for i := 0; i < len(usefulActions); i++ {
			if bestAction == nil || topUsefulness < usefulActions[i].Usefulness {
				bestAction = usefulActions[i].Action
				topUsefulness = usefulActions[i].Usefulness
			}
		}
		return bestAction
	}
	return ai.getActionSmartyRandom()
}
