package ai

import (
	"fmt"
	"math"
	"reflect"
	"time"

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

type Variants struct {
	PCards []*Card
	Color  game.CardColor
	Value  game.CardValue
}

func (ai *BaseAI) SetProbabilities_ConvergenceOfProbability(cards Cards, variants []Variants) {
	info := &ai.PlayerInfo
	delta := 0.0001
	needUpdate := true
	k := cards.Len()

	matrix := make([][]float64, k)
	for i := 0; i < k; i++ {
		matrix[i] = make([]float64, k)
		for j := 0; j < k; j++ {
			matrix[i][j] = cards[j].Probs[i]
		}
	}

	for i := 0; i < k; i++ {
		ok_AllProbs := false
		ok_AllCards := false
		for j := 0; j < k; j++ {
			ok_AllProbs = ok_AllProbs || cards[j].Probs[i] > 0
			ok_AllCards = ok_AllCards || cards[i].Probs[j] > 0
		}
		if !ok_AllCards || !ok_AllProbs {
			for i := 0; i < k; i++ {
				fmt.Print("[]float64{")
				for j := 0; j < k; j++ {
					if j > 0 {
						fmt.Print(", ")
					}
					fmt.Print(matrix[i][j])
				}
				fmt.Println("},")
			}
			panic("Bad data")
		}
	}

	for i := 0; i < k; i++ {
		card := cards[i].Card
		for color, _ := range card.ProbabilityColors {
			card.ProbabilityColors[color] = 0.0
		}

		for value, _ := range card.ProbabilityValues {
			card.ProbabilityValues[value] = 0.0
		}
	}

	now := time.Now().UTC().UnixNano()
	for needUpdate {
		if (time.Now().UTC().UnixNano()-now)/1000000000 > 10 {
			fmt.Println("Debug data:")
			for i := 0; i < k; i++ {
				fmt.Print("[]float64{")
				for j := 0; j < k; j++ {
					if j > 0 {
						fmt.Print(", ")
					}
					fmt.Print(matrix[i][j])
				}
				fmt.Println("},")
			}
			fmt.Println()
			panic("Looping")
		}
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

	deltaVerifyProbs := 0.001
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			c := cards[i]
			if c.Probs[j] < deltaVerifyProbs {
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
			if card.ProbabilityColors[color] < deltaVerifyProbs {
				delete(card.ProbabilityColors, color)
			}
		}

		if !card.KnownColor && len(card.ProbabilityColors) == 1 {
			keys := reflect.ValueOf(card.ProbabilityColors).MapKeys()
			color := keys[0].Interface().(game.CardColor)
			card.KnownColor = true
			card.Color = color
			if card.KnownValue {
				info.VariantsCount[game.ColorValue{Color: card.Color, Value: card.Value}]--
			}
		}

		for value, _ := range card.ProbabilityValues {
			if card.ProbabilityValues[value] < deltaVerifyProbs {
				delete(card.ProbabilityValues, value)
			}
		}

		if !card.KnownValue && len(card.ProbabilityValues) == 1 {
			keys := reflect.ValueOf(card.ProbabilityValues).MapKeys()
			value := keys[0].Interface().(game.CardValue)
			card.KnownValue = true
			card.Value = value
			if card.KnownColor {
				info.VariantsCount[game.ColorValue{Color: card.Color, Value: card.Value}]--
			}
		}
	}
}

func (ai *BaseAI) setCard(card *game.Card, cardVariants []Variants) *Card {
	newCardRef := &Card{
		Card:  card,
		Probs: make([]float64, len(cardVariants)),
	}
	card.ProbabilityCard = map[game.HashValue]float64{}
	count := 0.0
	for color, _ := range card.ProbabilityColors {
		card.ProbabilityColors[color] = 0.0
	}

	for value, _ := range card.ProbabilityValues {
		card.ProbabilityValues[value] = 0.0
	}

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
	info.VariantsCount = map[game.ColorValue]int{}
	variantsCount := info.VariantsCount
	for _, color := range game.ColorsTable {
		variantsCount[game.ColorValue{Color: color, Value: game.One}] = 3
		variantsCount[game.ColorValue{Color: color, Value: game.Two}] = 2
		variantsCount[game.ColorValue{Color: color, Value: game.Three}] = 2
		variantsCount[game.ColorValue{Color: color, Value: game.Four}] = 2
		variantsCount[game.ColorValue{Color: color, Value: game.Five}] = 1
	}

	knownCards := 0
	for _, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if !card.KnownColor || !card.KnownValue {
				continue
			}
			knownCards++
			if card.Color == game.NoneColor {
				panic("NoneColor")
			}
			if card.Value == game.NoneValue {
				panic("NoneValue")
			}

			card.ProbabilityColors = map[game.CardColor]float64{
				card.Color: 1.0,
			}

			card.ProbabilityValues = map[game.CardValue]float64{
				card.Value: 1.0,
			}

			card.ProbabilityCard = map[game.HashValue]float64{
				game.HashColorValue(card.Color, card.Value): 1.0,
			}
			variantsCount[game.ColorValue{Color: card.Color, Value: card.Value}]--
		}
	}

	for idx, _ := range info.UsedCards {
		card := &info.UsedCards[idx]
		card.ProbabilityColors = map[game.CardColor]float64{
			card.Color: 1.0,
		}

		card.ProbabilityValues = map[game.CardValue]float64{
			card.Value: 1.0,
		}

		card.ProbabilityCard = map[game.HashValue]float64{
			game.HashColorValue(card.Color, card.Value): 1.0,
		}
		variantsCount[game.ColorValue{Color: card.Color, Value: card.Value}]--
	}

	for color, card := range info.TableCards {
		for i := 1; i <= int(card.Value); i++ {
			variantsCount[game.ColorValue{Color: color, Value: game.CardValue(i)}]--
			card.ProbabilityColors = map[game.CardColor]float64{
				color: 1.0,
			}

			card.ProbabilityValues = map[game.CardValue]float64{
				card.Value: 1.0,
			}

			card.ProbabilityCard = map[game.HashValue]float64{
				game.HashColorValue(color, card.Value): 1.0,
			}
			info.TableCards[color] = card
		}
	}

	if ai.Type == Type_AICheater {
		for idx, _ := range info.Deck {
			card := &info.Deck[idx]
			card.ProbabilityColors = map[game.CardColor]float64{
				card.Color: 1.0,
			}

			card.ProbabilityValues = map[game.CardValue]float64{
				card.Value: 1.0,
			}

			card.ProbabilityCard = map[game.HashValue]float64{
				game.HashColorValue(card.Color, card.Value): 1.0,
			}
			variantsCount[game.ColorValue{Color: card.Color, Value: card.Value}]--
		}
		info.VariantsCount = map[game.ColorValue]int{}
		return
	}

	cardsRef := Cards{}
	cardVariants := []Variants{}
	for colorValue, count := range variantsCount {
		for i := 0; i < count; i++ {
			cardVariants = append(cardVariants, Variants{
				Color: colorValue.Color,
				Value: colorValue.Value,
			})
		}
	}

	isPreview := false
	lastPos := -1
	for pos, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if card.KnownColor && card.KnownValue {
				continue
			}

			if info.CurrentPostion != pos && lastPos != pos {
				if isPreview {
					panic("I don't know cards of two other players")
				}
				lastPos = pos
				isPreview = true
			}
			cardsRef = append(cardsRef, *ai.setCard(card, cardVariants))
		}
	}

	for idx, _ := range info.Deck {
		card := &info.Deck[idx]
		cardsRef = append(cardsRef, *ai.setCard(card, cardVariants))
	}

	info.VariantsCount = variantsCount
	ai.SetProbabilities_ConvergenceOfProbability(cardsRef, cardVariants)
}

func (ai *BaseAI) setAvailableInformation() {
	if ai.InfoIsSetted {
		return
	}

	info := &ai.PlayerInfo
	for pos, _ := range info.PlayerCards {
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
	}

	ai.setProbabilities()
	ai.InfoIsSetted = true
}
