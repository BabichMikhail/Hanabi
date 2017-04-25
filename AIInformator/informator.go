package informator

import (
	ai "github.com/BabichMikhail/Hanabi/AI"
	game "github.com/BabichMikhail/Hanabi/game"
)

type Informator struct {
	actions      []game.Action
	gameStates   map[int]game.GameState
	currentState *game.GameState
}

func NewInformator(currentGameState *game.GameState, initialGameState *game.GameState, actions []game.Action) *Informator {
	info := new(Informator)
	info.actions = actions
	info.currentState = currentGameState
	info.gameStates = map[int]game.GameState{}
	info.gameStates[0] = *initialGameState.Copy()
	info.gameStates[len(info.actions)] = *currentGameState.Copy()
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
	state := info.getCurrentState()
	if err := state.ApplyAction(action); err != nil {
		return err
	}

	info.actions = append(info.actions, *action)
	info.gameStates[len(info.actions)] = *state.Copy()
	return nil
}
