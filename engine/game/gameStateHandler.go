package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

const (
	MaxBlueTokens = 8
	MaxRedTokens  = 3
)

type GameState struct {
	Deck            []Card             `json:"deck"`
	Round           int                `json:"round"`
	Step            int                `json:"step"`
	BlueTokens      int                `json:"blue_tokens"`
	RedTokens       int                `json:"red_tokens"`
	CurrentPosition int                `json:"current_pos"`
	UsedCards       []Card             `json:"used_cards"`
	TableCards      map[CardColor]Card `json:"table_cards"`
	PlayerStates    []PlayerState      `json:"player_state"`
}

type Pair struct {
	Count int
	Index int
}

type Pairs []Pair

func (this Pairs) Len() int {
	return len(this)
}

func (this Pairs) Less(i, j int) bool {
	return this[i].Count > this[j].Count
}

func (this Pairs) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func SetMostColourfulPlayerCardsAtZeroPlace(pcards []*[]Card) [][]Card {
	pairs := Pairs{}
	for i := 0; i < len(pcards); i++ {
		colors := map[CardColor]int{}
		for j := 0; j < len(*pcards[i]); j++ {
			colors[(*pcards[i])[j].Color] = 1
		}
		pairs = append(pairs, Pair{len(colors), i})
	}
	sort.Sort(pairs)
	cards := [][]Card{}
	for i := 0; i < pairs.Len(); i++ {
		cards = append(cards, *pcards[pairs[i].Index])
	}
	return cards
}

func (this *GameState) GetCardCount() int {
	if len(this.PlayerStates) >= 4 {
		return 4
	}
	return 5
}

func NewGameState(ids []int, pcards []*Card) GameState {
	this := GameState{
		CurrentPosition: 0,
		BlueTokens:      MaxBlueTokens,
		RedTokens:       MaxRedTokens,
		Step:            0,
		Round:           0,
		UsedCards:       []Card{},
	}

	this.TableCards = map[CardColor]Card{
		Red:    *NewCard(Red, NoneValue, true),
		Blue:   *NewCard(Blue, NoneValue, true),
		Green:  *NewCard(Green, NoneValue, true),
		Yellow: *NewCard(Yellow, NoneValue, true),
		Orange: *NewCard(Orange, NoneValue, true),
	}

	cardCount := this.GetCardCount()
	allPlayerPCards := []*[]Card{}
	for i := 0; i < len(ids); i++ {
		userCards := pcards[0:cardCount]
		pcards = append(pcards[:0], pcards[cardCount:]...)
		cards := DereferenceCard(userCards)
		allPlayerPCards = append(allPlayerPCards, &cards)
	}
	allPlayerCards := SetMostColourfulPlayerCardsAtZeroPlace(allPlayerPCards)
	for i := 0; i < len(ids); i++ {
		this.PlayerStates = append(this.PlayerStates, NewPlayerState(allPlayerCards, i, ids[i]))
	}

	this.Deck = DereferenceCard(pcards)
	return this
}

func (this *GameState) Copy() GameState {
	newState := GameState{
		CurrentPosition: this.CurrentPosition,
		BlueTokens:      this.BlueTokens,
		RedTokens:       this.RedTokens,
		Step:            this.Step,
		Round:           this.Round,
	}

	newState.TableCards = map[CardColor]Card{
		Red:    this.TableCards[Red].Copy(),
		Blue:   this.TableCards[Blue].Copy(),
		Green:  this.TableCards[Green].Copy(),
		Yellow: this.TableCards[Yellow].Copy(),
		Orange: this.TableCards[Orange].Copy(),
	}

	for i := 0; i < len(this.UsedCards); i++ {
		newState.UsedCards = append(newState.UsedCards, this.UsedCards[i].Copy())
	}

	for i := 0; i < len(this.PlayerStates); i++ {
		newState.PlayerStates = append(newState.PlayerStates, this.PlayerStates[i].Copy())
	}

	for i := 0; i < len(this.Deck); i++ {
		newState.Deck = append(newState.Deck, this.Deck[i].Copy())
	}

	return newState
}

func (state *GameState) GetPlayerPositionById(id int) (pos int, err error) {
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == id {
			return i, nil
		}
	}
	return -1, errors.New("Player not found")
}

func (state *GameState) IsGameOver() bool {
	if state.RedTokens == 0 {
		return true
	}

	cardInHands := 0
	for i := 0; i < len(state.PlayerStates); i++ {
		cardInHands += len(state.PlayerStates[i].PlayerCards)
	}
	if cardInHands == len(state.PlayerStates)*(state.GetCardCount()-1) {
		return true
	}

	for _, card := range state.TableCards {
		if card.Value != Five {
			return false
		}
	}
	return true
}

func (state *GameState) GetPoints() (points int, err error) {
	for _, card := range state.TableCards {
		points += card.GetPoints()
	}
	return
}

func (state *GameState) Sprint() string {
	b, err := json.Marshal(state)
	if err != nil {
		return ""
	}
	return fmt.Sprint(string(b))
}
