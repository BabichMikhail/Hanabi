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
	for aiType, aiName := range AINames {
		if aiType == AIType {
			return AI_NamePrefix + aiName
		}
	}
	panic("Bad AIType")
}

func GetAITypeByUserNickName(nickname string) int {
	for aiType, aiName := range AINames {
		if ok, _ := regexp.MatchString(AI_NamePrefix+aiName+"_\\d", nickname); ok {
			return aiType
		}
	}
	panic("Bad nickname")
}
