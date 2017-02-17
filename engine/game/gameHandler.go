package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	PlayerCount  int       `json:"player_count"`
	InitState    GameState `json:"init_state"`
	CurrentState GameState `json:"current_state"`
	Actions      []Action  `json:"actions"`
	Seed         int64     `json:"seed"`
}

func NewGame(ids []int) Game {
	this := new(Game)
	this.Seed = time.Now().UTC().UnixNano()
	rand.Seed(this.Seed)
	cards := []*Card{}
	values := []CardValue{One, One, One, Two, Two, Three, Three, Four, Four, Five}
	colors := []CardColor{Red, Blue, Green, Yellow, Orange}
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
	points, err = this.CurrentState.GetPoints()
	return
}

func (this *Game) IsGameOver() bool {
	return this.CurrentState.IsGameOver()
}
