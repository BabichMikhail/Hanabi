package game

import (
	"encoding/json"
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
	Points       int       `json:"points"`
}

func NewGame(ids []int) *Game {
	this := new(Game)
	this.Seed = time.Now().UTC().UnixNano()
	this.Points = 0
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
	RandomIntPermutation(ids)
	this.PlayerCount = len(ids)
	this.Actions = []Action{}
	state := NewGameState(ids, cards, this.PlayerCount)
	this.InitState = state
	this.CurrentState = state.Copy()

	return this
}

func (this *Game) SprintGame() string {
	b, err := json.Marshal(this)
	if err != nil {
		return ""
	}
	return fmt.Sprintln(string(b))
}

func (game *Game) GetPlayerPositionById(id int) (pos int, err error) {
	return game.CurrentState.GetPlayerPositionById(id)
}

func (game *Game) GetPoints() (points int, err error) {
	if game.Points != 0 {
		return game.Points, nil
	}

	points, err = game.CurrentState.GetPoints()
	if game.IsGameOver() {
		game.Points = points
	}
	return
}

func (game *Game) IsGameOver() bool {
	if game.Points > 0 {
		return true
	}
	return game.CurrentState.IsGameOver()
}
