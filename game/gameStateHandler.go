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
	MaxStep         int                `json:"max_step"`
}

type Pair struct {
	Count int
	Index int
}

type Pairs []Pair

func (pairs Pairs) Len() int {
	return len(pairs)
}

func (pairs Pairs) Less(i, j int) bool {
	return pairs[i].Count > pairs[j].Count
}

func (pairs Pairs) Swap(i, j int) {
	pairs[i], pairs[j] = pairs[j], pairs[i]
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

func (state *GameState) GetCardCount() int {
	return state.GetCardCountByPlayerCount(len(state.PlayerStates))
}

func (state *GameState) GetCardCountByPlayerCount(count int) int {
	if count >= 4 {
		return 4
	}
	return 5
}

func NewGameState(ids []int, pcards []*Card) GameState {
	state := GameState{
		CurrentPosition: 0,
		BlueTokens:      MaxBlueTokens,
		RedTokens:       MaxRedTokens,
		Step:            0,
		Round:           0,
		UsedCards:       []Card{},
		MaxStep:         0,
	}

	state.TableCards = map[CardColor]Card{
		Red:    *NewCard(Red, NoneValue, true),
		Blue:   *NewCard(Blue, NoneValue, true),
		Green:  *NewCard(Green, NoneValue, true),
		Yellow: *NewCard(Yellow, NoneValue, true),
		Orange: *NewCard(Orange, NoneValue, true),
	}

	cardCount := state.GetCardCountByPlayerCount(len(ids))
	allPlayerPCards := []*[]Card{}
	for i := 0; i < len(ids); i++ {
		userCards := pcards[0:cardCount]
		cards := DereferenceCard(userCards)
		pcards = pcards[cardCount:]
		allPlayerPCards = append(allPlayerPCards, &cards)
	}
	allPlayerCards := SetMostColourfulPlayerCardsAtZeroPlace(allPlayerPCards)
	for i := 0; i < len(ids); i++ {
		state.PlayerStates = append(state.PlayerStates, *NewPlayerState(allPlayerCards, i, ids[i]))
	}

	state.Deck = DereferenceCard(pcards)
	return state
}

func (state *GameState) Copy() GameState {
	newState := GameState{
		CurrentPosition: state.CurrentPosition,
		BlueTokens:      state.BlueTokens,
		RedTokens:       state.RedTokens,
		Step:            state.Step,
		Round:           state.Round,
		MaxStep:         state.MaxStep,
	}

	newState.TableCards = map[CardColor]Card{
		Red:    state.TableCards[Red].Copy(),
		Blue:   state.TableCards[Blue].Copy(),
		Green:  state.TableCards[Green].Copy(),
		Yellow: state.TableCards[Yellow].Copy(),
		Orange: state.TableCards[Orange].Copy(),
	}

	for i := 0; i < len(state.UsedCards); i++ {
		newState.UsedCards = append(newState.UsedCards, state.UsedCards[i].Copy())
	}

	for i := 0; i < len(state.PlayerStates); i++ {
		newState.PlayerStates = append(newState.PlayerStates, state.PlayerStates[i].Copy())
	}

	for i := 0; i < len(state.Deck); i++ {
		newState.Deck = append(newState.Deck, state.Deck[i].Copy())
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
	if state.RedTokens == 0 || state.MaxStep != 0 && state.Step >= state.MaxStep {
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
