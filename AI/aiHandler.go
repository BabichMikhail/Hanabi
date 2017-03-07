package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

const (
	AI_RandomAction = iota
	AI_SmartyRandomAction
	AI_DiscardUsefullCardAction
	AI_UsefullInformationAction
)

type Action struct {
	game.Action
	UsefullCount int `json:"usefull_count"`
	Count        int `json:"count"`
}

type AI struct {
	Actions          []*Action           `json:"actions"`
	PlayActions      []*Action           `json:"playing_actions"`
	DiscardActions   []*Action           `json:"discard_actions"`
	InfoValueActions []*Action           `json:"info_value_actions"`
	InfoColorAcions  []*Action           `json:"info_color_actions"`
	History          []game.Action       `json:"history"`
	PlayerInfo       game.PlayerGameInfo `json:"player_info"`
	Type             int                 `json:"ai_type"`
}

type Card struct {
	game.Card
	ProbabilityValues map[game.CardValue]float64
	ProbabilityColors map[game.CardColor]float64
}

func NewCard(gameCard game.Card) *Card {
	return &Card{gameCard, map[game.CardValue]float64{}, map[game.CardColor]float64{}}
}

const (
	AI_NamePrefix = "AI_"

	AI_RandomName             = "RandomAction"
	AI_SmartyName             = "SmartyRandomAction"
	AI_DiscardUsefullCardName = "DiscardKnownCardAction"
	AI_UsefullInformationName = "UsefullInformationAction"
)

func NewAI(playerInfo game.PlayerGameInfo, actions []game.Action, aiType int) *AI {
	ai := new(AI)
	ai.History = actions
	ai.PlayerInfo = playerInfo
	ai.setAvailableActions()
	ai.Type = aiType
	return ai
}

func (ai *AI) GetAction() game.Action {
	switch ai.Type {
	case AI_RandomAction:
		return ai.getActionRandom()
	case AI_SmartyRandomAction:
		return ai.getActionSmartyRandom()
	case AI_DiscardUsefullCardAction:
		return ai.getActionDiscardUsefullCard()
	case AI_UsefullInformationAction:
		return ai.getActionUsefullInformation()
	default:
		panic("Missing AI_Type")
	}
}

func (ai *AI) setProbabilityValues(card *game.Card, cardsCount map[game.HashValue]int, color game.CardColor, delta int) {
	count := 0
	for value, _ := range card.ProbabilityValues {
		count += cardsCount[game.HashColorValue(color, value)]
	}
	count += cardsCount[game.HashColorValue(color, game.NoneValue)] + delta

	for value, _ := range card.ProbabilityValues {
		card.ProbabilityValues[value] = float64(cardsCount[game.HashColorValue(color, value)]) / float64(count)
	}
}

func (ai *AI) setProbabilityColors(card *game.Card, cardsCount map[game.HashValue]int, value game.CardValue, delta int) {
	count := 0
	for color, _ := range card.ProbabilityColors {
		count += cardsCount[game.HashColorValue(color, value)]
	}
	//count += cardsCount[game.HashColorValue(game.NoneColor, value)] + delta
	for color, _ := range card.ProbabilityColors {
		card.ProbabilityColors[color] = float64(cardsCount[game.HashColorValue(color, value)]) / float64(count)
	}
}

func (ai *AI) setProbabilities() {
	info := &ai.PlayerInfo

	cardsCount := map[game.HashValue]int{}
	for _, color := range append(game.Colors, game.NoneColor) {
		cardsCount[game.HashColorValue(color, game.One)] = 3
		cardsCount[game.HashColorValue(color, game.Two)] = 2
		cardsCount[game.HashColorValue(color, game.Three)] = 2
		cardsCount[game.HashColorValue(color, game.Four)] = 2
		cardsCount[game.HashColorValue(color, game.Five)] = 1
		cardsCount[game.HashColorValue(color, game.NoneValue)] = 0
	}

	for _, card := range info.UsedCards {
		cardsCount[game.HashColorValue(card.Color, card.Value)]--
	}

	for _, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if card.KnownColor {
				cardsCount[game.HashColorValue(card.Color, game.NoneValue)]--
			}

			if card.KnownValue {
				cardsCount[game.HashColorValue(game.NoneColor, card.Value)]--
			}
		}
	}

	for _, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if card.KnownColor && card.KnownValue {
				card.ProbabilityColors[card.Color] = 1.0
				card.ProbabilityValues[card.Value] = 1.0
				card.ProbabilityCard[game.HashColorValue(card.Color, card.Value)] = 1.0
			} else if card.KnownValue {
				card.ProbabilityValues[card.Value] = 1.0
				ai.setProbabilityColors(card, cardsCount, card.Value, 1)
			} else if card.KnownColor {
				card.ProbabilityColors[card.Color] = 1.0
				ai.setProbabilityValues(card, cardsCount, card.Color, 1)
			} else {
				for _, color := range game.Colors {
					ai.setProbabilityValues(card, cardsCount, color, 0)
				}

				for _, value := range game.Values {
					ai.setProbabilityColors(card, cardsCount, value, 0)
				}
			}

			for color, _ := range card.ProbabilityColors {
				for value, _ := range card.ProbabilityValues {
					card.ProbabilityCard[game.HashColorValue(color, value)] = card.ProbabilityColors[color] * card.ProbabilityValues[value]
				}
			}

		}
	}
}

func (ai *AI) setAvailableInfomation() {
	info := &ai.PlayerInfo
	for _, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if len(card.ProbabilityColors) == 1 {
				for color, _ := range card.ProbabilityColors {
					card.KnownColor = true
					card.Color = color
				}
			}

			if len(card.ProbabilityValues) == 1 {
				for value, _ := range card.ProbabilityValues {
					card.KnownValue = true
					card.Value = value
				}
			}
		}
	}

	ai.setProbabilities()
}
