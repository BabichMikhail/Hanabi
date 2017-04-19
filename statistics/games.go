package statistics

import (
	"math"
	"math/rand"

	info "github.com/BabichMikhail/Hanabi/AIInformator"
	"github.com/BabichMikhail/Hanabi/game"
)

type GameStat struct {
	Points    int                    `json:"points"`
	Step      int                    `json:"step"`
	RedTokens int                    `json:"red_tokens"`
	Values    map[game.CardColor]int `json:"table_values"`
}

type Stat struct {
	Count   int        `json:"count"`
	Wins    int        `json:"wins"`
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

func RunGames(aiTypes []int, playerIds []int, count int, fUpdateReady func(*int, int)) (Stat, []*game.Game) {
	playersCount := len(aiTypes)
	if playersCount > 5 && playersCount < 2 {
		panic("bad players count")
	}

	limit := 100
	if count < limit {
		limit = count
	}
	newCount := count / limit
	if count%limit > 0 {
		newCount++
	}
	count = newCount * limit

	stat := Stat{
		AITypes: aiTypes,
		Count:   count,
		Wins:    0,
		Games:   make([]GameStat, count, count),
	}

	posById := map[int]int{}
	for i := 0; i < len(playerIds); i++ {
		posById[playerIds[i]] = i
	}

	var bestGame *game.Game
	var worstGame *game.Game
	additionalGames := map[int]*game.Game{}
	if count > limit {
		for i := 0; i < 4; i++ {
			additionalGames[rand.Intn(count)] = nil
		}
	}
	maxPoints := -1
	minPoints := 26

	chans := make(chan struct{}, 2*limit)
	readyCount := 0
	go fUpdateReady(&readyCount, count)
	for j := 0; j < limit; j++ {
		go func(j int) {
			for k := 0; k < count/limit; k++ {
				i := k*limit + j
				g := game.NewGame(playerIds, game.Type_NormalGame)
				newAITypes := make([]int, len(aiTypes), len(aiTypes))
				for idx, state := range g.CurrentState.PlayerStates {
					newAITypes[posById[state.PlayerId]] = aiTypes[idx]
				}

				informator := info.NewInformator(g.CurrentState, g.InitState, g.Actions)
				for !g.IsGameOver() {
					pos := g.CurrentState.CurrentPosition
					AI := informator.NextAI(newAITypes[pos])
					action := AI.GetAction()
					err := informator.ApplyAction(action)
					if err != nil {
						panic(err)
					}
				}
				g.Actions = informator.GetActions()
				gamePoints, _ := g.GetPoints()
				stat.Games[i].Points = gamePoints
				stat.Games[i].RedTokens = g.CurrentState.RedTokens
				stat.Games[i].Step = len(g.Actions)
				stat.Games[i].Values = map[game.CardColor]int{}
				for _, color := range game.Colors {
					stat.Games[i].Values[color] = int(g.CurrentState.TableCards[color].Value)
				}

				if gamePoints > maxPoints {
					bestGame = g
					maxPoints = gamePoints
				}
				if gamePoints < minPoints {
					worstGame = g
					minPoints = gamePoints
				}

				if gamePoints == 25 {
					stat.Wins++
				}

				if _, ok := additionalGames[i]; ok {
					additionalGames[i] = g
				}
				readyCount++
			}
			chans <- struct{}{}
		}(j)
	}
	for i := 0; i < limit; i++ {
		<-chans
	}
	stat.SetCharacteristics()
	games := []*game.Game{worstGame, bestGame}
	for _, g := range additionalGames {
		games = append(games, g)
	}
	return stat, games
}
