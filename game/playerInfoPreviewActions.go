package game

type ResultPreviewPlayerInformations struct {
	Action  Action
	Results []ResultPlayerInfo
}

type ResultPlayerInfo struct {
	Probability float64
	Info        *PlayerGameInfo
}

func (info *PlayerGameInfo) PreviewActionDiscard(cardPosition int) ResultPreviewPlayerInformations {
	if info.BlueTokens == MaxBlueTokens {
		panic("Max blue tokens")
	}

	if info.DeckSize > 0 {
		panic("Not implemented")
	}

	newPlayerInfo := info.Copy()
	playerPosition := newPlayerInfo.CurrentPostion
	cards := newPlayerInfo.PlayerCards[playerPosition]
	if len(cards) >= cardPosition {
		panic("Card not found")
	}
	newPlayerInfo.PlayerCards[playerPosition] = append(cards[:cardPosition], cards[cardPosition+1:]...)

	cards = newPlayerInfo.PlayerCardsInfo[playerPosition]
	newPlayerInfo.PlayerCardsInfo[playerPosition] = append(cards[:cardPosition], cards[cardPosition+1:]...)

	return ResultPreviewPlayerInformations{
		Action: NewAction(TypeActionDiscard, playerPosition, cardPosition),
		Results: []ResultPlayerInfo{
			ResultPlayerInfo{
				Probability: 1.0,
				Info:        newPlayerInfo,
			},
		},
	}
}

func (info *PlayerGameInfo) PreviewActionPlaying(cardPosition int) ResultPreviewPlayerInformations {
	if info.DeckSize > 0 {
		panic("Not implemented")
	}

	newPlayerInfo := info.Copy()
	playerPosition := newPlayerInfo.CurrentPostion
	action := NewAction(TypeActionPlaying, playerPosition, cardPosition)

	updateFunc := func(playerInfo *PlayerGameInfo, cardValue CardValue, cardColor CardColor, probability float64) ResultPlayerInfo {
		newPlayerInfo = playerInfo.Copy()
		if newPlayerInfo.TableCards[cardColor].Value+1 == cardValue {
			newPlayerInfo.TableCards[cardColor] = *NewCard(cardColor, cardValue, true)
			playerCards := newPlayerInfo.PlayerCards[playerPosition]
			playerCardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
			newPlayerInfo.PlayerCards[playerPosition] = append(playerCards[0:cardPosition], playerCards[cardPosition+1:]...)
			newPlayerInfo.PlayerCardsInfo[playerPosition] = append(playerCardsInfo[0:cardPosition], playerCardsInfo[cardPosition+1:]...)
			if cardValue == Five && newPlayerInfo.BlueTokens < MaxBlueTokens {
				newPlayerInfo.BlueTokens++
			}
		} else {
			newPlayerInfo.RedTokens++
		}

		return ResultPlayerInfo{
			Probability: probability,
			Info:        newPlayerInfo,
		}
	}

	card := info.PlayerCards[playerPosition][cardPosition].Copy()

	if card.KnownColor && card.KnownValue {
		results := []ResultPlayerInfo{
			updateFunc(info, card.Value, card.Color, 1.0),
		}
		return ResultPreviewPlayerInformations{
			Action:  action,
			Results: results,
		}
	}

	if card.KnownColor {
		idx := 0
		results := make([]ResultPlayerInfo, len(card.ProbabilityValues), len(card.ProbabilityValues))
		for cardValue, probability := range card.ProbabilityValues {
			results[idx] = updateFunc(info, cardValue, card.Color, probability)
			idx++
		}
		return ResultPreviewPlayerInformations{
			Action:  action,
			Results: results,
		}
	}

	if card.KnownValue {
		results := make([]ResultPlayerInfo, len(card.ProbabilityColors), len(card.ProbabilityColors))
		idx := 0
		for cardColor, probability := range card.ProbabilityColors {
			results[idx] = updateFunc(info, card.Value, cardColor, probability)
			idx++
		}
		return ResultPreviewPlayerInformations{
			Action:  action,
			Results: results,
		}
	}

	idx := 0
	length := len(card.ProbabilityColors) * len(card.ProbabilityValues)
	results := make([]ResultPlayerInfo, length, length)
	for cardColor, probabilityColor := range card.ProbabilityColors {
		for cardValue, probabilityValue := range card.ProbabilityValues {
			results[idx] = updateFunc(info, cardValue, cardColor, probabilityColor*probabilityValue)
			idx++
		}
	}

	return ResultPreviewPlayerInformations{
		Action:  action,
		Results: results,
	}
}

func (info *PlayerGameInfo) PreviewActionInformationColor(playerPosition int, cardColor CardColor) ResultPreviewPlayerInformations {
	if info.BlueTokens == 0 {
		panic("No blue tokens")
	}

	newPlayerInfo := info.Copy()
	newPlayerInfo.BlueTokens--
	cards := newPlayerInfo.PlayerCards[playerPosition]
	cardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
	for i := 0; i < len(cards); i++ {
		if cards[i].Color == cardColor && !cardsInfo[i].KnownColor {
			cardsInfo[i].KnownColor = true
			cardsInfo[i].Color = cardColor
		}
	}

	return ResultPreviewPlayerInformations{
		Action: NewAction(TypeActionInformationColor, playerPosition, int(cardColor)),
		Results: []ResultPlayerInfo{
			ResultPlayerInfo{
				Probability: 1.0,
				Info:        newPlayerInfo,
			},
		},
	}
}

func (info *PlayerGameInfo) PreviewActionInformationValue(playerPosition int, cardValue CardValue) ResultPreviewPlayerInformations {
	if info.BlueTokens == 0 {
		panic("No blue tokens")
	}

	newPlayerInfo := info.Copy()
	newPlayerInfo.BlueTokens--
	cards := newPlayerInfo.PlayerCards[playerPosition]
	cardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
	for i := 0; i < len(cards); i++ {
		if cards[i].Value == cardValue && !cardsInfo[i].KnownValue {
			cardsInfo[i].KnownValue = true
			cardsInfo[i].Value = cardValue
		}
	}

	return ResultPreviewPlayerInformations{
		Action: NewAction(TypeActionInformationValue, playerPosition, int(cardValue)),
		Results: []ResultPlayerInfo{
			ResultPlayerInfo{
				Probability: 1.0,
				Info:        newPlayerInfo,
			},
		},
	}
}
