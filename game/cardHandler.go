package game

import (
	"fmt"
	"math/rand"
)

type HashValue int

func HashColorValue(color CardColor, value CardValue) HashValue {
	if color != NoneColor && value != NoneValue {
		return HashValue(int(color) - 1 + 5*(int(value)-1))
	}
	if color == NoneColor {
		return HashValue(int(color) + int(value) - 11)
	}
	return HashValue(int(color) + int(value) - 6)
}

var CardCodes map[HashValue]ColorValue

func init() {
	CardCodes = map[HashValue]ColorValue{}
	for _, value := range Values {
		for _, color := range Colors {
			fmt.Println(GetCardColor(color), value, HashColorValue(color, value))
			CardCodes[HashColorValue(color, value)] = ColorValue{Color: color, Value: value}
		}
	}
}

func ColorValueByHashColorValue(hashValue HashValue) (CardColor, CardValue) {
	cv := CardCodes[hashValue]
	return cv.Color, cv.Value
}

type ColorValue struct {
	Color CardColor
	Value CardValue
}

func (cv ColorValue) String() string {
	return fmt.Sprintf("{ %s %d }", GetCardColor(cv.Color), cv.Value)
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

var ColorsTable = Colors[1:]
var Colors = []CardColor{
	NoneColor,
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

var ValuesTable = Values[1:]
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

func (card *Card) CheckVisible() {
	if card.Color == NoneColor || card.Value == NoneValue {
		panic("Bad card")
	}
}

func (card *Card) IsCardPlayable(progress map[CardColor]CardValue) bool {
	return progress[card.Color]+1 == card.Value
}

func (card *Card) IsVisible() bool {
	return card.KnownColor && card.KnownValue
}

func (card *Card) SetValue(value CardValue) {
	card.Value = value
	card.KnownValue = true
	card.ProbabilityValues = map[CardValue]float64{value: 1.0}
}

func (card *Card) SetColor(color CardColor) {
	card.Color = color
	card.KnownColor = true
	card.ProbabilityColors = map[CardColor]float64{color: 1.0}
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

func (card Card) String() string {
	color := map[CardColor]string{
		NoneColor: "Unknown",
		Red:       "Red",
		Blue:      "Blue",
		Green:     "Green",
		Yellow:    "Yello",
		Orange:    "Orange",
	}[card.Color]
	value := map[CardValue]string{
		NoneValue: "Unknown",
		One:       "1",
		Two:       "2",
		Three:     "3",
		Four:      "4",
		Five:      "5",
	}[card.Value]
	return fmt.Sprintf("[ %s %s ]", color, value)
}

func (card *Card) NormalizeProbabilities(color CardColor, value CardValue, countLeft int) {
	if card.KnownValue && card.KnownColor {
		return
	}

	colorValue := HashColorValue(color, value)

	if probability, ok := card.ProbabilityCard[colorValue]; ok {
		probSum := 1.0
		if countLeft == 0 {
			probSum -= probability
			delete(card.ProbabilityCard, colorValue)
		} else {
			count := float64(countLeft)
			probSum -= probability / (count + 1)
			card.ProbabilityCard[colorValue] = probability / (count + 1) * count
		}

		for colorValue, _ := range card.ProbabilityCard {
			card.ProbabilityCard[colorValue] /= probSum
		}
	}
	return
}
