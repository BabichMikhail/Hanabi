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
