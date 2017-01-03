package engine

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

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

type ColorInfo struct {
	Color CardColor `json:"color"`
	Count int       `json:"count"`
}

type ValueInfo struct {
	Value CardColor `json:"value"`
	Count int       `json:"count"`
}

type Card struct {
	Color      CardColor `json:"color"`
	KnownColor bool      `json:"known_color"`
	Value      CardValue `json:"value"`
	KnownValue bool      `json:"known_value"`
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

func NewCard(color CardColor, value CardValue, known bool) *Card {
	return &Card{color, known, value, known}
}

type Information struct {
	PlayerId       int         `json:"player_id"`
	PlayerPosition int         `json:"pos"`
	PlayersCards   [][]Card    `json:"players_cards"`
	ColorInfo      []ColorInfo `json:"color_info"`
	ValueInfo      []ValueInfo `json:"value_info"`
}

func NewInformation(cards [][]Card, playerPosition int, playerId int) Information {
	this := new(Information)
	this.PlayerPosition = playerPosition
	this.ColorInfo = []ColorInfo{}
	this.ValueInfo = []ValueInfo{}
	this.PlayerId = playerId

	copyCards := [][]Card{}
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
	Deck            []Card        `json:"deck"`
	Round           int           `json:"round"`
	Step            int           `json:"step"`
	BlueTokens      int           `json:"blue_tokens"`
	RedTokens       int           `json:"red_tokens"`
	CurrentPosition int           `json:"current_pos"`
	UsedCards       []Card        `json:"used_cards"`
	TableCards      map[int]Card  `json:"table_cards"`
	PlayerStates    []Information `json:"information"`
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
	this.TableCards = map[int]Card{
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
		this.PlayerStates = append(this.PlayerStates, NewInformation(allPlayerCards, i, ids[i]))
	}
	this.Deck = DereferenceCard(pcards)
	return *this
}

type Game struct {
	PlayerCount int         `json:"player_count"`
	GameStatus  []GameState `json:"states"`
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
	colors := []CardColor{Red, Blue, Green, Gold, Black}
	for i := 0; i < len(colors); i++ {
		for j := 0; j < len(values); j++ {
			cards = append(cards, NewCard(colors[i], values[j], false))
		}
	}
	RandomCardsPermutation(cards)
	ids = RandomIntPermutation(ids)
	this.PlayerCount = len(ids)
	this.GameStatus = append(this.GameStatus, NewGameState(ids, cards, this.PlayerCount))
	return *this
}

type PlayerGameInfo struct {
	PlayerCount  int          `json:"player_count"`
	Position     int          `json:"pos"`
	Step         int          `json:"step"`
	Round        int          `json:"round"`
	PlayerId     int          `json:"player_id"`
	DeckSize     int          `json:"deck_size"`
	UsedCards    []Card       `json:"used_cards"`
	TableCards   map[int]Card `json:"table_cards"`
	PlayersCards [][]Card     `json:"players_cards"`
	BlueTokens   int          `json:"blue_tokens"`
	RedTokens    int          `json:"red_tokens"`
}

func (this *PlayerGameInfo) Red() Card {
	return this.TableCards[Red]
}

func (this *Game) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	state := this.GameStatus[len(this.GameStatus)-1]
	var playerState Information
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == playerId {
			playerState = state.PlayerStates[i]
		}
	}

	for i := 0; i < len(playerState.PlayersCards); i++ {
		if playerState.PlayerPosition == i {
			cards := &playerState.PlayersCards[i]
			for j := 0; j < len(*cards); j++ {
				card := &(*cards)[j]
				card.KnownColor = false
				card.Color = NoneColor
				card.KnownValue = false
				card.Value = NoneValue
			}
		}
	}

	return PlayerGameInfo{
		PlayerCount:  this.PlayerCount,
		Position:     playerState.PlayerPosition,
		Step:         state.Step,
		Round:        state.Round,
		PlayerId:     playerState.PlayerId,
		DeckSize:     len(state.Deck),
		UsedCards:    state.UsedCards,
		TableCards:   state.TableCards,
		PlayersCards: playerState.PlayersCards,
		BlueTokens:   state.BlueTokens,
		RedTokens:    state.RedTokens,
	}
}

func (this *Game) SprintGame() string {
	b, err := json.Marshal(this)
	if err != nil {
		return ""
	}
	return fmt.Sprintln(string(b))
}

func GetCardValue(value CardValue) string {
	return map[CardValue]string{
		NoneValue: "Unknown Value",
		One:       "1",
		Two:       "2",
		Three:     "3",
		Four:      "4",
		Five:      "5",
	}[value]
}

func GetCardColor(color CardColor) string {
	return map[CardColor]string{
		NoneColor: "Unknown Color",
		Red:       "Red",
		Blue:      "Blue",
		Green:     "Green",
		Gold:      "Gold",
		Black:     "Black",
	}[color]
}
