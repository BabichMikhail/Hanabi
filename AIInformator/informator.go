package informator

import (
	"fmt"
	"strconv"

	ai "github.com/BabichMikhail/Hanabi/AI"
	game "github.com/BabichMikhail/Hanabi/game"
)

type QReadFunc func(*game.PlayerGameInfo) float64
type QUpdateFunc func(*game.GameState, float64)

type KeyCacheActions struct {
	PlayerInfoHash string
	AIType         int
}

type Informator struct {
	actions       []game.Action
	gameStates    map[int]game.GameState
	currentState  *game.GameState
	QRead         func(*game.PlayerGameInfo) float64
	QUpdate       func(*game.GameState, float64)
	Cache         map[int]interface{}
	CachePInfo    map[string]*game.PlayerGameInfo
	CacheActions  map[KeyCacheActions]*game.Action
	isLearnOnStep bool
}

func NewInformator(currentGameState *game.GameState, initialGameState *game.GameState, actions []game.Action, qRead QReadFunc, qUpdate QUpdateFunc) *Informator {
	info := new(Informator)
	info.actions = actions
	info.currentState = currentGameState
	info.gameStates = map[int]game.GameState{}
	info.gameStates[0] = *initialGameState.Copy()
	info.gameStates[len(info.actions)] = *currentGameState.Copy()
	info.QRead = qRead
	info.QUpdate = qUpdate
	info.isLearnOnStep = false
	info.CachePInfo = map[string]*game.PlayerGameInfo{}
	info.CacheActions = map[KeyCacheActions]*game.Action{}
	return info
}

func (info *Informator) getCurrentState() *game.GameState {
	return info.currentState
}

func (info *Informator) GetActions() []game.Action {
	return info.actions
}

func (info *Informator) NextAI(aiType int) ai.AI {
	state := info.getCurrentState()
	infoType := game.InfoTypeUsually
	if aiType == ai.Type_AICheater {
		infoType = game.InfoTypeCheat
	} else if aiType == ai.Type_AIFullCheater {
		infoType = game.InfoTypeFullCheat
	}
	playerInfo := state.GetPlayerGameInfoByPos(state.CurrentPosition, infoType)
	return ai.NewAI(playerInfo, info.actions, aiType, info)
}

func (info *Informator) getPlayerState(step, infoType int, currentPosition ...int) game.PlayerGameInfo {
	state := info.gameStates[step]
	pos := state.CurrentPosition
	if len(currentPosition) == 1 {
		pos = currentPosition[0]
	}
	return state.GetPlayerGameInfoByPos(pos, infoType)
}

func (info *Informator) GetPlayerState(step int) game.PlayerGameInfo {
	_, ok := info.gameStates[step]
	if !ok {
		prevStep := step
		var prevState game.GameState
		for !ok {
			prevStep--
			prevState, ok = info.gameStates[prevStep]
		}

		prevState = *prevState.Copy()
		for i := prevStep; i < step; i++ {
			prevState.ApplyAction(&info.actions[i])
			info.gameStates[i+1] = *prevState.Copy()
		}
		_ = info.gameStates[step]
	}
	return info.getPlayerState(step, game.InfoTypeUsually, info.currentState.CurrentPosition)
}

func (info *Informator) ApplyAction(action *game.Action) error {
	if false && info.QUpdate != nil && !info.isLearnOnStep {
		newInformator := info.Copy()
		newInformator.QUpdate = nil
		saveState := newInformator.getCurrentState().Copy()
		for !newInformator.getCurrentState().IsGameOver() {
			AI := newInformator.NextAI(ai.Type_AIUsefulInfoAndMinMax)
			newAction := AI.GetAction()
			err := newInformator.ApplyAction(newAction)
			if err != nil {
				panic(err)
			}
		}
		info.isLearnOnStep = true
		state := newInformator.getCurrentState()
		points, err := state.GetPoints()
		if err == nil {
			info.QUpdate(saveState, float64(points))
		}
	}

	state := info.getCurrentState()
	if err := state.ApplyAction(action); err != nil {
		return err
	}

	info.actions = append(info.actions, *action)
	info.gameStates[len(info.actions)] = *state.Copy()
	return nil
}

func (info *Informator) Copy() *Informator {
	newInfo := new(Informator)
	newInfo.currentState = info.currentState.Copy()
	newInfo.gameStates = map[int]game.GameState{}
	newInfo.QRead = info.QRead
	newInfo.QUpdate = info.QUpdate
	for step, state := range info.gameStates {
		newInfo.gameStates[step] = *state.Copy()
	}

	newInfo.actions = make([]game.Action, len(info.actions))
	for i, action := range info.actions {
		newInfo.actions[i] = action
	}
	return newInfo
}

func (info *Informator) GetQualitativeAssessmentOfState(playerInfo *game.PlayerGameInfo) float64 {
	return info.QRead(playerInfo)
}

func (info *Informator) SetCache(data interface{}) {
	pos := info.currentState.CurrentPosition
	if info.Cache == nil {
		info.Cache = map[int]interface{}{}
	}
	info.Cache[pos] = data
}

func (info *Informator) GetCache() interface{} {
	pos := info.currentState.CurrentPosition
	return info.Cache[pos]
}

// Don't use this function in AI package. Use only for debugging
func (info *Informator) CheckAvailablePlayerInformation(availableGameInfo []*game.AvailablePlayerGameInfo, step int) int {
	okIdx := -1
	playerInfo := info.getPlayerState(step, game.InfoTypeCheat)
	fmt.Println("This Step:", playerInfo.Step)
	hashRes1 := info.PlayerInfoHash(&playerInfo)
	for idx, information := range availableGameInfo {
		availablePlayerInfo := information.PlayerInfo
		if availablePlayerInfo.CurrentPosition != playerInfo.CurrentPosition {
			panic(fmt.Sprint("Different current positions:", availablePlayerInfo.CurrentPosition, playerInfo.CurrentPosition))
		}

		hashRes2 := info.PlayerInfoHash(availablePlayerInfo)
		for pos, cards := range availablePlayerInfo.PlayerCards {
			for j := 0; j < len(cards); j++ {
				card1 := &cards[j]
				card2 := &playerInfo.PlayerCards[pos][j]
				if !card1.KnownColor || !card1.KnownValue {
					panic("1. Bad information about cards")
				}

				if len(card1.ProbabilityColors) != 1 || len(card1.ProbabilityValues) != 1 {
					panic("1. Bad probabilities for card")
				}

				if card1.Color == game.NoneColor || card1.Value == game.NoneValue {
					panic(fmt.Sprint("Bad Color Or Value", pos, j))
				}

				if !card2.KnownColor || !card2.KnownValue {
					panic("2. Bad information about cards 2")
				}

				if len(card2.ProbabilityColors) != 1 || len(card2.ProbabilityValues) != 1 {
					panic("2. Bad probabilities for card")
				}

				if card2.Color == game.NoneColor || card2.Value == game.NoneValue {
					panic("Bad Color Or Value")
				}

				if card1.Color != card2.Color || card1.Value != card2.Value {
					goto needContinue
				}
			}
		}

		for pos, cards := range availablePlayerInfo.PlayerCardsInfo {
			for j := 0; j < len(cards); j++ {
				card1 := &cards[j]
				card2 := &playerInfo.PlayerCardsInfo[pos][j]
				if card1.KnownColor != card2.KnownColor || card1.KnownValue != card2.KnownValue {
					panic("3. Bad information about cards ")
				}

				if len(card1.ProbabilityColors) != len(card2.ProbabilityColors) {
					panic("4. Bad information about cards")
				}

				for color, _ := range card1.ProbabilityColors {
					if _, ok := card2.ProbabilityColors[color]; !ok {
						panic("5. Bad information about cards")
					}
				}

				if len(card1.ProbabilityValues) != len(card2.ProbabilityValues) {
					panic("6. Bad information about cards")
				}

				for value, _ := range card1.ProbabilityValues {
					if _, ok := card2.ProbabilityValues[value]; !ok {
						panic("7. Bad information about cards")
					}
				}

				if card1.Color != card2.Color || card1.Value != card2.Value {
					goto needContinue
				}
			}
		}

		if availablePlayerInfo.BlueTokens != playerInfo.BlueTokens || availablePlayerInfo.RedTokens != playerInfo.RedTokens {
			continue
		}

		for color, card := range availablePlayerInfo.TableCards {
			if card.Value != playerInfo.TableCards[color].Value {
				panic(fmt.Sprint("Bad table cards", "\n", color, "\n", availablePlayerInfo.TableCards, "\n", playerInfo.TableCards))
			}
		}

		if okIdx != -1 {
			panic("Magic")
		}

		if hashRes1 != hashRes2 {
			panic("Different hash results\n" + fmt.Sprintf("Hash1: %s\n", hashRes1) + fmt.Sprintf("Hash2: %s", hashRes2))
		}
		okIdx = idx
	needContinue:
	}

	if okIdx == -1 {
		panic("Bad availablePlayerGameInformation")
	}
	return okIdx
}

func (info *Informator) ForcePlayerInfoHash(playerInfo *game.PlayerGameInfo) string {
	result := strconv.Itoa(playerInfo.CurrentPosition)
	for pos, cards := range playerInfo.PlayerCards {
		result += fmt.Sprintf("pl(%d)", pos)
		for i := 0; i < len(cards); i++ {
			card := &cards[i]
			result += strconv.Itoa(int(card.Color)) + strconv.Itoa(int(card.Value))
			probSum := 0
			for color, _ := range card.ProbabilityColors {
				probSum += 1 << uint(color)
			}
			for value, _ := range card.ProbabilityValues {
				probSum += 1 << (uint(value) + 5)
			}
			result += strconv.Itoa(probSum)
		}
	}
	result += "info"
	result += strconv.Itoa(playerInfo.BlueTokens) + strconv.Itoa(playerInfo.RedTokens) + strconv.Itoa(playerInfo.Step)
	result += "table"
	for _, color := range game.ColorsTable {
		value := playerInfo.TableCards[color].Value
		result += strconv.Itoa(int(value))
	}

	return result
}

func (info *Informator) PlayerInfoHash(playerInfo *game.PlayerGameInfo) string {
	if playerInfo.HashKey != nil {
		return *playerInfo.HashKey
	}

	result := info.ForcePlayerInfoHash(playerInfo)
	playerInfo.HashKey = &result
	return result
}

func (info *Informator) SetProbabilities(playerInfo *game.PlayerGameInfo) {
	hashPlayerInfo := info.PlayerInfoHash(playerInfo)
	if pinfo, ok := info.CachePInfo[hashPlayerInfo]; ok {
		copyPInfo := pinfo.Copy()
		*playerInfo = *copyPInfo
		return
	}
	playerInfo.SetProbabilities(false, false)
	info.CachePInfo[hashPlayerInfo] = playerInfo.Copy()
}

func (info *Informator) GetAction(playerInfo *game.PlayerGameInfo, aiType int, history []game.Action) *game.Action {
	key := KeyCacheActions{
		PlayerInfoHash: info.PlayerInfoHash(playerInfo),
		AIType:         aiType,
	}
	if action, ok := info.CacheActions[key]; ok {
		return action
	}
	newAI := ai.NewAI(*playerInfo, history, aiType, info)
	return newAI.GetAction()
}
