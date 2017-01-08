package game

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Game struct {
	PlayerCount  int       `json:"player_count"`
	InitState    GameState `json:"init_state"`
	CurrentState GameState `json:"current_state"`
	Actions      []Action  `json:"actions"`
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
	state := NewGameState(ids, cards, this.PlayerCount)
	this.InitState = state
	this.CurrentState = state.Copy()
	this.Actions = []Action{}
	return *this
}

func (this *Game) SprintGame() string {
	b, err := json.Marshal(this)
	if err != nil {
		return ""
	}
	return fmt.Sprintln(string(b))
}

func (this *Game) GetPlayerPositionById(id int) (pos int, err error) {
	for i := 0; i < len(this.CurrentState.PlayerStates); i++ {
		if this.CurrentState.PlayerStates[i].PlayerId == id {
			return i, nil
		}
	}
	return -1, errors.New("Player not found")
}

func (this *Game) GetPoints() (points int, err error) {
	if !this.IsGameOver() {
		return 0, errors.New("Game not is over")
	}

	state := &this.CurrentState
	for _, card := range state.TableCards {
		points += card.GetPoints()
	}
	return
}

func (this *Game) IsGameOver() bool {
	state := &this.CurrentState
	if state.RedTokens == 0 {
		return true
	}

	cardInHands := 0
	for i := 0; i < len(state.PlayerStates[0].PlayersCards); i++ {
		cardInHands += len(state.PlayerStates[0].PlayersCards[i])
	}
	if cardInHands == state.PlayerCount*(state.GetCardCount()-1) {
		return true
	}

	for _, card := range state.TableCards {
		if card.Value != Five {
			return false
		}
	}
	return true
}