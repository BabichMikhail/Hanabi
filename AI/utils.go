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
