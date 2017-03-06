package ai

import (
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

type Card struct {
	game.Card
	ProbabilityValues map[game.CardValue]float64
	ProbabilityColors map[game.CardColor]float64
}

func NewCard(gameCard game.Card) *Card {
	return &Card{gameCard, map[game.CardValue]float64{}, map[game.CardColor]float64{}}
}

const (
	AI_NamePrefix = "AI_"

	AI_RandomName             = "RandomAction"
	AI_SmartyName             = "SmartyRandomAction"
	AI_DiscardUsefullCardName = "DiscardKnownCardAction"
	AI_UsefullInformationName = "UsefullInformationAction"
)

func NewAI(playerInfo game.PlayerGameInfo, actions []game.Action, aiType int) *AI {
	ai := new(AI)
	ai.History = actions
	ai.PlayerInfo = playerInfo
	ai.setAvailableActions()
	ai.Type = aiType
	return ai
}

func (ai *AI) GetAction() game.Action {
	switch ai.Type {
	case AI_RandomAction:
		return ai.getActionRandom()
	case AI_SmartyRandomAction:
		return ai.getActionSmartyRandom()
	case AI_DiscardUsefullCardAction:
		return ai.getActionDiscardUsefullCard()
	case AI_UsefullInformationAction:
		return ai.getActionUsefullInformation()
	default:
		panic("Missing AI_Type")
	}
}

func (ai *AI) SetAvailableInfomation() {
	info := &ai.PlayerInfo
	for idx, card := range info.PlayerCards[info.Position] {
		if len(card.ProbabilityColors) == 1 {
			for color, _ := range card.ProbabilityColors {
				card := &info.PlayerCards[info.Position][idx]
				card.KnownColor = true
				card.Color = color
			}
		}

		if len(card.ProbabilityValues) == 1 {
			for value, _ := range card.ProbabilityValues {
				card := &info.PlayerCards[info.Position][idx]
				card.KnownValue = true
				card.Value = value
			}
		}
	}
}
