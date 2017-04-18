package statistics

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	ai "github.com/BabichMikhail/Hanabi/AI"
	info "github.com/BabichMikhail/Hanabi/AIInformator"
	"github.com/BabichMikhail/Hanabi/game"
)

type GeneticAlgorithm struct {
	Nodes       GeneticNodes
	Current     int
	GamesCounts []int
	N           []int
	LowValues   []float64
	Mutations   []float64
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
	gen.Current = 0
	gen.N = []int{120, 80, 40, 20}
	gen.GamesCounts = []int{1000, 4000, 12000, 40000}
	gen.LowValues = []float64{15.0, 15.6, 16.1, 16.5}
	gen.Mutations = []float64{1.3, 0.8, 0.4, 0.1}

	N := gen.N[gen.Current]
	a := -2.0
	b := 4.0
	h := (b - a) / float64(N)
	gen.Nodes = make([]GeneticNode, N+2, N+2)
	for i := 0; i <= N; i++ {
		k := a + h*float64(i)
		gen.Nodes[i] = GeneticNode{
			Coefs: []float64{
				k, k, k, k, k, k, k, k,
			},
			Result:  0.0,
			Percent: 0.0,
		}
	}
	gen.Nodes[N+1] = GeneticNode{
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
		informator := info.NewInformator(g.CurrentState, g.Actions)
		newAITypes := make([]int, len(stat.AITypes), len(stat.AITypes))
		for idx, state := range g.CurrentState.PlayerStates {
			newAITypes[posById[state.PlayerId]] = stat.AITypes[idx]
		}

		for !g.IsGameOver() {
			pos := g.CurrentState.CurrentPosition
			AI := informator.NextAI(newAITypes[pos])
			AI.(*ai.AIUsefulInfoV3AndParts).SetCoefs(0, coefs[0], coefs[1], coefs[2], coefs[3], coefs[4], coefs[5], coefs[6], coefs[7])
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

func (gen *GeneticAlgorithm) NewMutationAbsolute(idx int) *GeneticNode {
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

func (gen *GeneticAlgorithm) NewMutationRelative(idx int, k float64, N int) *GeneticNode {
	node := &gen.Nodes[idx]
	newNode := new(GeneticNode)
	newNode.Coefs = make([]float64, len(node.Coefs), len(node.Coefs))
	newNode.Percent = 0
	newNode.Result = 0
	for i := 0; i < N; i++ {
		newCoef := (rand.Float64()*gen.Mutations[gen.Current] - gen.Mutations[gen.Current]/2) * k
		newNode.Coefs[rand.Intn(len(node.Coefs))] += newCoef
	}
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

	updateAll := true
	repeats := 0

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
			if !updateAll && gen.Nodes[i].Result > 0 {
				gen.Nodes[i].Percent = 0
				readyCount++
				continue
			}

			f := func(i int) {
				node := &gen.Nodes[i]
				node.Result = gen.RunGamesWithCoefs(gen.GamesCounts[gen.Current], aiTypes, node.Coefs)
				if node.Result < gen.LowValues[gen.Current] {
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
		N := gen.N[gen.Current]
		for i := 0; i < N; i++ {
			node := &gen.Nodes[i]
			sum += 1 / (25 - node.Result)
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

		newNodes := make([]GeneticNode, N, N)
		for i := 0; i < N; i++ {
			node := &gen.Nodes[i]
			node.Percent = 1 / (25 - node.Result) / sum
		}

		for i := 0; i < N/4; i++ {
			newNodes[i] = gen.Nodes[i]
		}

		for i := N / 4; i < N/2; i++ {
			idx1 := gen.GetRandIdx()
			idx2 := gen.GetRandIdx()
			for idx1 == idx2 {
				idx2 = gen.GetRandIdx()
			}
			newNodes[i] = *gen.NewDescendant(idx1, idx2)
		}

		for i := N / 2; i < 2*N/3; i++ {
			idx := gen.GetRandIdx()
			newNodes[i] = *gen.NewMutationAbsolute(idx)
		}

		for i := 2 * N / 3; i < 5*N/6; i++ {
			idx := gen.GetRandIdx()
			newNodes[i] = *gen.NewMutationRelative(idx, 2, 1)
		}

		for i := 5 * N / 6; i < N; i++ {
			idx := gen.GetRandIdx()
			newNodes[i] = *gen.NewMutationRelative(idx, 0.2, 4)
		}

		fmt.Println("Minimum: ", min, gen.Nodes[minIdx].Coefs)
		fmt.Println("Maximum: ", max, gen.Nodes[maxIdx].Coefs)
		fmt.Println()
		updateAll = false
		if !isContinue {
			if gen.Current == len(gen.N)-1 {
				if repeats > 2 {
					break
				} else {
					repeats++
				}
			} else {
				gen.Current++
			}
			updateAll = true
		}
		gen.Nodes = newNodes
	}

	fmt.Println("Results:")
	for i := 0; i < len(gen.Nodes); i++ {
		fmt.Println(gen.Nodes[i].Result, gen.Nodes[i].Coefs)
	}

	fmt.Println("OK")
}
