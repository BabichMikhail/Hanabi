package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	ai "github.com/BabichMikhail/Hanabi/AI"
	"github.com/BabichMikhail/Hanabi/game"
)

type GeneticAlgorithm struct {
	Nodes GeneticNodes
	N     int
}

type GeneticNode struct {
	Coefs   []float64
	Result  float64
	Percent float64
}

type GeneticNodes []GeneticNode

func (nodes GeneticNodes) Len() int {
	return len(nodes)
}

func (nodes GeneticNodes) Less(i, j int) bool {
	return nodes[i].Result > nodes[j].Result
}

func (nodes GeneticNodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func NewGeneticAlgorithm() *GeneticAlgorithm {
	gen := new(GeneticAlgorithm)
	gen.N = 120
	a := -2.0
	b := 4.0
	h := (b - a) / float64(gen.N)
	gen.Nodes = make([]GeneticNode, gen.N+2, gen.N+2)
	for i := 0; i <= gen.N; i++ {
		k := a + h*float64(i)
		gen.Nodes[i] = GeneticNode{
			Coefs: []float64{
				k, k, k, k, k, k, k, k,
			},
			Result:  0.0,
			Percent: 0.0,
		}
	}
	gen.Nodes[gen.N+1] = GeneticNode{
		Coefs: []float64{
			2.1, -0.9, 1.05, 1.0, 0.1, 0.04, 0.01, 0.07,
		},
		Result:  0.0,
		Percent: 0.0,
	}
	return gen
}

func (gen *GeneticAlgorithm) RunGamesWithCoefs(count int, aiTypes []int, coefs []float64) float64 {
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
		AITypes: aiTypes,
		Count:   count,
		Games:   make([]GameStat, count, count),
	}

	for step := 0; step < count; step++ {
		g := game.NewGame(playerIds)
		actions := []game.Action{}
		newAITypes := make([]int, len(stat.AITypes), len(stat.AITypes))
		for idx, state := range g.CurrentState.PlayerStates {
			newAITypes[posById[state.PlayerId]] = stat.AITypes[idx]
		}

		for !g.IsGameOver() {
			pos := g.CurrentState.CurrentPosition
			playerInfo := g.CurrentState.GetPlayerGameInfoByPos(pos)
			AI := ai.NewAI(playerInfo, actions, newAITypes[pos])
			AI.(*ai.AIUsefulInformationV3).SetCoefs(coefs[0], coefs[1], coefs[2], coefs[3], coefs[4], coefs[5], coefs[6], coefs[7])
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
		stat.Games[step].Points = gamePoints
		stat.Games[step].RedTokens = g.CurrentState.RedTokens
		stat.Games[step].Step = len(g.Actions)
	}
	stat.SetCharacteristics()
	return stat.Medium
}

func (gen *GeneticAlgorithm) NewDescendant(idx1, idx2 int) *GeneticNode {
	node1 := &gen.Nodes[idx1]
	node2 := &gen.Nodes[idx2]
	crossGen := rand.Intn(len(node1.Coefs))
	newNode := new(GeneticNode)
	newNode.Coefs = make([]float64, len(node1.Coefs), len(node1.Coefs))
	newNode.Percent = 0
	newNode.Result = 0
	crossType := rand.Intn(2)
	if crossType == 0 {
		copy(newNode.Coefs, node2.Coefs)
		newNode.Coefs[crossGen] = node1.Coefs[crossGen]
	} else {
		copy(newNode.Coefs, node1.Coefs)
		newNode.Coefs[crossGen] = node2.Coefs[crossGen]
	}
	return newNode
}

func (gen *GeneticAlgorithm) NewMutation(idx int) *GeneticNode {
	node := &gen.Nodes[idx]
	newNode := new(GeneticNode)
	newNode.Coefs = make([]float64, len(node.Coefs), len(node.Coefs))
	newNode.Percent = 0
	newNode.Result = 0
	newCoef := rand.Float64()*8 - 2
	newNode.Coefs[rand.Intn(len(node.Coefs))] = newCoef
	copy(newNode.Coefs, node.Coefs)
	return newNode
}

func (gen *GeneticAlgorithm) GetRandIdx() int {
	idx := -1
	sumPercent := 0.0
	for sumPercent <= rand.Float64() {
		idx++
		sumPercent += gen.Nodes[idx].Percent
	}
	return idx
}

func (gen *GeneticAlgorithm) FindUsefulInfoV3Coefs() {
	time.Sleep(5 * time.Second)
	fmt.Println("Start algorithm")
	aiTypes := []int{
		ai.Type_AIUsefulInformationV3,
		ai.Type_AIUsefulInformationV3,
		ai.Type_AIUsefulInformationV3,
		ai.Type_AIUsefulInformationV3,
		ai.Type_AIUsefulInformationV3,
	}

	gamesCount := 500
	lowValue := 15.5

	for {
		isContinue := false
		sum := 0.0
		min := 26.0
		var minIdx int
		max := -1.0
		var maxIdx int
		chans := make(chan int, len(gen.Nodes))
		readyCount := 0
		for i := 0; i < len(gen.Nodes); i++ {
			if gen.Nodes[i].Result > 0 {
				gen.Nodes[i].Percent = 0
				readyCount++
				continue
			}

			f := func(i int) {
				node := &gen.Nodes[i]
				node.Result = gen.RunGamesWithCoefs(gamesCount, aiTypes, node.Coefs)
				if node.Result < lowValue {
					isContinue = true
				}
			}
			go func(i int) {
				f(i)
				chans <- i
			}(i)
		}

		fmt.Println("Results:")
		for j := 0; j < len(gen.Nodes)-readyCount; j++ {
			<-chans
		}

		sort.Sort(gen.Nodes)
		N := gen.N
		for i := 0; i < N/2; i++ {
			node := &gen.Nodes[i]
			sum += math.Pow(1/(25-node.Result), 2)
			if node.Result < min {
				min = node.Result
				minIdx = i
			}
			if node.Result > max {
				max = node.Result
				maxIdx = i
			}
			fmt.Println("Result: ", node.Result, node.Coefs)
		}

		for i := 0; i < len(gen.Nodes); i++ { // i<N/2
			node := &gen.Nodes[i]
			node.Percent = math.Pow(1/(25-node.Result), 2) / sum
		}

		newNodes := make([]GeneticNode, gen.N, gen.N)
		for i := 0; i < N/2; i++ {
			newNodes[i] = gen.Nodes[i]
		}

		for i := N / 2; i < N/2+N/4; i++ {
			idx1 := gen.GetRandIdx()
			idx2 := gen.GetRandIdx()
			for idx1 == idx2 {
				idx2 = gen.GetRandIdx()
			}
			newNodes[i] = *gen.NewDescendant(idx1, idx2)
		}

		for i := N/2 + N/4; i < N; i++ {
			idx := gen.GetRandIdx()
			newNodes[i] = *gen.NewMutation(idx)
		}

		fmt.Println("Minimum: ", min, gen.Nodes[minIdx].Coefs)
		fmt.Println("Maximum: ", max, gen.Nodes[maxIdx].Coefs)
		fmt.Println()
		if !isContinue {
			if gen.N == 60 {
				break
			}
			gen.N = 60
			gamesCount = 3000
			lowValue = 16.2
		}
		gen.Nodes = newNodes
	}

	fmt.Println("Results:")
	for i := 0; i < len(gen.Nodes); i++ {
		fmt.Println(gen.Nodes[i].Result, gen.Nodes[i].Coefs)
	}

	fmt.Println("OK")
}
