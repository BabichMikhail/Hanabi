package game

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

func (state *GameState) NewAction(actionType ActionType, playerPosition int, value int) *Action {
	action := NewAction(actionType, playerPosition, value)
	state.IncreaseStep()
	return action
}

func NewAction(actionType ActionType, playerPosition int, value int) *Action {
	return &Action{
		ActionType:     actionType,
		PlayerPosition: playerPosition,
		Value:          value,
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
