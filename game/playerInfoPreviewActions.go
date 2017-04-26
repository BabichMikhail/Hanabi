package game

import (
	"errors"
	"fmt"
)

type ResultPreviewPlayerInformations struct {
	Action  *Action
	Max     int
	Min     int
	Med     float64
	Results []ResultPlayerInfo
}

type ResultPlayerInfo struct {
	Probability float64
	Info        *PlayerGameInfo
}

func (info *PlayerGameInfo) MoveCardFromDeckToPlayer(playerPosition int) {
	card := info.Deck[0]
	info.Deck = info.Deck[1:]
	info.DeckSize--
	info.PlayerCards[playerPosition] = append(info.PlayerCards[playerPosition], card.Copy())
	info.PlayerCardsInfo[playerPosition] = append(info.PlayerCardsInfo[playerPosition], card.Copy())
}

func (info *PlayerGameInfo) PreviewAction(action *Action) (*ResultPreviewPlayerInformations, error) {
	switch action.ActionType {
	case TypeActionDiscard:
		return info.PreviewActionDiscard(action.Value)
	case TypeActionInformationColor:
		return info.PreviewActionInformationColor(action.PlayerPosition, CardColor(action.Value))
	case TypeActionInformationValue:
		return info.PreviewActionInformationValue(action.PlayerPosition, CardValue(action.Value))
	case TypeActionPlaying:
		return info.PreviewActionPlaying(action.Value)
	}
	panic("Unknown ActionType")
}

func (info *PlayerGameInfo) PreviewActionDiscard(cardPosition int) (*ResultPreviewPlayerInformations, error) {
	if info.BlueTokens == MaxBlueTokens {
		return nil, errors.New("Max blue tokens")
	}

	newPlayerInfo := info.Copy()
	playerPosition := newPlayerInfo.CurrentPostion
	cards := newPlayerInfo.PlayerCards[playerPosition]
	action := NewAction(TypeActionDiscard, playerPosition, cardPosition)
	if len(cards) <= cardPosition {
		return nil, errors.New("Card not found")
	}

	max := -1
	min := 26
	med := 0.0

	updateFunc := func(playerInfo *PlayerGameInfo, cardValue CardValue, cardColor CardColor, probability float64) *ResultPlayerInfo {
		newPlayerInfo = playerInfo.Copy()
		playerCards := newPlayerInfo.PlayerCards[playerPosition]
		card := playerCards[cardPosition].Copy()
		newPlayerInfo.PlayerCards[playerPosition] = append(playerCards[:cardPosition], playerCards[cardPosition+1:]...)

		playerCardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
		newPlayerInfo.PlayerCardsInfo[playerPosition] = append(playerCardsInfo[:cardPosition], playerCardsInfo[cardPosition+1:]...)

		card.KnownColor = true
		card.KnownValue = true
		card.Value = cardValue
		card.Color = cardColor
		card.ProbabilityColors = map[CardColor]float64{
			cardColor: 0.0,
		}
		card.ProbabilityValues = map[CardValue]float64{
			cardValue: 0.0,
		}
		newPlayerInfo.UsedCards = append(newPlayerInfo.UsedCards, card)
		newPlayerInfo.BlueTokens++

		if newPlayerInfo.DeckSize > 0 {
			newPlayerInfo.MoveCardFromDeckToPlayer(playerPosition)
		}

		points := newPlayerInfo.GetPoints()
		med += float64(points) * probability
		if points > max {
			max = points
		}
		if points < min {
			min = points
		}

		newPlayerInfo.IncreasePosition()
		return &ResultPlayerInfo{
			Probability: probability,
			Info:        newPlayerInfo,
		}
	}

	card := &cards[cardPosition]
	if card.KnownColor && card.KnownValue {
		result := updateFunc(info, card.Value, card.Color, 1.0)
		if result == nil {
			return nil, nil
		}
		results := []ResultPlayerInfo{
			*result,
		}
		return &ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}, nil
	}

	if card.KnownColor {
		results := []ResultPlayerInfo{}
		for cardValue, probability := range card.ProbabilityValues {
			result := updateFunc(info, cardValue, card.Color, probability)
			if result != nil {
				results = append(results, *result)
			}
		}

		return &ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}, nil
	}

	if card.KnownValue {
		results := []ResultPlayerInfo{}
		for cardColor, probability := range card.ProbabilityColors {
			result := updateFunc(info, card.Value, cardColor, probability)
			if result != nil {
				results = append(results, *result)
			}
		}

		return &ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}, nil
	}

	results := []ResultPlayerInfo{}
	for colorValue, probability := range card.ProbabilityCard {
		cardColor, cardValue := ColorValueByHashColorValue(colorValue)
		result := updateFunc(info, cardValue, cardColor, probability)
		if result != nil {
			results = append(results, *result)
		}
	}

	return &ResultPreviewPlayerInformations{
		Action:  action,
		Max:     max,
		Min:     min,
		Med:     med,
		Results: results,
	}, nil
}

func (info *PlayerGameInfo) PreviewActionPlaying(cardPosition int) (*ResultPreviewPlayerInformations, error) {
	if info.IsGameOver() {
		panic("GameOver")
	}

	newPlayerInfo := info.Copy()
	playerPosition := newPlayerInfo.CurrentPostion
	action := NewAction(TypeActionPlaying, playerPosition, cardPosition)

	max := -1
	min := 26
	med := 0.0

	updateFunc := func(playerInfo *PlayerGameInfo, cardValue CardValue, cardColor CardColor, probability float64) *ResultPlayerInfo {
		newPlayerInfo = playerInfo.Copy()
		playerCards := newPlayerInfo.PlayerCards[playerPosition]
		card := playerCards[cardPosition].Copy()
		newPlayerInfo.PlayerCards[playerPosition] = append(playerCards[:cardPosition], playerCards[cardPosition+1:]...)

		playerCardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
		newPlayerInfo.PlayerCardsInfo[playerPosition] = append(playerCardsInfo[:cardPosition], playerCardsInfo[cardPosition+1:]...)

		if newPlayerInfo.TableCards[cardColor].Value+1 == cardValue {
			newPlayerInfo.TableCards[cardColor] = *NewCard(cardColor, cardValue, true)
			if cardValue == Five && newPlayerInfo.BlueTokens < MaxBlueTokens {
				newPlayerInfo.BlueTokens++
			}
		} else {
			card.KnownColor = true
			card.KnownValue = true
			card.Value = cardValue
			card.Color = cardColor
			card.ProbabilityColors = map[CardColor]float64{
				cardColor: 1.0,
			}
			card.ProbabilityValues = map[CardValue]float64{
				cardValue: 1.0,
			}
			newPlayerInfo.UsedCards = append(newPlayerInfo.UsedCards, card)
			newPlayerInfo.RedTokens++
		}

		if newPlayerInfo.DeckSize > 0 {
			newPlayerInfo.MoveCardFromDeckToPlayer(playerPosition)
		}

		points := newPlayerInfo.GetPoints()
		med += float64(points) * probability
		if points > max {
			max = points
		}
		if points < min {
			min = points
		}

		newPlayerInfo.IncreasePosition()
		return &ResultPlayerInfo{
			Probability: probability,
			Info:        newPlayerInfo,
		}
	}

	cards := info.PlayerCards[playerPosition]

	if len(cards) <= cardPosition {
		return nil, errors.New("Card not found")
	}

	card := cards[cardPosition].Copy()

	if card.KnownColor && card.KnownValue {
		result := updateFunc(info, card.Value, card.Color, 1.0)
		if result == nil {
			return nil, errors.New("Fail for optimize")
		}
		results := []ResultPlayerInfo{
			*result,
		}
		return &ResultPreviewPlayerInformations{
			Action:  action,
			Results: results,
		}, nil
	}

	if card.KnownColor {
		results := []ResultPlayerInfo{}
		for colorValue, probability := range card.ProbabilityCard {
			cardColor, cardValue := ColorValueByHashColorValue(colorValue)
			result := updateFunc(info, cardValue, cardColor, probability)
			if result != nil {
				results = append(results, *result)
			}
		}

		return &ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}, nil
	}

	if card.KnownValue {
		results := []ResultPlayerInfo{}
		for colorValue, probability := range card.ProbabilityCard {
			cardColor, cardValue := ColorValueByHashColorValue(colorValue)
			result := updateFunc(info, cardValue, cardColor, probability)
			if result != nil {
				results = append(results, *result)
			}
		}

		return &ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}, nil
	}

	results := []ResultPlayerInfo{}
	for colorValue, probability := range card.ProbabilityCard {
		cardColor, cardValue := ColorValueByHashColorValue(colorValue)
		result := updateFunc(info, cardValue, cardColor, probability)
		if result != nil {
			results = append(results, *result)
		}
	}

	return &ResultPreviewPlayerInformations{
		Action:  action,
		Max:     max,
		Min:     min,
		Med:     med,
		Results: results,
	}, nil
}

func (info *PlayerGameInfo) PreviewActionInformationColor(playerPosition int, cardColor CardColor) (*ResultPreviewPlayerInformations, error) {
	if info.BlueTokens == 0 {
		return nil, errors.New("No blue tokens")
	}

	if info.VariantsCount == nil {
		return nil, errors.New("Need setAvailableInformation()")
	}

	newPlayerInfo := info.Copy()
	newPlayerInfo.BlueTokens--
	cards := newPlayerInfo.PlayerCards[playerPosition]
	cardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
	for i := 0; i < len(cards); i++ {
		if !cards[i].KnownColor {
			fmt.Println(cards[i])
			panic("error Color")
		}
		if cards[i].Color == cardColor && !cardsInfo[i].KnownColor {
			cardsInfo[i].KnownColor = true
			cardsInfo[i].Color = cardColor
		}
	}

	points := newPlayerInfo.GetPoints()
	newPlayerInfo.IncreasePosition()
	return &ResultPreviewPlayerInformations{
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
	}, nil
}

func (info *PlayerGameInfo) PreviewActionInformationValue(playerPosition int, cardValue CardValue) (*ResultPreviewPlayerInformations, error) {
	if info.BlueTokens == 0 {
		return nil, errors.New("No blue tokens")
	}

	if info.VariantsCount == nil {
		return nil, errors.New("Need setAvailableInformation()")
	}

	newPlayerInfo := info.Copy()
	newPlayerInfo.BlueTokens--
	cards := newPlayerInfo.PlayerCards[playerPosition]
	cardsInfo := newPlayerInfo.PlayerCardsInfo[playerPosition]
	for i := 0; i < len(cards); i++ {
		if !cards[i].KnownValue {
			fmt.Println(cards[i])
			panic("error Value")
		}
		if cards[i].Value == cardValue && !cardsInfo[i].KnownValue {
			cardsInfo[i].KnownValue = true
			cardsInfo[i].Value = cardValue
		}
	}

	points := newPlayerInfo.GetPoints()
	newPlayerInfo.IncreasePosition()
	return &ResultPreviewPlayerInformations{
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
	}, nil
}
