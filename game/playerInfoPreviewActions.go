package game

import (
	"errors"
	"fmt"
	"strconv"
)

type ResultPreviewPlayerInformations struct {
	Action  *Action
	Max     int
	Min     int
	Med     float64
	Results []*ResultPlayerInfo
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
	playerPosition := newPlayerInfo.CurrentPosition
	cards := newPlayerInfo.PlayerCards[playerPosition]
	action := NewAction(TypeActionDiscard, playerPosition, cardPosition)
	if len(cards) <= cardPosition {
		return nil, errors.New("Card not found")
	}

	max := -1
	min := 26
	med := 0.0

	updateFunc := func(playerInfo *PlayerGameInfo, cardValue CardValue, cardColor CardColor, probability float64) []*ResultPlayerInfo {
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
			cardColor: 1.0,
		}
		card.ProbabilityValues = map[CardValue]float64{
			cardValue: 1.0,
		}
		card.ProbabilityCard = map[HashValue]float64{
			HashColorValue(cardColor, cardValue): 1.0,
		}
		colorValue := ColorValue{Color: cardColor, Value: cardValue}
		newPlayerInfo.VariantsCount[colorValue]--
		newPlayerInfo.UsedCards = append(newPlayerInfo.UsedCards, card)
		newPlayerInfo.BlueTokens++

		points := newPlayerInfo.GetPoints()
		med += float64(points) * probability
		if points > max {
			max = points
		}
		if points < min {
			min = points
		}

		newPlayerInfo.IncreasePosition()

		for i := 0; i < len(playerCards); i++ {
			playerCards[i].NormalizeProbabilities(cardColor, cardValue, newPlayerInfo.VariantsCount[colorValue])
		}
		for i := 0; i < len(newPlayerInfo.Deck); i++ {
			newPlayerInfo.Deck[i].NormalizeProbabilities(cardColor, cardValue, newPlayerInfo.VariantsCount[colorValue])
		}

		results := []*ResultPlayerInfo{}
		if newPlayerInfo.DeckSize > 0 {
			newPlayerInfo.MoveCardFromDeckToPlayer(playerPosition)
			cardPos := len(cards) - 1
			if cards[cardPos].KnownColor && cards[cardPos].KnownValue {
				results = append(results, &ResultPlayerInfo{
					Probability: probability,
					Info:        newPlayerInfo,
				})
			} else {
				for hashColorValue, prob := range cards[cardPos].ProbabilityCard {
					newInfo := newPlayerInfo.Copy()
					playerCards = newInfo.PlayerCards[playerPosition]
					card := &playerCards[cardPos]
					card.KnownValue = true
					card.KnownColor = true
					color, value := ColorValueByHashColorValue(hashColorValue)
					card.ProbabilityCard = map[HashValue]float64{hashColorValue: 1.0}
					card.ProbabilityColors = map[CardColor]float64{color: 1.0}
					card.ProbabilityValues = map[CardValue]float64{value: 1.0}
					colorValue := ColorValue{Color: color, Value: value}
					newInfo.VariantsCount[colorValue]--
					count := newInfo.VariantsCount[colorValue]
					for i := 0; i < len(playerCards)-1; i++ {
						playerCards[i].NormalizeProbabilities(color, value, count)
					}
					for i := 0; i < len(newInfo.Deck); i++ {
						newInfo.Deck[i].NormalizeProbabilities(color, value, count)
					}
					results = append(results, &ResultPlayerInfo{
						Probability: probability * prob,
						Info:        newInfo,
					})
				}
			}
		} else {
			results = append(results, &ResultPlayerInfo{
				Probability: probability,
				Info:        newPlayerInfo,
			})
		}

		return results
	}

	card := &cards[cardPosition]
	if card.KnownColor && card.KnownValue {
		results := updateFunc(info, card.Value, card.Color, 1.0)
		if len(results) == 0 {
			return nil, nil
		}
		return &ResultPreviewPlayerInformations{
			Action:  action,
			Max:     max,
			Min:     min,
			Med:     med,
			Results: results,
		}, nil
	}

	results := []*ResultPlayerInfo{}
	for colorValue, probability := range card.ProbabilityCard {
		cardColor, cardValue := ColorValueByHashColorValue(colorValue)
		newResults := updateFunc(info, cardValue, cardColor, probability)
		if len(newResults) != 0 {
			results = append(results, newResults...)
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
	playerPosition := newPlayerInfo.CurrentPosition
	action := NewAction(TypeActionPlaying, playerPosition, cardPosition)

	max := -1
	min := 26
	med := 0.0

	updateFunc := func(playerInfo *PlayerGameInfo, cardValue CardValue, cardColor CardColor, probability float64) []*ResultPlayerInfo {
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
			card.ProbabilityColors = map[CardColor]float64{cardColor: 1.0}
			card.ProbabilityValues = map[CardValue]float64{cardValue: 1.0}
			card.ProbabilityCard = map[HashValue]float64{
				HashColorValue(cardColor, cardValue): 1.0,
			}
			newPlayerInfo.UsedCards = append(newPlayerInfo.UsedCards, card)
			newPlayerInfo.RedTokens++
		}
		colorValue := ColorValue{Color: cardColor, Value: cardValue}
		newPlayerInfo.VariantsCount[colorValue]--

		points := newPlayerInfo.GetPoints()
		med += float64(points) * probability
		if points > max {
			max = points
		}
		if points < min {
			min = points
		}

		newPlayerInfo.IncreasePosition()

		for i := 0; i < len(playerCards); i++ {
			playerCards[i].NormalizeProbabilities(cardColor, cardValue, newPlayerInfo.VariantsCount[colorValue])
		}
		for i := 0; i < len(newPlayerInfo.Deck); i++ {
			newPlayerInfo.Deck[i].NormalizeProbabilities(cardColor, cardValue, newPlayerInfo.VariantsCount[colorValue])
		}

		results := []*ResultPlayerInfo{}
		if newPlayerInfo.DeckSize > 0 {
			newPlayerInfo.MoveCardFromDeckToPlayer(playerPosition)
			cards := newPlayerInfo.PlayerCards[playerPosition]
			cardPos := len(cards) - 1
			for hashColorValue, prob := range cards[cardPos].ProbabilityCard {
				newInfo := newPlayerInfo.Copy()
				playerCards = newInfo.PlayerCards[playerPosition]
				card := &playerCards[cardPos]
				card.KnownValue = true
				card.KnownColor = true
				color, value := ColorValueByHashColorValue(hashColorValue)
				card.ProbabilityCard = map[HashValue]float64{hashColorValue: 1.0}
				card.ProbabilityColors = map[CardColor]float64{color: 1.0}
				card.ProbabilityValues = map[CardValue]float64{value: 1.0}
				colorValue := ColorValue{Color: color, Value: value}

				for i := 0; i < len(playerCards); i++ {
					playerCards[i].NormalizeProbabilities(color, value, newInfo.VariantsCount[colorValue])
				}
				for i := 0; i < len(newInfo.Deck); i++ {
					newInfo.Deck[i].NormalizeProbabilities(color, value, newInfo.VariantsCount[colorValue])
				}
				results = append(results, &ResultPlayerInfo{
					Probability: probability * prob,
					Info:        newInfo,
				})
			}
		} else {
			results = append(results, &ResultPlayerInfo{
				Probability: probability,
				Info:        newPlayerInfo,
			})
		}

		return results
	}

	cards := info.PlayerCards[playerPosition]

	if len(cards) <= cardPosition {
		return nil, errors.New("Card not found")
	}

	card := cards[cardPosition].Copy()

	if card.KnownColor && card.KnownValue {
		results := updateFunc(info, card.Value, card.Color, 1.0)
		if len(results) == 0 {
			return nil, errors.New("Fail for optimize")
		}

		return &ResultPreviewPlayerInformations{
			Action:  action,
			Results: results,
		}, nil
	}

	results := []*ResultPlayerInfo{}
	for colorValue, probability := range card.ProbabilityCard {
		cardColor, cardValue := ColorValueByHashColorValue(colorValue)
		newResults := updateFunc(info, cardValue, cardColor, probability)
		if len(newResults) != 0 {
			results = append(results, newResults...)
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
			panic("error Color " + strconv.Itoa(playerPosition))
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
		Results: []*ResultPlayerInfo{
			&ResultPlayerInfo{
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
		Results: []*ResultPlayerInfo{
			&ResultPlayerInfo{
				Probability: 1.0,
				Info:        newPlayerInfo,
			},
		},
	}, nil
}
