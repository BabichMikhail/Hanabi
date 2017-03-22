package statistics

import (
	"math"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/game"
)

type GameStat struct {
	Points    int `json:"points"`
	Step      int `json:"step"`
	RedTokens int `json:"red_tokens"`
}

type Stat struct {
	Count   int        `json:"count"`
	Medium  float64    `json:"medium"`
	Disp    float64    `json:"disp"`
	Asym    float64    `json:"asymmetry"`
	Kurt    float64    `json:"kurtosis"`
	AITypes []int      `json:"ai_types"`
	Games   []GameStat `json:"game_stat"`
}

func Medium(games []GameStat) float64 {
	sum := 0.0
	for _, g := range games {
		sum += float64(g.Points)
	}
	return sum / float64(len(games))
}

func CentralMoment(games []GameStat, med float64, pow float64) float64 {
	sum := 0.0
	for _, g := range games {
		sum += math.Pow(float64(g.Points)-med, pow)
	}
	return sum / float64(len(games))
}

func Dispersion(games []GameStat, med float64) float64 {
	return CentralMoment(games, med, 2)
}

func (stat *Stat) SetCharacteristics() {
	stat.Medium = Medium(stat.Games)
	stat.Disp = Dispersion(stat.Games, stat.Medium)
	stat.Asym = CentralMoment(stat.Games, stat.Medium, 3) / math.Pow(stat.Disp, 1.5)
	stat.Kurt = CentralMoment(stat.Games, stat.Medium, 4)/math.Pow(stat.Disp, 2) - 3
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
		Games:   make([]GameStat, count, count),
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
		stat.Games[i].Points = gamePoints
		stat.Games[i].RedTokens = g.CurrentState.RedTokens
		stat.Games[i].Step = len(g.Actions)
		if i == 0 || gamePoints > maxPoints {
			bestGame = g
			maxPoints = gamePoints
		}
	}
	stat.SetCharacteristics()
	return stat, bestGame
}
