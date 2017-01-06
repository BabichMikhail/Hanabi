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

func (this *Game) NewAction(actionType ActionType, playerPosition int, value int) {
	action := Action{
		ActionType:     actionType,
		PlayerPosition: playerPosition,
		Value:          value,
	}
	this.Actions = append(this.Actions, action)
	this.IncreaseStep()
}

func (this *Game) IncreaseStep() {
	state := &this.CurrentState
	state.Step++
	state.CurrentPosition++
	if state.CurrentPosition/state.PlayerCount == 1 {
		state.CurrentPosition = 0
		state.Round++
	}
}
