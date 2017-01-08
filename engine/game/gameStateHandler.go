package game

import "sort"

const (
	MaxBlueTokens = 8
	MaxRedTokens  = 3
)

type GameState struct {
	Deck            []Card        `json:"deck"`
	Round           int           `json:"round"`
	PlayerCount     int           `json:"player_count"`
	Step            int           `json:"step"`
	BlueTokens      int           `json:"blue_tokens"`
	RedTokens       int           `json:"red_tokens"`
	CurrentPosition int           `json:"current_pos"`
	UsedCards       []Card        `json:"used_cards"`
	TableCards      map[int]Card  `json:"table_cards"`
	PlayerStates    []PlayerState `json:"player_state"`
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
	if this.PlayerCount >= 4 {
		return 4
	}
	return 5
}

func NewGameState(ids []int, pcards []*Card, playerCount int) GameState {
	this := GameState{
		CurrentPosition: 0,
		BlueTokens:      MaxBlueTokens,
		RedTokens:       MaxRedTokens,
		Step:            0,
		Round:           0,
		PlayerCount:     playerCount,
		UsedCards:       []Card{},
	}

	this.TableCards = map[int]Card{
		Red:   *NewCard(Red, NoneValue, true),
		Blue:  *NewCard(Blue, NoneValue, true),
		Green: *NewCard(Green, NoneValue, true),
		Gold:  *NewCard(Gold, NoneValue, true),
		Black: *NewCard(Black, NoneValue, true),
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

func (this GameState) Copy() GameState {
	newState := GameState{
		CurrentPosition: this.CurrentPosition,
		BlueTokens:      this.BlueTokens,
		RedTokens:       this.RedTokens,
		Step:            this.Step,
		Round:           this.Round,
		PlayerCount:     this.PlayerCount,
	}

	newState.TableCards = map[int]Card{
		Red:   this.TableCards[Red].Copy(),
		Blue:  this.TableCards[Blue].Copy(),
		Green: this.TableCards[Green].Copy(),
		Gold:  this.TableCards[Gold].Copy(),
		Black: this.TableCards[Black].Copy(),
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
