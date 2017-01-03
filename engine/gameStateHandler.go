package engine

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

func NewGameState(ids []int, pcards []*Card, playerCount int) GameState {
	this := new(GameState)
	this.CurrentPosition = 0
	this.BlueTokens = 8
	this.RedTokens = 3
	this.Step = 0
	this.Round = 0
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
