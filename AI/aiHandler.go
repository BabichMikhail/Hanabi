package ai

import "github.com/BabichMikhail/Hanabi/game"

const (
	AI_RandomAction = iota
	AI_SmartyRandomAction
	AI_DiscardUsefullCardAction
)

type Action struct {
	game.Action
	UsefullCount int `json:"usefull_count"`
	Count        int `json:"count"`
}

type AI struct {
	Actions          []*Action           `json:"actions"`
	PlayActions      []*Action           `json:"playing_actions"`
	DiscardActions   []*Action           `json:"discard_actions"`
	InfoValueActions []*Action           `json:"info_value_actions"`
	InfoColorAcions  []*Action           `json:"info_color_actions"`
	PlayerInfo       game.PlayerGameInfo `json:"player_info"`
	Type             int                 `json:"ai_type"`
}

const AI_NamePrefix = "AI_"

func DefaultUsernamePrefix(AIType int) string {
	switch AIType {
	case AI_RandomAction:
		return AI_NamePrefix + "RandomAction"
	case AI_SmartyRandomAction:
		return AI_NamePrefix + "SmartyRandomAction"
	case AI_DiscardUsefullCardAction:
		return AI_NamePrefix + "DiscardKnownCardAction"
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
	var action *Action
	switch ai.Type {
	case AI_RandomAction:
		action = ai.getRandomAction()
	case AI_SmartyRandomAction:
		action = ai.getSmartyRandomAction()
	case AI_DiscardUsefullCardAction:
		action = ai.getDiscardUsefullCardAction()
	}
	return action.Action
}

func (ai *AI) SetAvailableInfomation() {
	info := &ai.PlayerInfo
	for idx, card := range info.PlayerCards[info.Position] {
		count := 0
		var color game.CardColor
		for cardColor, known := range card.AvailableColors {
			if known {
				count++
				color = cardColor
			}
		}
		if count == 1 {
			card := &info.PlayerCards[info.Position][idx]
			card.KnownColor = true
			card.Color = color
		}

		count = 0
		var value game.CardValue
		for cardValue, known := range card.AvailableValues {
			if known {
				count++
				value = cardValue
			}
		}
		if count == 1 {
			card := &info.PlayerCards[info.Position][idx]
			card.KnownValue = true
			card.Value = value
		}
	}
}
