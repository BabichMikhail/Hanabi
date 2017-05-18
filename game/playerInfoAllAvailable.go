package game

import "fmt"

type AvailablePlayerGameInfos []*AvailablePlayerGameInfo

type AvailablePlayerGameInfo struct {
	PlayerInfo  *PlayerGameInfo
	Probability float64
}

func (info *PlayerGameInfo) setPlayerCardOnPosition(color CardColor, value CardValue, cardPos int) {
	myPos := info.Position
	colorValue := ColorValue{Color: color, Value: value}
	info.VariantsCount[colorValue]--
	count := float64(info.VariantsCount[colorValue])

	card := &info.PlayerCards[myPos][cardPos]
	card.ProbabilityCard = map[HashValue]float64{
		HashColorValue(color, value): 1.0,
	}
	card.ProbabilityColors = map[CardColor]float64{color: 1.0}
	card.ProbabilityValues = map[CardValue]float64{value: 1.0}
	card.Color = color
	card.Value = value
	card.KnownColor = true
	card.KnownValue = true

	for i := cardPos + 1; i < len(info.PlayerCards[myPos]); i++ {
		card := &info.PlayerCards[myPos][i]
		card.NormalizeProbabilities(color, value, int(count))
	}
	for i := 0; i < len(info.Deck); i++ {
		info.Deck[i].NormalizeProbabilities(color, value, int(count))
	}
}

func (info *PlayerGameInfo) AvailablePlayerInformations() []*AvailablePlayerGameInfo {
	myPos := info.Position
	results := []*AvailablePlayerGameInfo{
		&AvailablePlayerGameInfo{
			PlayerInfo:  info.Copy(),
			Probability: 1.0,
		},
	}

	for i := 0; i < len(info.PlayerCards[myPos]); i++ {
		length := len(results)
		fmt.Println(length)
		for j := 0; j < length; j++ {
			playerInfo := results[j].PlayerInfo
			probability := results[j].Probability
			myCards := playerInfo.PlayerCards[myPos]
			card := &myCards[i]
			if card.KnownColor && card.KnownValue {
				results = append(results, &AvailablePlayerGameInfo{
					PlayerInfo:  playerInfo,
					Probability: probability,
				})
				continue
			}

			fmt.Println(card.KnownColor, card.KnownValue)
			if len(card.ProbabilityCard) == 0 {
				panic("Need set probablities")
			}

			for colorValue, probabilityCard := range card.ProbabilityCard {
				color, value := ColorValueByHashColorValue(colorValue)
				copyPlayerInfo := playerInfo.Copy()
				copyPlayerInfo.setPlayerCardOnPosition(color, value, i)
				fmt.Println("append", probabilityCard)
				results = append(results, &AvailablePlayerGameInfo{
					PlayerInfo:  copyPlayerInfo,
					Probability: probability * probabilityCard,
				})
			}
		}
		results = results[length:]
	}
	return results
}
