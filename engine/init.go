package engine

import "math/rand"

type Color int

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

type ColorInfo struct {
	Color Color `json:"color"`
	Count int   `json:"count"`
}

type ValueInfo struct {
	Value Color `json:"value"`
	Count int   `json:"count"`
}

type Card struct {
	Color      Color     `json:"color"`
	KnownColor bool      `json:"known_color"`
	Value      CardValue `json:"value"`
	KnownValue bool      `json:"known_value"`
}

func (this *Card) SetKnown(known bool) {
	this.KnownColor = known
	this.KnownValue = known
}

func NewCard(color Color, value CardValue, known bool) *Card {
	return &Card{color, known, value, known}
}

type Information struct {
	PlayerId       int         `json:"player_id"`
	PlayerPosition int         `json:"pos"`
	PlayersCards   [][]Card    `json:"players_cards"`
	ColorInfo      []ColorInfo `json:"color_info"`
	ValueInfo      []ValueInfo `json:"value_info"`
}

func NewInformation(cards [][]Card, playerPosition int) Information {
	this := new(Information)
	this.PlayerPosition = playerPosition
	this.ColorInfo = []ColorInfo{}
	this.ValueInfo = []ValueInfo{}

	copyCards := make([][]Card, len(cards))
	for i := 0; i < len(cards); i++ {
		copyPlayerCards := make([]Card, len(cards[i]))
		copy(copyPlayerCards, cards[i])
		for j := 0; j < len(cards[playerPosition]); j++ {
			copyPlayerCards[j].SetKnown(i != playerPosition)
		}
		copyCards = append(copyCards, copyPlayerCards)
	}
	this.PlayersCards = copyCards
	return *this
}

type Deck struct {
	Cards []Card `json:"cards"`
}

type GameState struct {
	Deck            []Card         `json:"deck"`
	Round           int            `json:"round"`
	Step            int            `json:"step"`
	BlueTokens      int            `json:"blue_tokens"`
	RedTokens       int            `json:"red_tokens"`
	CurrentPosition int            `json:"current_pos"`
	UsedCards       []Card         `json:"used_cards"`
	TableCards      map[Color]Card `json:"table_cards"`
	PlayerStates    []Information  `json:"information"`
}

func DereferenceCard(pcards []*Card) []Card {
	cards := []Card{}
	for i := 0; i < len(pcards); i++ {
		cards = append(cards, *pcards[i])
	}
	return cards
}

func NewGameState(ids []int, pcards []*Card, playerCount int) GameState {
	this := new(GameState)
	this.CurrentPosition = 0
	this.BlueTokens = 8
	this.RedTokens = 3
	this.TableCards = map[Color]Card{
		Red:   *NewCard(Red, NoneValue, true),
		Blue:  *NewCard(Blue, NoneValue, true),
		Green: *NewCard(Green, NoneValue, true),
		Gold:  *NewCard(Gold, NoneValue, true),
		Black: *NewCard(Black, NoneValue, true),
	}
	this.Step = 0
	this.Round = 0
	cardCount := 5
	if playerCount >= 4 {
		cardCount = 4
	}
	this.UsedCards = []Card{}
	allPlayerCards := [][]Card{}

	for i := 0; i < len(ids); i++ {
		userCards := pcards[0:cardCount]
		pcards = append(pcards[:0], pcards[cardCount:]...)
		allPlayerCards = append(allPlayerCards, DereferenceCard(userCards))
	}
	for i := 0; i < len(ids); i++ {
		this.PlayerStates = append(this.PlayerStates, NewInformation(allPlayerCards, i))
	}
	this.Deck = DereferenceCard(pcards)
	return *this
}

type Game struct {
	PlayerCount int         `json:"player_count"`
	GameStates  []GameState `json:"states"`
}

func RandomCardsPermutation(cards []*Card) {
	for i := len(cards) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		card := cards[i]
		cards[i] = cards[j]
		cards[j] = card
	}
}

func RandomIntPermutation(values []int) []int {
	for i := len(values) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		value := values[i]
		values[i] = values[j]
		values[j] = value
	}
	return values
}

func NewGame(ids []int) Game {
	this := new(Game)
	cards := []*Card{}
	values := []CardValue{One, One, One, Two, Two, Three, Three, Four, Four, Five}
	colors := []Color{Red, Blue, Green, Gold, Black}
	for i := 0; i < len(colors); i++ {
		for j := 0; j < len(values); j++ {
			cards = append(cards, NewCard(colors[i], values[j], false))
		}
	}
	RandomCardsPermutation(cards)
	ids = RandomIntPermutation(ids)
	this.PlayerCount = len(ids)
	this.GameStates = append(this.GameStates, NewGameState(ids, cards, this.PlayerCount))
	return *this
}
