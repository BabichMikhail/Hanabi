package ai

import "github.com/BabichMikhail/Hanabi/game"

const (
	AI_RandomAction = iota
)

type Action struct {
	game.Action
	UsefullCount int `json:"usefull_count"`
	Count        int `json:"count"`
}

type AI struct {
	Actions    []Action            `json:"actions"`
	PlayerInfo game.PlayerGameInfo `json:"player_info"`
	Type       int                 `json:"ai_type"`
}

const AI_NamePrefix = "AI_"

func DefaultUsernamePrefix(AIType int) string {
	switch AIType {
	case AI_RandomAction:
		return AI_NamePrefix + "RandomAction"
	default:
		return AI_NamePrefix + "Any"
	}
}

func NewAI(playerInfo game.PlayerGameInfo, aiType int) *AI {
	ai := new(AI)
	ai.PlayerInfo = playerInfo
	ai.setAvailableActions()
	ai.Type = aiType
	return ai
}

func (ai *AI) GetAction() game.Action {
	idx := 0
	switch ai.Type {
	case AI_RandomAction:
		idx = ai.getRandomActionIdx()
	}
	return ai.Actions[idx].Action
}