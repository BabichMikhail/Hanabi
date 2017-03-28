package ai

import (
	"regexp"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func DefaultUsernamePrefix(AIType int) string {
	switch AIType {
	case AI_RandomAction:
		return AI_NamePrefix + AI_RandomName
	case AI_SmartyRandomAction:
		return AI_NamePrefix + AI_SmartyName
	case AI_DiscardUsefulCardAction:
		return AI_NamePrefix + AI_DiscardUsefulCardName
	case AI_UsefulInformationAction:
		return AI_NamePrefix + AI_UsefulInformationName
	default:
		return AI_NamePrefix + "Any"
	}
}

func GetAITypeByUserNickName(nickname string) int {
	if ok, _ := regexp.MatchString(AI_NamePrefix+AI_RandomName+"_\\d", nickname); ok {
		return AI_RandomAction
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+AI_SmartyName+"_\\d", nickname); ok {
		return AI_SmartyRandomAction
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+AI_DiscardUsefulCardName+"_\\d", nickname); ok {
		return AI_DiscardUsefulCardAction
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+AI_UsefulInformationName+"_\\d", nickname); ok {
		return AI_UsefulInformationAction
	}
	return AI_UsefulInformationAction
}
