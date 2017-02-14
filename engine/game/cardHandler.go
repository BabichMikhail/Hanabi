package game

import "math/rand"

type Card struct {
	Color      CardColor `json:"color"`
	KnownColor bool      `json:"known_color"`
	Value      CardValue `json:"value"`
	KnownValue bool      `json:"known_value"`
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

func (this *Card) SetKnown(known bool) {
	this.KnownColor = known
	this.KnownValue = known
}

func (this *Card) GetColors() map[CardColor]string {
	return map[CardColor]string{
		NoneColor: "",
		Red:       "Red",
		Blue:      "Blue",
		Green:     "Green",
		Yellow:    "Yellow",
		Orange:    "Orange",
	}
}

func (this *Card) GetValues() map[CardValue]string {
	return map[CardValue]string{
		NoneValue: "",
		One:       "1",
		Two:       "2",
		Three:     "3",
		Four:      "4",
		Five:      "5",
	}
}

func (this *Card) GetPoints() int {
	return map[CardValue]int{
		NoneValue: 0,
		One:       1,
		Two:       2,
		Three:     3,
		Four:      4,
		Five:      5,
	}[this.Value]
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

func NewCard(color CardColor, value CardValue, known bool) *Card {
	return &Card{color, known, value, known}
}

func (this Card) Copy() Card {
	return Card{
		Color:      this.Color,
		KnownColor: this.KnownColor,
		Value:      this.Value,
		KnownValue: this.KnownValue,
	}
}
