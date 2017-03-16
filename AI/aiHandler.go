package ai

import (
	"reflect"

	"github.com/BabichMikhail/Hanabi/game"
)

const (
	AI_RandomAction = iota
	AI_SmartyRandomAction
	AI_DiscardUsefullCardAction
	AI_UsefullInformationAction
)

var AITypes = []int{
	AI_RandomAction,
	AI_SmartyRandomAction,
	AI_DiscardUsefullCardAction,
	AI_UsefullInformationAction,
}

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

var AINames = map[int]string{
	AI_RandomAction:             AI_RandomName,
	AI_SmartyRandomAction:       AI_SmartyName,
	AI_DiscardUsefullCardAction: AI_DiscardUsefullCardName,
	AI_UsefullInformationAction: AI_UsefullInformationName,
}

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

func (ai *AI) setProbabilityValues(card *game.Card, cardsCount map[game.HashValue]int, color game.CardColor, delta int) bool {
	count := 0
	for value, _ := range card.ProbabilityValues {
		count += cardsCount[game.HashColorValue(color, value)]
	}
	count += cardsCount[game.HashColorValue(color, game.NoneValue)] + delta
	correction := float64(count) / float64(count-cardsCount[game.HashColorValue(color, game.NoneValue)]-delta)

	for value, _ := range card.ProbabilityValues {
		if cardsCount[game.HashColorValue(color, value)] == 0 {
			delete(card.ProbabilityValues, value)
		}
		card.ProbabilityValues[value] = correction * float64(cardsCount[game.HashColorValue(color, value)]) / float64(count)
	}

	if len(card.ProbabilityValues) == 1 {
		keys := reflect.ValueOf(card.ProbabilityValues).MapKeys()
		value := keys[0].Interface().(game.CardValue)
		card.KnownValue = true
		card.Value = value
		card.ProbabilityValues[value] = 1.0
		return true
	}
	return false
}

func (ai *AI) setProbabilityColors(card *game.Card, cardsCount map[game.HashValue]int, value game.CardValue, delta int) bool {
	count := 0
	for color, _ := range card.ProbabilityColors {
		count += cardsCount[game.HashColorValue(color, value)]
	}

	count += cardsCount[game.HashColorValue(game.NoneColor, value)] + delta
	correction := float64(count) / float64(count-cardsCount[game.HashColorValue(game.NoneColor, value)]-delta)

	for color, _ := range card.ProbabilityColors {
		if cardsCount[game.HashColorValue(color, value)] == 0 {
			delete(card.ProbabilityColors, color)
		}
		card.ProbabilityColors[color] = correction * float64(cardsCount[game.HashColorValue(color, value)]) / float64(count)
	}

	if len(card.ProbabilityColors) == 1 {
		keys := reflect.ValueOf(card.ProbabilityColors).MapKeys()
		color := keys[0].Interface().(game.CardColor)
		card.KnownColor = true
		card.Color = color
		card.ProbabilityColors[color] = 1.0
		return true
	}
	return false
}

func (ai *AI) setProbabilities() {
	info := &ai.PlayerInfo
	copyCardsCount := map[game.HashValue]int{}
	for _, color := range append(game.Colors, game.NoneColor) {
		val := 1
		if color == game.NoneColor {
			val = 0
		}
		copyCardsCount[game.HashColorValue(color, game.One)] = 3 * val
		copyCardsCount[game.HashColorValue(color, game.Two)] = 2 * val
		copyCardsCount[game.HashColorValue(color, game.Three)] = 2 * val
		copyCardsCount[game.HashColorValue(color, game.Four)] = 2 * val
		copyCardsCount[game.HashColorValue(color, game.Five)] = 1 * val
		copyCardsCount[game.HashColorValue(color, game.NoneValue)] = 0
	}

	for pos, cards := range info.PlayerCards {
		if pos == info.Position {
			continue
		}
		for idx, _ := range cards {
			card := &cards[idx]
			card.ProbabilityColors[card.Color] = 1.0
			card.ProbabilityValues[card.Value] = 1.0
			card.ProbabilityCard[game.HashColorValue(card.Color, card.Value)] = 1.0
		}
	}

	for pos, cards := range info.PlayerCards {
		if pos == info.Position {
			continue
		}
		for idx, _ := range cards {
			card := &cards[idx]
			if card.KnownColor && card.KnownValue {
				copyCardsCount[game.HashColorValue(card.Color, card.Value)]--
			} else if card.KnownColor {
				copyCardsCount[game.HashColorValue(card.Color, game.NoneValue)]--
			} else if card.KnownValue {
				copyCardsCount[game.HashColorValue(game.NoneColor, card.Value)]--
			}
		}
	}

	for _, card := range info.UsedCards {
		copyCardsCount[game.HashColorValue(card.Color, card.Value)]--
	}

	needUpdate := true
	for needUpdate {
		needUpdate = false
		cardsCount := map[game.HashValue]int{}
		for k, v := range copyCardsCount {
			cardsCount[k] = v
		}
		for idx, _ := range info.PlayerCards[info.Position] {
			card := &info.PlayerCards[info.Position][idx]
			if card.KnownColor && card.KnownValue {
				cardsCount[game.HashColorValue(card.Color, card.Value)]--
			} else if card.KnownColor {
				cardsCount[game.HashColorValue(card.Color, game.NoneValue)]--
			} else if card.KnownValue {
				cardsCount[game.HashColorValue(game.NoneColor, card.Value)]--
			}
		}

		for idx, _ := range info.PlayerCards[info.Position] {
			card := &info.PlayerCards[info.Position][idx]
			if card.KnownColor && card.KnownValue {
				card.ProbabilityColors[card.Color] = 1.0
				card.ProbabilityValues[card.Value] = 1.0
				card.ProbabilityCard[game.HashColorValue(card.Color, card.Value)] = 1.0
				continue
			} else if card.KnownValue {
				card.ProbabilityValues[card.Value] = 1.0
				needUpdate = needUpdate || ai.setProbabilityColors(card, cardsCount, card.Value, 1)
			} else if card.KnownColor {
				card.ProbabilityColors[card.Color] = 1.0
				needUpdate = needUpdate || ai.setProbabilityValues(card, cardsCount, card.Color, 1)
			} else {
				for _, color := range game.Colors {
					needUpdate = needUpdate || ai.setProbabilityValues(card, cardsCount, color, 0)
				}

				for _, value := range game.Values {
					needUpdate = needUpdate || ai.setProbabilityColors(card, cardsCount, value, 0)
				}
			}

			for color, _ := range card.ProbabilityColors {
				for value, _ := range card.ProbabilityValues {
					card.ProbabilityCard[game.HashColorValue(color, value)] = card.ProbabilityColors[color] * card.ProbabilityValues[value]
				}
			}
		}
	}

	return
}

func (ai *AI) setAvailableInfomation() {
	info := &ai.PlayerInfo
	for idx, _ := range info.PlayerCards[info.Position] {
		card := &info.PlayerCards[info.Position][idx]
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

	ai.setProbabilities()
}
