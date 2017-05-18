package game

import "reflect"

func (availableInfo AvailablePlayerGameInfos) NormalizeProbablities() {
	probSum := 0.0
	for i, _ := range availableInfo {
		probSum += availableInfo[i].Probability
	}
	for i, _ := range availableInfo {
		availableInfo[i].Probability /= probSum
	}
}

func (card *Card) EmptyProbabilities() {
	card.ProbabilityCard = map[HashValue]float64{}
}

func (card *Card) ApplyProbabilities() {
	probSum := 0.0
	for _, prob := range card.ProbabilityCard {
		probSum += prob
	}
	for colorValue, _ := range card.ProbabilityCard {
		card.ProbabilityCard[colorValue] /= probSum
	}

	probSum = 0.0
	for _, prob := range card.ProbabilityColors {
		probSum += prob
	}
	for color, _ := range card.ProbabilityColors {
		card.ProbabilityColors[color] /= probSum
	}

	probSum = 0.0
	for _, prob := range card.ProbabilityValues {
		probSum += prob
	}
	for value, _ := range card.ProbabilityValues {
		card.ProbabilityValues[value] /= probSum
	}

	if len(card.ProbabilityValues) == 1 {
		keys := reflect.ValueOf(card.ProbabilityValues).MapKeys()
		value := keys[0].Interface().(CardValue)
		card.KnownValue = true
		card.Value = value
	}

	if len(card.ProbabilityColors) == 1 {
		keys := reflect.ValueOf(card.ProbabilityColors).MapKeys()
		color := keys[0].Interface().(CardColor)
		card.KnownColor = true
		card.Color = color
	}
}

func (availableInfo AvailablePlayerGameInfos) UpdatePlayerInformation(playerInfo *PlayerGameInfo) {
	availableInfo.NormalizeProbablities()
	for _, cards := range playerInfo.PlayerCards {
		for i := 0; i < len(cards); i++ {
			cards[i].EmptyProbabilities()
		}
	}

	for i := 0; i < len(playerInfo.Deck); i++ {
		playerInfo.Deck[i].EmptyProbabilities()
	}

	playerCards := playerInfo.PlayerCards
	for _, information := range availableInfo {
		info := information.PlayerInfo
		probability := information.Probability
		for pos, cards := range info.PlayerCards {
			for i := 0; i < len(cards); i++ {
				card := &cards[i]
				playerCard := &playerCards[pos][i]
				for colorValue, prob := range card.ProbabilityCard {
					playerCard.ProbabilityCard[colorValue] += prob * probability
				}
				for color, prob := range card.ProbabilityColors {
					playerCard.ProbabilityColors[color] += prob * probability
				}
				for value, prob := range card.ProbabilityValues {
					playerCard.ProbabilityValues[value] += prob * probability
				}
			}
		}

		for i, _ := range info.Deck {
			card := &info.Deck[i]
			playerCard := &playerInfo.Deck[i]
			for colorValue, prob := range card.ProbabilityCard {
				playerCard.ProbabilityCard[colorValue] += prob * probability
			}
			for color, prob := range card.ProbabilityColors {
				playerCard.ProbabilityColors[color] += prob * probability
			}
			for value, prob := range card.ProbabilityValues {
				playerCard.ProbabilityValues[value] += prob * probability
			}
		}
	}

	for _, cards := range playerInfo.PlayerCards {
		for i := 0; i < len(cards); i++ {
			card := &cards[i]
			card.ApplyProbabilities()
		}
	}

	for i := 0; i < len(playerInfo.Deck); i++ {
		playerInfo.Deck[i].ApplyProbabilities()
	}
	playerInfo.InfoIsSetted = true
}
