package ai

import (
	"math"
	"reflect"

	"github.com/BabichMikhail/Hanabi/game"
)

type Card struct {
	Card  *game.Card
	Probs []float64
}

type Cards []Card

func (c Cards) Len() int {
	return len(c)
}

func (c Cards) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CardColorValue struct {
	Color game.CardColor
	Value game.CardValue
}

type Variants struct {
	PCards []*Card
	Color  game.CardColor
	Value  game.CardValue
}

func (ai *BaseAI) SetProbabilities_ConvergenceOfProbability(cards Cards, variants []Variants) {
	delta := 0.0001
	needUpdate := true
	k := cards.Len()

	for i := 0; i < k; i++ {
		card := cards[i].Card
		for color, _ := range card.ProbabilityColors {
			card.ProbabilityColors[color] = 0.0
		}

		for value, _ := range card.ProbabilityValues {
			card.ProbabilityValues[value] = 0.0
		}
	}

	for needUpdate {
		needUpdate = false
		for i := 0; i < k; i++ {
			sum := 0.0
			for j := 0; j < k; j++ {
				sum += cards[i].Probs[j]
			}

			if math.Abs(sum-1.0) < delta {
				continue
			}
			for j := 0; j < k; j++ {
				cards[i].Probs[j] /= sum
			}
		}

		for j := 0; j < k; j++ {
			sum := 0.0
			for i := 0; i < k; i++ {
				sum += cards[i].Probs[j]
			}

			if math.Abs(sum-1.0) < delta {
				continue
			}

			needUpdate = true
			for i := 0; i < k; i++ {
				cards[i].Probs[j] /= sum
			}
		}
	}

	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			c := cards[i]
			if c.Probs[j] < delta {
				continue
			}
			v := variants[j]
			card := c.Card
			card.ProbabilityColors[v.Color] += c.Probs[j]
			card.ProbabilityValues[v.Value] += c.Probs[j]
			card.ProbabilityCard[game.HashColorValue(v.Color, v.Value)] += c.Probs[j]
		}
	}

	for i := 0; i < k; i++ {
		card := cards[i].Card
		for color, _ := range card.ProbabilityColors {
			if card.ProbabilityColors[color] == 0 {
				delete(card.ProbabilityColors, color)
			}
		}

		if !card.KnownColor && len(card.ProbabilityColors) == 1 {
			keys := reflect.ValueOf(card.ProbabilityColors).MapKeys()
			color := keys[0].Interface().(game.CardColor)
			card.KnownColor = true
			card.Color = color
		}

		for value, _ := range card.ProbabilityValues {
			if card.ProbabilityValues[value] == 0 {
				delete(card.ProbabilityValues, value)
			}
		}

		if !card.KnownValue && len(card.ProbabilityValues) == 1 {
			keys := reflect.ValueOf(card.ProbabilityValues).MapKeys()
			value := keys[0].Interface().(game.CardValue)
			card.KnownValue = true
			card.Value = value
		}
	}
}

func (ai *BaseAI) setCard(card *game.Card, cardVariants []Variants) *Card {
	newCardRef := &Card{
		Card:  card,
		Probs: make([]float64, len(cardVariants)),
	}

	count := 0.0
	for k := 0; k < len(cardVariants); k++ {
		_, colorOK := card.ProbabilityColors[cardVariants[k].Color]
		_, valueOK := card.ProbabilityValues[cardVariants[k].Value]
		if colorOK && valueOK {
			count++
			newCardRef.Probs[k] = 1.0
		} else {
			newCardRef.Probs[k] = 0.0
		}
	}

	for k := 0; k < len(cardVariants); k++ {
		if newCardRef.Probs[k] == 1.0 {
			newCardRef.Probs[k] /= count
		}
	}
	return newCardRef
}

func (ai *BaseAI) setProbabilities() {
	info := &ai.PlayerInfo
	cardsCount := map[CardColorValue]int{}
	for _, color := range game.Colors {
		if color == game.NoneColor {
			continue
		}
		cardsCount[CardColorValue{color, game.One}] = 3
		cardsCount[CardColorValue{color, game.Two}] = 2
		cardsCount[CardColorValue{color, game.Three}] = 2
		cardsCount[CardColorValue{color, game.Four}] = 2
		cardsCount[CardColorValue{color, game.Five}] = 1
	}

	for _, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if !card.KnownColor || !card.KnownValue {
				continue
			}

			card.ProbabilityColors = map[game.CardColor]float64{
				card.Color: 1.0,
			}

			card.ProbabilityValues = map[game.CardValue]float64{
				card.Value: 1.0,
			}
			cardsCount[CardColorValue{card.Color, card.Value}]--
		}
	}

	for _, card := range info.UsedCards {
		cardsCount[CardColorValue{card.Color, card.Value}]--
	}

	for _, card := range info.TableCards {
		for i := 1; i <= int(card.Value); i++ {
			cardsCount[CardColorValue{card.Color, game.CardValue(i)}]--
		}
	}

	cardsRef := Cards{}
	cardVariants := []Variants{}
	for colorValue, count := range cardsCount {
		for i := 0; i < count; i++ {
			cardVariants = append(cardVariants, Variants{
				Color: colorValue.Color,
				Value: colorValue.Value,
			})
		}
	}

	for _, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if card.KnownColor && card.KnownValue {
				continue
			}

			cardsRef = append(cardsRef, *ai.setCard(card, cardVariants))
		}
	}

	for idx, _ := range info.Deck {
		card := &info.Deck[idx]
		if card.KnownColor && card.KnownValue {
			continue
		}

		cardsRef = append(cardsRef, *ai.setCard(card, cardVariants))
	}

	ai.SetProbabilities_ConvergenceOfProbability(cardsRef, cardVariants)
	return
}

func (ai *BaseAI) setAvailableInformation() {
	if ai.InfoIsSetted {
		return
	}

	info := &ai.PlayerInfo
	pos := info.CurrentPostion
	for idx, _ := range info.PlayerCards[pos] {
		card := &info.PlayerCards[pos][idx]
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
	ai.InfoIsSetted = true
}
