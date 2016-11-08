package engine

import (
	"errors"
	"fmt"

	models "github.com/BabichMikhail/Hanabi/engine/models"
	wetalk "github.com/beego/wetalk/modules/models"
)

func GetAllStatuses() []int {
	return []int{models.GameWait, models.GameActive, models.GameInActive}
}

func MakeGame(id int, user wetalk.User) (game models.Game, err error) {
	if id <= 0 {
		errors.New(fmt.Sprintf("Can't make game with Id = %d", id))
	}
	game.Id = id
	game.Status = models.GameStatusName(models.GameWait)
	game.Owner = user.NickName
	player := models.Player{user.Id, user.NickName}
	game.Players = []models.Player{player}
	return game, nil
}
