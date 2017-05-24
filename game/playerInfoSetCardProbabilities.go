package game

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func (info *PlayerGameInfo) SetEasyKnownAboutCards() {
	for pos, _ := range info.PlayerCards {
		for idx, _ := range info.PlayerCards[pos] {
			card := &info.PlayerCards[pos][idx]
			if card.Color == NoneColor && len(card.ProbabilityColors) == 1 {
				for color, _ := range card.ProbabilityColors {
					card.SetColor(color)
				}
			}

			if card.Value == NoneValue && len(card.ProbabilityValues) == 1 {
				for value, _ := range card.ProbabilityValues {
					card.SetValue(value)
				}
			}
		}
	}
}

type CardProbs struct {
	Card  *Card
	Probs []float64
}

type Cards []CardProbs

func (c Cards) Len() int {
	return len(c)
}

func (c Cards) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Variants struct {
	PCards []*CardProbs
	Color  CardColor
	Value  CardValue
}

func (info *PlayerGameInfo) SetProbabilities_ConvergenceOfProbability(cards Cards, variants []Variants) {
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
			fmt.Println("Debug data:")
			fmt.Println(info.VariantsCount)
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
			card.ProbabilityCard[HashColorValue(v.Color, v.Value)] += c.Probs[j]
		}
	}

	for i := 0; i < k; i++ {
		card := cards[i].Card
		colorProbSum := 1.0
		for color, prob := range card.ProbabilityColors {
			if card.ProbabilityColors[color] < deltaVerifyProbs {
				delete(card.ProbabilityColors, color)
			} else {
				colorProbSum += prob
			}
		}

		for color, prob := range card.ProbabilityColors {
			card.ProbabilityColors[color] = prob / colorProbSum
		}

		if !card.KnownColor && len(card.ProbabilityColors) == 1 {
			keys := reflect.ValueOf(card.ProbabilityColors).MapKeys()
			color := keys[0].Interface().(CardColor)
			card.KnownColor = true
			card.Color = color
			card.ProbabilityColors[color] = 1.0
			if card.KnownValue {
				info.VariantsCount[ColorValue{Color: card.Color, Value: card.Value}]--
			}
		}

		valueProbSum := 1.0
		for value, prob := range card.ProbabilityValues {
			if card.ProbabilityValues[value] < deltaVerifyProbs {
				delete(card.ProbabilityValues, value)
			} else {
				valueProbSum += prob
			}
		}

		for value, prob := range card.ProbabilityValues {
			card.ProbabilityValues[value] = prob / valueProbSum
		}

		if !card.KnownValue && len(card.ProbabilityValues) == 1 {
			keys := reflect.ValueOf(card.ProbabilityValues).MapKeys()
			value := keys[0].Interface().(CardValue)
			card.KnownValue = true
			card.Value = value
			card.ProbabilityValues[value] = 1.0
			if card.KnownColor {
				info.VariantsCount[ColorValue{Color: card.Color, Value: card.Value}]--
			}
		}

		probSum := 0.0
		for key, prob := range card.ProbabilityCard {
			if prob < math.Pow(deltaVerifyProbs, 2) {
				delete(card.ProbabilityCard, key)
			} else {
				probSum += prob
			}
		}

		for colorValue, prob := range card.ProbabilityCard {
			newProb := prob / probSum
			card.ProbabilityCard[colorValue] = prob / probSum
			if newProb == 1.0 {
				card.KnownColor = true
				card.KnownValue = true
				color, value := ColorValueByHashColorValue(colorValue)
				card.Color = color
				card.Value = value
				card.ProbabilityColors = map[CardColor]float64{color: 1.0}
				card.ProbabilityValues = map[CardValue]float64{value: 1.0}
			}
		}
	}
}

func (info *PlayerGameInfo) SetCard(card *Card, cardVariants []Variants) *CardProbs {
	newCardRef := &CardProbs{
		Card:  card,
		Probs: make([]float64, len(cardVariants)),
	}
	card.ProbabilityCard = map[HashValue]float64{}
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

func (info *PlayerGameInfo) SetVariantsCount(isCheater, isFullCheater bool) {
	variantsCount := info.GetDefaultDeck()

	for _, cards := range info.PlayerCards {
		for _, card := range cards {
			if !card.KnownColor || !card.KnownValue {
				continue
			}
			if card.Color == NoneColor {
				panic("NoneColor")
			}
			if card.Value == NoneValue {
				panic("NoneValue")
			}

			card.ProbabilityColors = map[CardColor]float64{card.Color: 1.0}
			card.ProbabilityValues = map[CardValue]float64{card.Value: 1.0}
			card.ProbabilityCard = map[HashValue]float64{
				HashColorValue(card.Color, card.Value): 1.0,
			}
			variantsCount[ColorValue{Color: card.Color, Value: card.Value}]--
		}
	}

	for _, card := range info.UsedCards {
		card.ProbabilityColors = map[CardColor]float64{card.Color: 1.0}
		card.ProbabilityValues = map[CardValue]float64{card.Value: 1.0}
		card.ProbabilityCard = map[HashValue]float64{
			HashColorValue(card.Color, card.Value): 1.0,
		}
		variantsCount[ColorValue{Color: card.Color, Value: card.Value}]--
	}

	for color, card := range info.TableCards {
		for i := 1; i <= int(card.Value); i++ {
			variantsCount[ColorValue{Color: color, Value: CardValue(i)}]--

			card.ProbabilityColors = map[CardColor]float64{color: 1.0}
			card.ProbabilityValues = map[CardValue]float64{card.Value: 1.0}
			card.ProbabilityCard = map[HashValue]float64{
				HashColorValue(color, card.Value): 1.0,
			}
			info.TableCards[color] = card
		}
	}

	if isFullCheater {
		for _, card := range info.Deck {
			card.ProbabilityColors = map[CardColor]float64{card.Color: 1.0}
			card.ProbabilityValues = map[CardValue]float64{card.Value: 1.0}
			card.ProbabilityCard = map[HashValue]float64{
				HashColorValue(card.Color, card.Value): 1.0,
			}
			variantsCount[ColorValue{Color: card.Color, Value: card.Value}]--
		}
		info.VariantsCount = map[ColorValue]int{}
		return
	}
	info.VariantsCount = variantsCount
}

func (info *PlayerGameInfo) GetCardVariants() []Variants {
	cardVariants := []Variants{}
	for colorValue, count := range info.VariantsCount {
		for i := 0; i < count; i++ {
			cardVariants = append(cardVariants, Variants{
				Color: colorValue.Color,
				Value: colorValue.Value,
			})
		}
	}
	return cardVariants
}

func (info *PlayerGameInfo) SetProbabilities(isCheater, isFullCheater bool) {
	if info.InfoIsSetted {
		return
	}
	info.SetEasyKnownAboutCards()
	info.SetVariantsCount(isCheater, isFullCheater)
	cardVariants := info.GetCardVariants()

	isPreview := false
	lastPos := -1
	cardsRef := Cards{}
	for pos, cards := range info.PlayerCards {
		for idx, _ := range cards {
			card := &cards[idx]
			if card.KnownColor && card.KnownValue {
				continue
			}

			if info.CurrentPosition != pos && lastPos != pos {
				if isPreview {
					panic("I don't know cards of two other players")
				}
				lastPos = pos
				isPreview = true
			}
			cardsRef = append(cardsRef, *info.SetCard(card, cardVariants))
		}
	}

	for idx, _ := range info.Deck {
		card := &info.Deck[idx]
		cardsRef = append(cardsRef, *info.SetCard(card, cardVariants))
	}

	info.SetProbabilities_ConvergenceOfProbability(cardsRef, cardVariants)
	info.InfoIsSetted = true
}
