package engine

import (
	"encoding/json"
	"fmt"
)

type Game struct {
	PlayerCount int         `json:"player_count"`
	GameStatus  []GameState `json:"states"`
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

func (this *Game) SprintGame() string {
	b, err := json.Marshal(this)
	if err != nil {
		return ""
	}
	return fmt.Sprintln(string(b))
}
