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
	case Type_AIRandom:
		return AI_NamePrefix + Name_AIRandom
	case Type_AISmartyRandom:
		return AI_NamePrefix + Name_AISmartyRandom
	case Type_AIDiscardUsefulCard:
		return AI_NamePrefix + Name_AIDiscardUsefulCard
	case Type_AIUsefulInformation:
		return AI_NamePrefix + Name_AIUsefulInformation
	case Type_AIUsefulInformationV2:
		return AI_NamePrefix + Name_AIUsefulInformationV2
	}
	panic("Bad AI_Type")
}

func GetAITypeByUserNickName(nickname string) int {
	if ok, _ := regexp.MatchString(AI_NamePrefix+Name_AIRandom+"_\\d", nickname); ok {
		return Type_AIRandom
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+Name_AISmartyRandom+"_\\d", nickname); ok {
		return Type_AISmartyRandom
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+Name_AIDiscardUsefulCard+"_\\d", nickname); ok {
		return Type_AIDiscardUsefulCard
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+Name_AIUsefulInformation+"_\\d", nickname); ok {
		return Type_AIUsefulInformation
	} else if ok, _ := regexp.MatchString(AI_NamePrefix+Name_AIUsefulInformationV2+"_\\d", nickname); ok {
		return Type_AIUsefulInformationV2
	}
	panic("Bad nickname")
}
