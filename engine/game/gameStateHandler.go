package game

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

	cardCount := 5
	if playerCount >= 4 {
		cardCount = 4
	}

	allPlayerCards := [][]Card{}
	for i := 0; i < len(ids); i++ {
		userCards := pcards[0:cardCount]
		pcards = append(pcards[:0], pcards[cardCount:]...)
		allPlayerCards = append(allPlayerCards, DereferenceCard(userCards))
	}
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
