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
	Gold
	Black
)

type CardValue int

const (
	NoneValue = iota
	One
	Two
	Three
	Four
	Five
)

func RandomCardsPermutation(cards []*Card) {
	for i := len(cards) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		card := cards[i]
		cards[i] = cards[j]
		cards[j] = card
	}
}

func (this *Card) SetKnown(known bool) {
	this.KnownColor = known
	this.KnownValue = known
}

func (this *Card) GetColors() map[CardColor]string {
	colors := map[CardColor]string{}
	colors[NoneColor] = ""
	colors[Red] = "Red"
	colors[Blue] = "Blue"
	colors[Green] = "Green"
	colors[Gold] = "Gold"
	colors[Black] = "Black"
	return colors
}

func (this *Card) GetValues() map[CardValue]string {
	values := map[CardValue]string{}
	values[NoneValue] = ""
	values[One] = "1"
	values[Two] = "2"
	values[Three] = "3"
	values[Four] = "4"
	values[Five] = "5"
	return values
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
