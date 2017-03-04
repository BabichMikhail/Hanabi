package ai

import (
	"regexp"

	"github.com/BabichMikhail/Hanabi/game"
)

const (
	AI_RandomAction = iota
	AI_SmartyRandomAction
	AI_DiscardUsefullCardAction
	AI_UsefullInformationAction
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
	History          []game.Action       `json:"history"`
	PlayerInfo       game.PlayerGameInfo `json:"player_info"`
	Type             int                 `json:"ai_type"`
}

const (
	AI_NamePrefix = "AI_"

	AI_RandomName             = "RandomAction"
	AI_SmartyName             = "SmartyRandomAction"
	AI_DiscardUsefullCardName = "DiscardKnownCardAction"
	AI_UsefullInformationName = "UsefullInformationAction"
)

func DefaultUsernamePrefix(AIType int) string {
	switch AIType {
	case AI_RandomAction:
		return AI_NamePrefix + AI_RandomName
	case AI_SmartyRandomAction:
		return AI_NamePrefix + AI_SmartyName
	case AI_DiscardUsefullCardAction:
		return AI_NamePrefix + AI_DiscardUsefullCardName
	case AI_UsefullInformationAction:
		return AI_NamePrefix + AI_UsefullInformationName
	default:
		return AI_NamePrefix + "Any"
	}
}

func GetAITypeByUserNickName(nickname string) int {
	if ok, _ := regexp.MatchString(AI_NamePrefix+AI_RandomName+"_\\d", nickname); ok {
		return AI_RandomAction
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+AI_SmartyName+"_\\d", nickname); ok {
		return AI_SmartyRandomAction
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+AI_DiscardUsefullCardName+"_\\d", nickname); ok {
		return AI_DiscardUsefullCardAction
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+AI_UsefullInformationName+"_\\d", nickname); ok {
		return AI_UsefullInformationAction
	}
	return AI_UsefullInformationAction
}

func NewAI(playerInfo game.PlayerGameInfo, actions []game.Action, aiType int) *AI {
	ai := new(AI)
	ai.History = actions
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
	case AI_UsefullInformationAction:
		action = ai.getUsefullInformationAction()
	}
	return action.Action
}

func (ai *AI) SetAvailableInfomation() {
	info := &ai.PlayerInfo
	for idx, card := range info.PlayerCards[info.Position] {
		if len(card.AvailableColors) == 1 {
			for color, _ := range card.AvailableColors {
				card := &info.PlayerCards[info.Position][idx]
				card.KnownColor = true
				card.Color = color
			}
		}

		if len(card.AvailableValues) == 1 {
			for value, _ := range card.AvailableValues {
				card := &info.PlayerCards[info.Position][idx]
				card.KnownValue = true
				card.Value = value
			}
		}
	}
}
