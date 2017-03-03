package ai

import "github.com/BabichMikhail/Hanabi/game"

type key struct {
	Value    int
	Position int
}

type value struct {
	Count        int
	UsefullCount int
}

func updateValue(value value, card game.Card) value {
	value.Count++
	if !card.KnownValue {
		value.UsefullCount++
	}
	return value
}

func (ai *AI) AppendAction(actionType game.ActionType, playerPosition, actionValue, cardCount, usefullCardCount int) {
	action := &Action{
		Action: game.Action{
			ActionType:     actionType,
			PlayerPosition: playerPosition,
			Value:          actionValue,
		},
		Count:        cardCount,
		UsefullCount: usefullCardCount,
	}
	ai.Actions = append(ai.Actions, action)
	switch actionType {
	case game.TypeActionDiscard:
		ai.DiscardActions = append(ai.DiscardActions, action)
	case game.TypeActionPlaying:
		ai.PlayActions = append(ai.PlayActions, action)
	case game.TypeActionInformationColor:
		ai.InfoColorAcions = append(ai.InfoColorAcions, action)
	case game.TypeActionInformationValue:
		ai.InfoValueActions = append(ai.InfoValueActions, action)
	}
}

func getParams(actionType game.ActionType) (game.ActionType, func(*game.Card) int) {
	if actionType == game.TypeActionInformationColor {
		return actionType, func(card *game.Card) int {
			return int(card.Color)
		}
	} else {
		return actionType, func(card *game.Card) int {
			return int(card.Value)
		}
	}
}

func (ai *AI) setAvailableInfomationActions() {
	actionTypes := []game.ActionType{game.TypeActionInformationColor, game.TypeActionInformationValue}
	if ai.PlayerInfo.BlueTokens == 0 {
		return
	}
	for _, actionType := range actionTypes {
		actionType, cardF := getParams(actionType)
		values := map[key]struct {
			Count        int
			UsefullCount int
		}{}
		playerInfo := &ai.PlayerInfo
		playersCardInfo := playerInfo.PlayerCardsInfo
		for i, cards := range playerInfo.PlayerCards {
			if i == playerInfo.Position {
				continue
			}
			for j, card := range cards {
				values[key{cardF(&card), i}] = updateValue(values[key{cardF(&card), i}], playersCardInfo[i][j])
			}
		}

		for key, value := range values {
			ai.AppendAction(actionType, key.Position, key.Value, value.Count, value.UsefullCount)
		}
	}
}

func (ai *AI) setAvailablePlayingAndDiscardActions() {
	playerInfo := &ai.PlayerInfo
	pos := playerInfo.Position
	for i, _ := range playerInfo.PlayerCards[pos] {
		if ai.PlayerInfo.BlueTokens > 0 {
			ai.AppendAction(game.TypeActionPlaying, pos, i, 1, 1)
		}

		if ai.PlayerInfo.BlueTokens < game.MaxBlueTokens {
			ai.AppendAction(game.TypeActionDiscard, pos, i, 1, 1)
		}
	}
}

func (ai *AI) setAvailableActions() {
	ai.setAvailableInfomationActions()
	ai.setAvailablePlayingAndDiscardActions()
}
