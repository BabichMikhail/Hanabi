package ai

func (ai *BaseAI) setAvailableInformation() {
	if ai.InfoIsSetted {
		return
	}
	info := &ai.PlayerInfo
	info.SetProbabilities(ai.Type == Type_AICheater, ai.Type == Type_AIFullCheater)
	ai.InfoIsSetted = true
}
