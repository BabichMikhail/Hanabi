package ai

import (
	"math"
	"math/rand"
	"sort"

	"github.com/BabichMikhail/Hanabi/game"
)

type AIUsefulInfoV3Coefs struct {
	CoefPlayByValue           float64
	CoefPlayByColor           float64
	CoefInfoValue             float64
	CoefInfoColor             float64
	CoefDiscardUselessCard    float64
	CoefDiscardUnknownCard    float64
	CoefDiscardUsefulCard     float64
	CoefDiscardMaybeUsefuCard float64
}

type AIUsefulInformationV3 struct {
	BaseAI
	PartOfGame int
	Coefs      []AIUsefulInfoV3Coefs
}

func NewAIUsefulInformationV3(baseAI *BaseAI) *AIUsefulInformationV3 {
	ai := new(AIUsefulInformationV3)
	ai.BaseAI = *baseAI
	ai.Coefs = []AIUsefulInfoV3Coefs{
		AIUsefulInfoV3Coefs{
			CoefPlayByValue:           2.1,
			CoefPlayByColor:           -0.9,
			CoefInfoValue:             1.05,
			CoefInfoColor:             1.0,
			CoefDiscardUsefulCard:     0.1,
			CoefDiscardMaybeUsefuCard: 0.04,
			CoefDiscardUselessCard:    0.01,
			CoefDiscardUnknownCard:    0.07,
		},
	}

	return ai
}

func (ai *AIUsefulInformationV3) GetCoefs(part int) []float64 {
	if part >= len(ai.Coefs) || part < 0 {
		panic("Bad part for ai.GetCoefs()")
	}
	coefs := ai.Coefs[part]
	return []float64{
		coefs.CoefPlayByValue,
		coefs.CoefPlayByColor,
		coefs.CoefInfoValue,
		coefs.CoefInfoColor,
		coefs.CoefDiscardUsefulCard,
		coefs.CoefDiscardMaybeUsefuCard,
		coefs.CoefDiscardUselessCard,
		coefs.CoefDiscardUnknownCard,
	}
}

func (ai *AIUsefulInformationV3) SetCoefs(part int, coefs ...float64) {
	if part >= len(ai.Coefs) || part < 0 {
		panic("Bad part for ai.SetCoefs()")
	}

	ai.Coefs[part] = AIUsefulInfoV3Coefs{
		CoefPlayByValue:           coefs[0],
		CoefPlayByColor:           coefs[1],
		CoefInfoValue:             coefs[2],
		CoefInfoColor:             coefs[3],
		CoefDiscardUsefulCard:     coefs[4],
		CoefDiscardMaybeUsefuCard: coefs[5],
		CoefDiscardUselessCard:    coefs[6],
		CoefDiscardUnknownCard:    coefs[7],
	}
}

func (ai *AIUsefulInformationV3) GetPartOfGame() int {
	return 0
}

func (ai *AIUsefulInformationV3) GetAction() *game.Action {
	info := &ai.PlayerInfo
	myPos := info.CurrentPostion
	oldPlayerCards := info.PlayerCards[myPos]
	info.PlayerCards[myPos] = info.PlayerCardsInfo[myPos]
	defer func() {
		info.PlayerCards[myPos] = oldPlayerCards
	}()

	ai.setAvailableInformation()
	for color, tableCard := range info.TableCards {
		for idx, card := range info.PlayerCards[myPos] {
			if card.KnownColor && card.KnownValue && card.Color == color && card.Value == tableCard.Value+1 {
				return game.NewAction(game.TypeActionPlaying, myPos, idx)
			}
		}
	}

	usefulActions := Actions{}
	coefs := ai.Coefs[ai.GetPartOfGame()]
	subHistory := ai.History[Max(len(ai.History)-len(info.PlayerCards)+1, 0):]
	historyLength := len(subHistory)
	for i, action := range subHistory {
		if action.ActionType == game.TypeActionInformationValue && action.PlayerPosition == myPos {
			count := 0.0
			for _, card := range info.PlayerCards[myPos] {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					count++
				}
			}

			if count == 0 {
				continue
			}

			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByValue / float64(historyLength-i) / math.Sqrt(count),
					}
					usefulActions = append(usefulActions, action)
				}
			}
		}

		if action.ActionType == game.TypeActionInformationColor && action.PlayerPosition == myPos {
			count := 0.0
			for _, card := range info.PlayerCards[myPos] {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					count++
				}
			}

			if count == 0 {
				continue
			}

			for idx, card := range info.PlayerCards[myPos] {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					action := UsefulAction{
						Action:     game.NewAction(game.TypeActionPlaying, myPos, idx),
						Usefulness: coefs.CoefPlayByColor / float64(historyLength-i) / math.Sqrt(count),
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
								Usefulness: coefs.CoefInfoValue * (1.0 - float64(i)/float64(len(info.PlayerCards))),
							}
							usefulActions = append(usefulActions, action)
						}

						if !cardInfo.KnownColor {
							action := UsefulAction{
								Action:     game.NewAction(game.TypeActionInformationColor, nextPos, int(card.Color)),
								Usefulness: coefs.CoefInfoColor * (1.0 - float64(i)/float64(len(info.PlayerCards))),
							}
							usefulActions = append(usefulActions, action)
						}
					}
				}
			}
		}
	}

	if info.BlueTokens < game.MaxBlueTokens {
		for idx, card := range info.PlayerCards[myPos] {
			var coef float64
			if card.KnownColor && card.KnownValue {
				if card.Value > info.TableCards[card.Color].Value {
					coef = coefs.CoefDiscardUsefulCard
				} else {
					coef = coefs.CoefDiscardUselessCard
				}
			} else if card.KnownValue {
				coef = coefs.CoefDiscardUselessCard
				for _, card := range info.TableCards {
					if card.Value+1 == card.Value {
						coef = coefs.CoefDiscardMaybeUsefuCard
					}
				}
			} else if card.KnownColor {
				if info.TableCards[card.Color].Value == 5 {
					coef = coefs.CoefDiscardUselessCard
				} else {
					coef = coefs.CoefDiscardMaybeUsefuCard
				}
			} else {
				coef = coefs.CoefDiscardUnknownCard
			}
			action := UsefulAction{
				Action:     game.NewAction(game.TypeActionDiscard, myPos, idx),
				Usefulness: coef,
			}
			usefulActions = append(usefulActions, action)
		}
	}

	if len(usefulActions) > 0 {
		sort.Sort(usefulActions)
		return usefulActions[0].Action
	}

	if info.BlueTokens < game.MaxBlueTokens {
		return game.NewAction(game.TypeActionDiscard, myPos, rand.Intn(len(info.PlayerCards[myPos])))
	}

	return ai.getActionSmartyRandom()
}
