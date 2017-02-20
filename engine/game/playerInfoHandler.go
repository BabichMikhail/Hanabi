package game

type PlayerGameInfo struct {
	MyTurn          bool               `json:"my_turn"`
	PlayerCount     int                `json:"player_count"`
	Position        int                `json:"pos"`
	Step            int                `json:"step"`
	Round           int                `json:"round"`
	PlayerId        int                `json:"player_id"`
	DeckSize        int                `json:"deck_size"`
	UsedCards       []Card             `json:"used_cards"`
	TableCards      map[CardColor]Card `json:"table_cards"`
	PlayerCards     [][]Card           `json:"player_cards"`
	PlayerCardsInfo [][]Card           `json:"player_cards_info"`
	BlueTokens      int                `json:"blue_tokens"`
	RedTokens       int                `json:"red_tokens"`
}

func (game *Game) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	return game.CurrentState.GetPlayerGameInfo(playerId)
}

func (state *GameState) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	var playerPosition int
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == playerId {
			playerPosition = i
		}
	}

	playerCardsInfo := [][]Card{}
	playerCards := [][]Card{}
	for i := 0; i < len(state.PlayerStates); i++ {
		cards1 := make([]Card, len(state.PlayerStates[i].PlayerCards))
		cards2 := make([]Card, len(state.PlayerStates[i].PlayerCards))
		copy(cards1, state.PlayerStates[i].PlayerCards)
		playerCardsInfo = append(playerCardsInfo, cards1)
		copy(cards2, state.PlayerStates[i].PlayerCards)
		playerCards = append(playerCards, cards2)
	}

	for i := 0; i < len(playerCards[playerPosition]); i++ {
		card := &playerCards[playerPosition][i]
		if !(*card).KnownColor {
			(*card).Color = NoneColor
		}
		if !(*card).KnownValue {
			(*card).Value = NoneValue
		}
	}

	for i := 0; i < len(playerCardsInfo); i++ {
		for j := 0; j < len(playerCardsInfo[i]); j++ {
			card := &playerCardsInfo[i][j]
			if !(*card).KnownColor {
				(*card).Color = NoneColor
			}
			if !(*card).KnownValue {
				(*card).Value = NoneValue
			}
		}
	}

	return PlayerGameInfo{
		MyTurn:          state.CurrentPosition == playerPosition,
		PlayerCount:     state.PlayerCount,
		Position:        playerPosition,
		Step:            state.Step,
		Round:           state.Round,
		PlayerId:        state.PlayerStates[playerPosition].PlayerId,
		DeckSize:        len(state.Deck),
		UsedCards:       state.UsedCards,
		TableCards:      state.TableCards,
		PlayerCards:     playerCards,
		PlayerCardsInfo: playerCardsInfo,
		BlueTokens:      state.BlueTokens,
		RedTokens:       MaxRedTokens - state.RedTokens,
	}
}

func NewPlayerInfo(game *Game, playerId int) PlayerGameInfo {
	return game.GetPlayerGameInfo(playerId)
}
