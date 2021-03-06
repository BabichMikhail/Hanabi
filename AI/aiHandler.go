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
	Type_AIUsefulInfoAndMaxMax
	Type_AIUsefulInfoAndMinMax
	Type_AIUsefulInfoAndMedMax
	Type_AIUsefulInfoV4AndParts
	Type_AIUsefulInfoV3AndParts
	Type_AIUsefulInformationV4
	Type_AICheater
	Type_AIFullCheater
	Type_AI1
	Type_AI2
	Type_AI6
	Type_AIHat
	Type_AI7
	Type_AI8
	Type_AI9
)

var AITypes = []int{
	Type_AIRandom,
	Type_AISmartyRandom,
	Type_AIDiscardUsefulCard,
	Type_AIUsefulInformation,
	Type_AIUsefulInformationV2,
	Type_AIUsefulInformationV3,
	Type_AIUsefulInfoAndMaxMax,
	Type_AIUsefulInfoAndMinMax,
	Type_AIUsefulInfoAndMedMax,
	Type_AIUsefulInfoV4AndParts,
	Type_AIUsefulInfoV3AndParts,
	Type_AIUsefulInformationV4,
	Type_AICheater,
	Type_AIFullCheater,
	Type_AI1,
	Type_AI2,
	Type_AI6,
	Type_AIHat,
	Type_AI7,
	Type_AI8,
	Type_AI9,
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
	GetAction() *game.Action
}

const (
	AI_NamePrefix = "AI_"

	Name_AIRandom               = "Random"
	Name_AISmartyRandom         = "SmartyRandom"
	Name_AIDiscardUsefulCard    = "DiscardKnownCard"
	Name_AIUsefulInformation    = "UsefulInformation"
	Name_AIUsefulInformationV2  = "UsefulInformationV2"
	Name_AIUsefulInformationV3  = "UsefulInformationV3"
	Name_AIUsefulInfoAndMaxMax  = "UsefulInfo&MaxMax"
	Name_AIUsefulInfoAndMinMax  = "UsefulInfo&MinMax"
	Name_AIUsefulInfoAndMedMax  = "UsefulInfo&MedMax"
	Name_AIUsefulInfoV4AndParts = "UsefulInfoV4AndParts"
	Name_AIUsefulInfoV3AndParts = "UsefulInfoV3AndParts"
	Name_AIUsefulInformationV4  = "UsefulInformationV4"
	Name_AICheater              = "Cheater"
	Name_AIFullCheater          = "FullCheater"
	Name_AI1                    = "AI1"
	Name_AI2                    = "AI2"
	Name_AI6                    = "AI6"
	Name_AIHat                  = "Hat"
	Name_AI7                    = "AI7"
	Name_AI8                    = "AI8"
	Name_AI9                    = "AI9"
)

var AINames = map[int]string{
	Type_AIRandom:               Name_AIRandom,
	Type_AISmartyRandom:         Name_AISmartyRandom,
	Type_AIDiscardUsefulCard:    Name_AIDiscardUsefulCard,
	Type_AIUsefulInformation:    Name_AIUsefulInformation,
	Type_AIUsefulInformationV2:  Name_AIUsefulInformationV2,
	Type_AIUsefulInformationV3:  Name_AIUsefulInformationV3,
	Type_AIUsefulInfoAndMaxMax:  Name_AIUsefulInfoAndMaxMax,
	Type_AIUsefulInfoAndMinMax:  Name_AIUsefulInfoAndMinMax,
	Type_AIUsefulInfoAndMedMax:  Name_AIUsefulInfoAndMedMax,
	Type_AIUsefulInfoV4AndParts: Name_AIUsefulInfoV4AndParts,
	Type_AIUsefulInfoV3AndParts: Name_AIUsefulInfoV3AndParts,
	Type_AIUsefulInformationV4:  Name_AIUsefulInformationV4,
	Type_AICheater:              Name_AICheater,
	Type_AIFullCheater:          Name_AIFullCheater,
	Type_AI1:                    Name_AI1,
	Type_AI2:                    Name_AI2,
	Type_AI6:                    Name_AI6,
	Type_AIHat:                  Name_AIHat, /* hat_player: https://github.com/chikinn/hanabi */
	Type_AI7:                    Name_AI7,
	Type_AI8:                    Name_AI8,
	Type_AI9:                    Name_AI9,
}

type AIInformator interface {
	GetPlayerState(step int) game.PlayerGameInfo
	GetQualitativeAssessmentOfState(*game.PlayerGameInfo) float64
	GetCache() interface{}
	SetCache(interface{})
	CheckAvailablePlayerInformation([]*game.AvailablePlayerGameInfo, int) int
	SetProbabilities(*game.PlayerGameInfo)
	PlayerInfoHash(*game.PlayerGameInfo) string
	ForcePlayerInfoHash(*game.PlayerGameInfo) string
	GetAction(*game.PlayerGameInfo, int, []game.Action) *game.Action
	SetPseudoPlayerInfo(*game.PlayerGameInfo)
}

func NewAI(playerInfo game.PlayerGameInfo, history []game.Action, aiType int, informator AIInformator) AI {
	baseAI := new(BaseAI)
	baseAI.Informator = informator
	baseAI.History = history
	baseAI.PlayerInfo = playerInfo
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
		ai = NewAIUsefulInfoV3AndParts(baseAI, true)
	case Type_AIUsefulInfoAndMaxMax:
		ai = NewAIUsefulInfoAndMaxMax(baseAI)
	case Type_AIUsefulInfoAndMinMax:
		ai = NewAIUsefulInfoAndMinMax(baseAI)
	case Type_AIUsefulInfoAndMedMax:
		ai = NewAIUsefulInfoAndMedMax(baseAI)
	case Type_AIUsefulInfoV4AndParts:
		ai = NewAIUsefulInfoV4AndParts(baseAI, false)
	case Type_AIUsefulInfoV3AndParts:
		ai = NewAIUsefulInfoV3AndParts(baseAI, false)
	case Type_AIUsefulInformationV4:
		ai = NewAIUsefulInfoV4AndParts(baseAI, true)
	case Type_AICheater:
		ai = NewAICheater(baseAI)
	case Type_AI1:
		ai = NewAI1(baseAI)
	case Type_AI2:
		ai = NewAI2(baseAI)
	case Type_AI6:
		ai = NewAI6(baseAI)
	case Type_AIHat:
		ai = NewAIHat(baseAI)
	case Type_AI7:
		ai = NewAI7(baseAI)
	case Type_AI8:
		ai = NewAI8(baseAI)
	case Type_AI9:
		ai = NewAI9(baseAI)
	default:
		panic("Unknown aiType")
	}
	return ai
}
