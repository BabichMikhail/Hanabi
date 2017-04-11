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

func NewInformator(currentGameState *game.GameState, actions []game.Action) *Informator {
	info := new(Informator)
	info.actions = actions
	info.currentState = currentGameState
	info.gameStates = map[int]game.GameState{}
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
	playerInfo := state.GetPlayerGameInfoByPos(state.CurrentPosition)
	return ai.NewAI(playerInfo, info.actions, aiType, info)
}

func (info *Informator) ApplyAction(action game.Action) error {
	state := info.getCurrentState()
	if err := state.ApplyAction(action); err != nil {
		return err
	}

	info.actions = append(info.actions, action)
	info.gameStates[len(info.actions)] = *state.Copy()
	return nil
}
