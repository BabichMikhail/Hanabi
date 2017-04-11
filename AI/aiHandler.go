package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

const (
	Type_AIRandom = iota
	Type_AISmartyRandom
	Type_AIDiscardUsefulCard
	Type_AIUsefulInformation
	Type_AIUsefulInformationV2
	Type_AIUsefulInformationV3
)

var AITypes = []int{
	Type_AIRandom,
	Type_AISmartyRandom,
	Type_AIDiscardUsefulCard,
	Type_AIUsefulInformation,
	Type_AIUsefulInformationV2,
	Type_AIUsefulInformationV3,
}

type Action struct {
	game.Action
	UsefullCount int `json:"usefull_count"`
	Count        int `json:"count"`
}

type BaseAI struct {
	Actions          []*Action           `json:"actions"`
	PlayActions      []*Action           `json:"playing_actions"`
	DiscardActions   []*Action           `json:"discard_actions"`
	InfoValueActions []*Action           `json:"info_value_actions"`
	InfoColorAcions  []*Action           `json:"info_color_actions"`
	History          []game.Action       `json:"history"`
	PlayerInfo       game.PlayerGameInfo `json:"player_info"`
	Type             int                 `json:"ai_type"`
	Informator       AIInformator        `json:"informator"`
	InfoIsSetted     bool                `json:"info_is_setted"`
}

type AI interface {
	GetAction() game.Action
}

const (
	AI_NamePrefix = "AI_"

	Name_AIRandom              = "Random"
	Name_AISmartyRandom        = "SmartyRandom"
	Name_AIDiscardUsefulCard   = "DiscardKnownCard"
	Name_AIUsefulInformation   = "UsefulInformation"
	Name_AIUsefulInformationV2 = "UsefulInformationV2"
	Name_AIUsefulInformationV3 = "UsefulInformationV3"
)

var AINames = map[int]string{
	Type_AIRandom:              Name_AIRandom,
	Type_AISmartyRandom:        Name_AISmartyRandom,
	Type_AIDiscardUsefulCard:   Name_AIDiscardUsefulCard,
	Type_AIUsefulInformation:   Name_AIUsefulInformation,
	Type_AIUsefulInformationV2: Name_AIUsefulInformationV2,
	Type_AIUsefulInformationV3: Name_AIUsefulInformationV3,
}

type AIInformator interface {
}

func NewAI(playerInfo game.PlayerGameInfo, history []game.Action, aiType int, informator AIInformator) AI {
	baseAI := new(BaseAI)
	baseAI.Informator = informator
	baseAI.History = history
	baseAI.PlayerInfo = playerInfo
	baseAI.setAvailableActions()
	baseAI.Type = aiType
	baseAI.InfoIsSetted = false

	var ai AI
	switch aiType {
	case Type_AIRandom:
		ai = NewAIRandom(baseAI)
	case Type_AISmartyRandom:
		ai = NewAISmartyRandom(baseAI)
	case Type_AIDiscardUsefulCard:
		ai = NewAIDiscardKnownCard(baseAI)
	case Type_AIUsefulInformation:
		ai = NewAIUsefulInformation(baseAI)
	case Type_AIUsefulInformationV2:
		ai = NewAIUsefulInformationV2(baseAI)
	case Type_AIUsefulInformationV3:
		ai = NewAIUsefulInformationV3(baseAI)
	default:
		panic("Unknown aiType")
	}
	return ai
}
