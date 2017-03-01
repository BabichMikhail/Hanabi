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

func (game *Game) AppendAction(action Action, err error) (Action, error) {
	if err == nil {
		game.Actions = append(game.Actions, action)
	}
	return action, err
}

func (state *GameState) NewAction(actionType ActionType, playerPosition int, value int) Action {
	action := Action{
		ActionType:     actionType,
		PlayerPosition: playerPosition,
		Value:          value,
	}
	state.IncreaseStep()
	return action
}

func (state *GameState) IncreaseStep() {
	state.Step++
	state.CurrentPosition++
	if state.CurrentPosition/len(state.PlayerStates) == 1 {
		state.CurrentPosition = 0
		state.Round++
	}
}
