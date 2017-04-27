package informator

import (
	ai "github.com/BabichMikhail/Hanabi/AI"
	game "github.com/BabichMikhail/Hanabi/game"
)

type QReadFunc func(*game.PlayerGameInfo) float64
type QUpdateFunc func(*game.GameState, float64)

type Informator struct {
	actions       []game.Action
	gameStates    map[int]game.GameState
	currentState  *game.GameState
	QRead         func(*game.PlayerGameInfo) float64
	QUpdate       func(*game.GameState, float64)
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

func (info *Informator) GetPlayerState(step int) game.PlayerGameInfo {
	state, ok := info.gameStates[step]
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
		state = info.gameStates[step]
	}
	return state.GetPlayerGameInfoByPos(info.currentState.CurrentPosition, game.InfoTypeUsually)
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
