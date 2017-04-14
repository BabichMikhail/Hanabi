package game

import (
	"math/rand"
)

type HashValue int

func HashColorValue(color CardColor, value CardValue) HashValue {
	return HashValue(int(color) + 10*int(value))
}

func ColorValueByHashColorValue(colorValue HashValue) (CardColor, CardValue) {
	val := int(colorValue)
	return CardColor(val % 10), CardValue(val / 10)
}

type ColorValue struct {
	Color CardColor
	Value CardValue
}

type Card struct {
	Color             CardColor             `json:"color"`
	KnownColor        bool                  `json:"known_color"`
	ProbabilityColors map[CardColor]float64 `json:"probability_colors"`
	Value             CardValue             `json:"value"`
	KnownValue        bool                  `json:"known_value"`
	ProbabilityValues map[CardValue]float64 `json:"probability_values"`
	ProbabilityCard   map[HashValue]float64 `json:"probability_card"`
}

type CardColor int

const (
	NoneColor = iota
	Red
	Blue
	Green
	Yellow
	Orange
)

var Colors = []CardColor{
	NoneColor,
	Red,
	Blue,
	Green,
	Yellow,
	Orange,
}

var ColorsTable = []CardColor{
	Red,
	Blue,
	Green,
	Yellow,
	Orange,
}

type CardValue int

const (
	NoneValue = iota
	One
	Two
	Three
	Four
	Five
)

var Values = []CardValue{
	NoneValue,
	One,
	Two,
	Three,
	Four,
	Five,
}

func RandomCardsPermutation(cards []*Card) {
	for i := len(cards) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func (card *Card) SetKnown(known bool) {
	card.KnownColor = known
	card.KnownValue = known
}

func (card *Card) GetColors() map[CardColor]string {
	return map[CardColor]string{
		NoneColor: "",
		Red:       "Red",
		Blue:      "Blue",
		Green:     "Green",
		Yellow:    "Yellow",
		Orange:    "Orange",
	}
}

func (card *Card) GetValues() map[CardValue]string {
	return map[CardValue]string{
		NoneValue: "",
		One:       "1",
		Two:       "2",
		Three:     "3",
		Four:      "4",
		Five:      "5",
	}
}

func (card *Card) GetPoints() int {
	return map[CardValue]int{
		NoneValue: 0,
		One:       1,
		Two:       2,
		Three:     3,
		Four:      4,
		Five:      5,
	}[card.Value]
}

func GetTableColorOrder() map[string]CardColor {
	return map[string]CardColor{
		"red":    Red,
		"blue":   Blue,
		"green":  Green,
		"yellow": Yellow,
		"orange": Orange,
	}
}

func DereferenceCard(pcards []*Card) []Card {
	cards := []Card{}
	for i := 0; i < len(pcards); i++ {
		cards = append(cards, *pcards[i])
	}
	return cards
}

func (card *Card) UpdateProbability() {
	if !card.KnownColor || !card.KnownValue {
		return
	}
	card.ProbabilityValues = map[CardValue]float64{
		card.Value: 1.0,
	}
	card.ProbabilityColors = map[CardColor]float64{
		card.Color: 1.0,
	}
}

func NewCard(color CardColor, value CardValue, known bool) *Card {
	values := map[CardValue]bool{
		One:   !known || value == One,
		Two:   !known || value == Two,
		Three: !known || value == Three,
		Four:  !known || value == Four,
		Five:  !known || value == Five,
	}
	probValues := map[CardValue]float64{}
	for value, isAvailable := range values {
		if isAvailable {
			probValues[value] = 0.0
		}
	}

	colors := map[CardColor]bool{
		Red:    !known || color == Red,
		Blue:   !known || color == Blue,
		Green:  !known || color == Green,
		Yellow: !known || color == Yellow,
		Orange: !known || color == Orange,
	}
	probColors := map[CardColor]float64{}
	for color, isAvailable := range colors {
		if isAvailable {
			probColors[color] = 0.0
		}
	}
	return &Card{color, known, probColors, value, known, probValues, map[HashValue]float64{}}
}

func (card Card) Copy() Card {
	probValues := map[CardValue]float64{}
	for k, v := range card.ProbabilityValues {
		probValues[k] = v
	}

	probColors := map[CardColor]float64{}
	for k, v := range card.ProbabilityColors {
		probColors[k] = v
	}

	probCards := map[HashValue]float64{}
	for k, v := range card.ProbabilityCard {
		probCards[k] = v
	}

	return Card{
		Color:             card.Color,
		KnownColor:        card.KnownColor,
		ProbabilityColors: probColors,
		Value:             card.Value,
		KnownValue:        card.KnownValue,
		ProbabilityValues: probValues,
		ProbabilityCard:   probCards,
	}
}
