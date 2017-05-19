package game

import "fmt"

type ActionType int

const (
	TypeActionInformationColor = iota
	TypeActionInformationValue
	TypeActionDiscard
	TypeActionPlaying
)

type Action struct {
	ActionType     ActionType `json:"action_type"`
	PlayerPosition int        `json:"player_position"`
	Value          int        `json:"value"`
}

func (game *Game) AppendAction(action *Action, err error) (*Action, error) {
	if err == nil {
		game.Actions = append(game.Actions, *action)
	}
	return action, err
}

func (action Action) String() string {
	names := map[int]string{
		TypeActionInformationColor: "InfoColor",
		TypeActionInformationValue: "InfoValue",
		TypeActionPlaying:          "Play",
		TypeActionDiscard:          "Discard",
	}
	return fmt.Sprintf("{%d %s %d}", action.PlayerPosition, names[int(action.ActionType)], action.Value)
}

func (a1 *Action) Equal(a2 *Action) bool {
	return a1.ActionType == a2.ActionType && a1.PlayerPosition == a2.PlayerPosition && a1.Value == a2.Value
}

func (state *GameState) NewAction(actionType ActionType, playerPosition int, value int) *Action {
	action := NewAction(actionType, playerPosition, value)
	state.IncreaseStep()
	return action
}

func (action Action) IsInfoAction() bool {
	return action.ActionType == TypeActionInformationColor || action.ActionType == TypeActionInformationValue
}

func NewAction(actionType ActionType, playerPosition int, value int) *Action {
	return &Action{
		ActionType:     actionType,
		PlayerPosition: playerPosition,
		Value:          value,
	}
}

func (action *Action) Copy() *Action {
	return &Action{
		ActionType:     action.ActionType,
		PlayerPosition: action.PlayerPosition,
		Value:          action.Value,
	}
}

func (game *Game) ApplyAction(action Action) (err error) {
	switch action.ActionType {
	case TypeActionDiscard:
		_, err = game.NewActionDiscard(action.PlayerPosition, action.Value)
	case TypeActionInformationColor:
		_, err = game.NewActionInformationColor(action.PlayerPosition, CardColor(action.Value))
	case TypeActionInformationValue:
		_, err = game.NewActionInformationValue(action.PlayerPosition, CardValue(action.Value))
	case TypeActionPlaying:
		_, err = game.NewActionPlaying(action.PlayerPosition, action.Value)
	}
	return
}

func (state *GameState) ApplyAction(action *Action) (err error) {
	//fmt.Println("apply action")
	if action.IsInfoAction() {
		if action.PlayerPosition == state.CurrentPosition {
			panic("Bad Action")
		}
	} else {
		if action.PlayerPosition != state.CurrentPosition {
			panic("Bad Action")
		}
	}
	switch action.ActionType {
	case TypeActionDiscard:
		_, err = state.NewActionDiscard(action.PlayerPosition, action.Value)
	case TypeActionInformationColor:
		_, err = state.NewActionInformationColor(action.PlayerPosition, CardColor(action.Value))
	case TypeActionInformationValue:
		_, err = state.NewActionInformationValue(action.PlayerPosition, CardValue(action.Value))
	case TypeActionPlaying:
		_, err = state.NewActionPlaying(action.PlayerPosition, action.Value)
	}
	return
}

func (state *GameState) IncreaseStep() {
	state.Step++
	state.CurrentPosition++
	if state.CurrentPosition/len(state.PlayerStates) == 1 {
		state.CurrentPosition = 0
		state.Round++
	}
}
