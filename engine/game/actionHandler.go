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

func (state *GameState) NewAction(actionType ActionType, playerPosition int, value int) Action {
	action := Action{
		ActionType:     actionType,
		PlayerPosition: playerPosition,
		Value:          value,
	}
	state.Actions = append(state.Actions, action)
	state.IncreaseStep()
	return action
}

func (state *GameState) IncreaseStep() {
	state.Step++
	state.CurrentPosition++
	if state.CurrentPosition/state.PlayerCount == 1 {
		state.CurrentPosition = 0
		state.Round++
	}
}
