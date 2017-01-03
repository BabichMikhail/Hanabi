package engine

type PlayerGameInfo struct {
	PlayerCount  int          `json:"player_count"`
	Position     int          `json:"pos"`
	Step         int          `json:"step"`
	Round        int          `json:"round"`
	PlayerId     int          `json:"player_id"`
	DeckSize     int          `json:"deck_size"`
	UsedCards    []Card       `json:"used_cards"`
	TableCards   map[int]Card `json:"table_cards"`
	PlayersCards [][]Card     `json:"players_cards"`
	BlueTokens   int          `json:"blue_tokens"`
	RedTokens    int          `json:"red_tokens"`
}

func (this *Game) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	state := this.GameStatus[len(this.GameStatus)-1]
	var playerState PlayerState
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == playerId {
			playerState = state.PlayerStates[i]
		}
	}

	for i := 0; i < len(playerState.PlayersCards); i++ {
		if playerState.PlayerPosition == i {
			cards := &playerState.PlayersCards[i]
			for j := 0; j < len(*cards); j++ {
				card := &(*cards)[j]
				card.KnownColor = false
				card.Color = NoneColor
				card.KnownValue = false
				card.Value = NoneValue
			}
		}
	}

	return PlayerGameInfo{
		PlayerCount:  this.PlayerCount,
		Position:     playerState.PlayerPosition,
		Step:         state.Step,
		Round:        state.Round,
		PlayerId:     playerState.PlayerId,
		DeckSize:     len(state.Deck),
		UsedCards:    state.UsedCards,
		TableCards:   state.TableCards,
		PlayersCards: playerState.PlayersCards,
		BlueTokens:   state.BlueTokens,
		RedTokens:    state.RedTokens,
	}
}

func NewPlayerInfo(game *Game, playerId int) PlayerGameInfo {
	return game.GetPlayerGameInfo(playerId)
}
