package statistics

import (
	"math"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/game"
)

type Stat struct {
	Count   int       `json:"count"`
	Medium  float64   `json:"medium"`
	Disp    float64   `json:"disp"`
	Asym    float64   `json:"asymmetry"`
	Kurt    float64   `json:"kurtosis"`
	AITypes []int     `json:"ai_types"`
	Values  []float64 `json:"values"`
}

func Medium(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func CentralMoment(values []float64, med float64, pow float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += math.Pow(v-med, pow)
	}
	return sum / float64(len(values))
}

func Dispersion(values []float64, med float64) float64 {
	return CentralMoment(values, med, 2)
}

func (stat *Stat) SetCharacteristics() {
	stat.Medium = Medium(stat.Values)
	stat.Disp = Dispersion(stat.Values, stat.Medium)
	stat.Asym = CentralMoment(stat.Values, stat.Medium, 3) / math.Pow(stat.Disp, 1.5)
	stat.Kurt = CentralMoment(stat.Values, stat.Medium, 4)/math.Pow(stat.Disp, 2) - 3
	return
}

func RunGames(aiTypes []int, playerIds []int, count int) (Stat, *game.Game) {
	playersCount := len(aiTypes)
	if playersCount > 5 && playersCount < 2 {
		panic("bad players count")
	}

	pseudoIds := make([]int, playersCount, playersCount)
	for i := 0; i < playersCount; i++ {
		pseudoIds[i] = i + 1
	}

	stat := Stat{
		AITypes: aiTypes,
		Count:   count,
		Values:  make([]float64, count, count),
	}

	posById := map[int]int{}
	for i := 0; i < len(playerIds); i++ {
		posById[playerIds[i]] = i
	}

	var bestGame *game.Game
	maxPoints := -1
	for i := 0; i < count; i++ {
		g := game.NewGame(playerIds)
		actions := []game.Action{}
		newAITypes := make([]int, len(aiTypes), len(aiTypes))
		for idx, state := range g.CurrentState.PlayerStates {
			newAITypes[posById[state.PlayerId]] = aiTypes[idx]
		}

		for !g.IsGameOver() {
			pos := g.CurrentState.CurrentPosition
			playerInfo := g.CurrentState.GetPlayerGameInfoByPos(pos)
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
		stat.Values[i] = float64(gamePoints)
		if i == 0 || gamePoints > maxPoints {
			bestGame = g
			maxPoints = gamePoints
		}
	}
	stat.SetCharacteristics()
	return stat, bestGame
}
