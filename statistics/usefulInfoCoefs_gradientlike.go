package statistics

import (
	"fmt"
	"time"

	info "github.com/BabichMikhail/Hanabi/AIInformator"
	"github.com/BabichMikhail/Hanabi/game"
)

type AIWithCoefs interface {
	SetCoefs(part int, coefs ...float64)
	GetCoefs(part int) []float64
}

func RunGamesWithCoefs(count int, part int, aiType int, coefs []float64) float64 {
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
			aiType,
			aiType,
			aiType,
			aiType,
			aiType,
		},
		Count: count,
		Games: make([]GameStat, count, count),
	}

	for step := 0; step < count; step++ {
		g := game.NewGame(playerIds, game.Type_NormalGame)
		informator := info.NewInformator(g.CurrentState, g.InitState, g.Actions)

		for !g.IsGameOver() {
			AI := informator.NextAI(aiType)
			AI.(AIWithCoefs).SetCoefs(part, coefs...)
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

func GetUsefulCoefs(part, aiType int) []float64 {
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
	g := game.NewGame(playerIds, game.Type_NormalGame)
	informator := info.NewInformator(g.CurrentState, g.InitState, g.Actions)
	newAI := informator.NextAI(aiType)
	return newAI.(AIWithCoefs).GetCoefs(part)
}

func FindUsefulInfoCoefs_Gradient(part, aiType int) {
	time.Sleep(5 * time.Second)
	fmt.Println("Start FindUsefulInfoCoefs_Gradient", part, aiType)
	delta := []float64{1.0, 0.7, 0.5, 0.3, 0.1, 0.05, 0.03}
	N := []int{1000, 5000, 10000, 10000, 10000, 12000, 15000}
	usefulCoefs := GetUsefulCoefs(part, aiType)

	for idx, d := range delta {
		max := RunGamesWithCoefs(N[idx], part, aiType, usefulCoefs)
		for {
			newMax := max
			length := len(usefulCoefs)
			newK := make([][]float64, 2*length)
			for i := 0; i < length; i++ {
				newKi1 := make([]float64, length)
				newKi2 := make([]float64, length)
				for j := 0; j < length; j++ {
					newKi1[j] = usefulCoefs[j]
					newKi2[j] = usefulCoefs[j]
				}
				newKi1[i] += d
				newKi2[i] -= d

				newK[2*i] = newKi1
				newK[2*i+1] = newKi2
			}

			chans := make(chan struct{}, 2*length)
			for _, k := range newK {
				f := func(k []float64) {
					if result := RunGamesWithCoefs(N[idx], part, aiType, k); result > newMax {
						usefulCoefs = k
						newMax = result
						fmt.Print("NewMax:", result)
					} else {
						fmt.Print("Fail  :", result)
					}
					for i := 0; i < len(k); i++ {
						fmt.Print(" ", k[i])
					}
					fmt.Println()
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
			fmt.Print("Max   :", max)
			for i := 0; i < len(usefulCoefs); i++ {
				fmt.Print(" ", usefulCoefs[i])
			}
			fmt.Println()
		}
	}

	fmt.Println("OK")
}
