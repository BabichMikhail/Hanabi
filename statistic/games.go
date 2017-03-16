package statistic

import (
	"fmt"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/game"
)

type Stat struct {
	Count  int
	Medium float64
	AIType []int
}

func RunGames(aiTypes []int, count int) Stat {
	playersCount := len(aiTypes)
	if playersCount > 5 && playersCount < 2 {
		panic("bad players count")
	}

	pseudoIds := make([]int, playersCount, playersCount)
	for i := 0; i < playersCount; i++ {
		pseudoIds[i] = i + 1
	}

	points := 0
	stat := Stat{
		AIType: aiTypes,
		Count:  count,
	}

	for i := 0; i < count; i++ {
		g := game.NewGame(pseudoIds)
		actions := []game.Action{}
		newAITypes := make([]int, len(aiTypes), len(aiTypes))
		for idx, state := range g.CurrentState.PlayerStates {
			newAITypes[state.PlayerId-1] = aiTypes[idx]
		}

		for !g.IsGameOver() {
			pos := g.CurrentState.CurrentPosition
			playerInfo := g.GetPlayerGameInfo(pos)
			AI := ai.NewAI(playerInfo, actions, newAITypes[pos])
			action := AI.GetAction()
			actions = append(actions, action)
			switch action.ActionType {
			case game.TypeActionDiscard:
				g.NewActionDiscard(action.PlayerPosition, action.Value)
			case game.TypeActionInformationColor:
				g.NewActionInformationColor(action.PlayerPosition, game.CardColor(action.Value))
			case game.TypeActionInformationValue:
				g.NewActionInformationValue(action.PlayerPosition, game.CardValue(action.Value))
			case game.TypeActionPlaying:
				g.NewActionPlaying(action.PlayerPosition, action.Value)
			}
		}
		gamePoints, _ := g.GetPoints()
		points += gamePoints
	}
	stat.Medium = float64(points) / float64(count)
	fmt.Printf("Result: %d %f\n\n\n\n\n", stat.Count, stat.Medium)
	return stat
}
