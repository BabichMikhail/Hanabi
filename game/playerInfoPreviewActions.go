package game

type ResultPreviewPlayerInformations struct {
	Action  Action
	Max     int
	Min     int
	Med     float64
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

	points := newPlayerInfo.GetPoints()
	return ResultPreviewPlayerInformations{
		Action: NewAction(TypeActionDiscard, playerPosition, cardPosition),
		Max:    points,
		Min:    points,
		Med:    float64(points),
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

	max := -1
	min := 26
	med := 0.0

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

		points := newPlayerInfo.GetPoints()
		med += float64(points) * probability
		if points > max {
			max = points
		}
		if points < min {
			min = points
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
		length := len(card.ProbabilityValues)
		results := make([]ResultPlayerInfo, length, length)
		for cardValue, probability := range card.ProbabilityValues {
			results[idx] = updateFunc(info, cardValue, card.Color, probability)
			idx++
		}

		return ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}
	}

	if card.KnownValue {
		idx := 0
		length := len(card.ProbabilityColors)
		results := make([]ResultPlayerInfo, length, length)
		for cardColor, probability := range card.ProbabilityColors {
			results[idx] = updateFunc(info, card.Value, cardColor, probability)
			idx++
		}

		return ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
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
		Max:     max,
		Min:     min,
		Med:     med,
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

	points := newPlayerInfo.GetPoints()
	return ResultPreviewPlayerInformations{
		Action: NewAction(TypeActionInformationColor, playerPosition, int(cardColor)),
		Max:    points,
		Min:    points,
		Med:    float64(points),
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

	points := newPlayerInfo.GetPoints()
	return ResultPreviewPlayerInformations{
		Action: NewAction(TypeActionInformationValue, playerPosition, int(cardValue)),
		Max:    points,
		Min:    points,
		Med:    float64(points),
		Results: []ResultPlayerInfo{
			ResultPlayerInfo{
				Probability: 1.0,
				Info:        newPlayerInfo,
			},
		},
	}
}
