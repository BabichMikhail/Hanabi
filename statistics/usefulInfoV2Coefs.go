package statistics

import (
	"fmt"

	ai "github.com/BabichMikhail/Hanabi/AI"
	info "github.com/BabichMikhail/Hanabi/AIInformator"
	"github.com/BabichMikhail/Hanabi/game"
)

func RunGamesWithCoefs(count int, kPlayByValue, kPlayByColor, kInfoValue, kInfoColor float64) float64 {
	playersCount := 5
	pseudoIds := make([]int, playersCount, playersCount)
	for i := 0; i < playersCount; i++ {
		pseudoIds[i] = i + 1
	}

	playerIds := pseudoIds
	posById := map[int]int{}
	for i := 0; i < len(playerIds); i++ {
		posById[playerIds[i]] = i
	}

	stat := Stat{
		AITypes: []int{
			ai.Type_AIUsefulInformationV2,
			ai.Type_AIUsefulInformationV2,
			ai.Type_AIUsefulInformationV2,
			ai.Type_AIUsefulInformationV2,
			ai.Type_AIUsefulInformationV2,
		},
		Count: count,
		Games: make([]GameStat, count, count),
	}

	for step := 0; step < count; step++ {
		g := game.NewGame(playerIds)
		informator := info.NewInformator(g.CurrentState, g.Actions)

		newAITypes := make([]int, len(stat.AITypes), len(stat.AITypes))
		for idx, state := range g.CurrentState.PlayerStates {
			newAITypes[posById[state.PlayerId]] = stat.AITypes[idx]
		}

		for !g.IsGameOver() {
			pos := g.CurrentState.CurrentPosition
			AI := informator.NextAI(newAITypes[pos])
			AI.(*ai.AIUsefulInformationV2).SetCoefs(kPlayByValue, kPlayByColor, kInfoValue, kInfoColor)
			action := AI.GetAction()
			informator.ApplyAction(action)
		}
		gamePoints, _ := g.GetPoints()
		stat.Games[step].Points = gamePoints
		stat.Games[step].RedTokens = g.CurrentState.RedTokens
		stat.Games[step].Step = len(g.Actions)
	}
	stat.SetCharacteristics()
	return stat.Medium
}

func FindUsefulInfoV2Coefs() {
	delta := []float64{1.0, 0.5, 0.1, 0.05}
	N := []int{10000, 10000, 10000, 10000}
	usefulCoefs := []float64{2.1, -0.9, 1.05, 1.0}
	max := RunGamesWithCoefs(10000, usefulCoefs[0], usefulCoefs[1], usefulCoefs[2], usefulCoefs[3])
	for idx, d := range delta {
		for {
			newMax := max
			newK := [][]float64{
				{usefulCoefs[0] + d, usefulCoefs[1], usefulCoefs[2], usefulCoefs[3]},
				{usefulCoefs[0] - d, usefulCoefs[1], usefulCoefs[2], usefulCoefs[3]},
				{usefulCoefs[0], usefulCoefs[1] + d, usefulCoefs[2], usefulCoefs[3]},
				{usefulCoefs[0], usefulCoefs[1] - d, usefulCoefs[2], usefulCoefs[3]},
				{usefulCoefs[0], usefulCoefs[1], usefulCoefs[2] + d, usefulCoefs[3]},
				{usefulCoefs[0], usefulCoefs[1], usefulCoefs[2] - d, usefulCoefs[3]},
				{usefulCoefs[0], usefulCoefs[1], usefulCoefs[2], usefulCoefs[3] + d},
				{usefulCoefs[0], usefulCoefs[1], usefulCoefs[2], usefulCoefs[3] - d},
			}

			chans := make(chan struct{}, 8)
			for _, k := range newK {
				f := func(k []float64) {
					if result := RunGamesWithCoefs(N[idx], k[0], k[1], k[2], k[3]); result > newMax {
						usefulCoefs = k
						newMax = result
						fmt.Println("NewMax:", result, k[0], k[1], k[2], k[3])
					} else {
						fmt.Println("Fail  :", result, k[0], k[1], k[2], k[3])
					}
					chans <- struct{}{}
				}
				go f(k)
			}

			for i := 0; i < len(newK); i++ {
				<-chans
			}

			if newMax <= max {
				fmt.Println(newMax, max, d)
				break
			}
			max = newMax
			fmt.Println("Max   :", max, usefulCoefs[0], usefulCoefs[1], usefulCoefs[2], usefulCoefs[3])
		}
	}

	fmt.Println("OK")
}
