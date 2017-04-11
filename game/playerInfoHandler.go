package game

type PlayerGameInfo struct {
	MyTurn          bool               `json:"my_turn"`
	CurrentPostion  int                `json:"current_position"`
	PlayerCount     int                `json:"player_count"`
	Position        int                `json:"pos"`
	Step            int                `json:"step"`
	MaxStep         int                `json:"max_step"`
	Round           int                `json:"round"`
	PlayerId        int                `json:"player_id"`
	DeckSize        int                `json:"deck_size"`
	UsedCards       []Card             `json:"used_cards"`
	TableCards      map[CardColor]Card `json:"table_cards"`
	PlayerCards     [][]Card           `json:"player_cards"`
	PlayerCardsInfo [][]Card           `json:"player_cards_info"`
	BlueTokens      int                `json:"blue_tokens"`
	RedTokens       int                `json:"red_tokens"`
	Points          int                `json:"points"`
}

func (game *Game) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	return game.CurrentState.GetPlayerGameInfo(playerId)
}

func (state *GameState) GetPlayerGameInfoByPos(playerPosition int) PlayerGameInfo {
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

	for i := 0; i < len(playerCards); i++ {
		if i == playerPosition {
			continue
		}
		for j := 0; j < len(playerCards[i]); j++ {
			card := &playerCards[i][j]
			card.KnownColor = true
			card.KnownValue = true
			card.UpdateProbability()
		}
	}

	for i := 0; i < len(playerCardsInfo); i++ {
		for j := 0; j < len(playerCardsInfo[i]); j++ {
			card := &playerCardsInfo[i][j]
			if !card.KnownColor {
				(*card).Color = NoneColor
			}
			if !card.KnownValue {
				(*card).Value = NoneValue
			}
		}
	}

	return PlayerGameInfo{
		MyTurn:          state.CurrentPosition == playerPosition,
		CurrentPostion:  state.CurrentPosition,
		PlayerCount:     len(state.PlayerStates),
		Position:        playerPosition,
		Step:            state.Step,
		MaxStep:         state.MaxStep,
		Round:           state.Round,
		PlayerId:        state.PlayerStates[playerPosition].PlayerId,
		DeckSize:        len(state.Deck),
		UsedCards:       state.UsedCards,
		TableCards:      state.TableCards,
		PlayerCards:     playerCards,
		PlayerCardsInfo: playerCardsInfo,
		BlueTokens:      state.BlueTokens,
		RedTokens:       MaxRedTokens - state.RedTokens,
		Points:          0,
	}
}

func (info *PlayerGameInfo) GetPoints() int {
	if info.Points > 0 {
		return info.Points
	}

	for _, card := range info.TableCards {
		info.Points += int(card.Value)
	}
	return info.Points
}

func (state *GameState) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	var playerPosition int
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == playerId {
			playerPosition = i
		}
	}
	return state.GetPlayerGameInfoByPos(playerPosition)
}

func NewPlayerInfo(game *Game, playerId int) PlayerGameInfo {
	return game.GetPlayerGameInfo(playerId)
}

func (info *PlayerGameInfo) Copy() *PlayerGameInfo {
	newInfo := new(PlayerGameInfo)

	newInfo.MyTurn = info.MyTurn
	newInfo.CurrentPostion = info.CurrentPostion
	newInfo.PlayerCount = info.PlayerCount
	newInfo.Position = info.Position
	newInfo.Step = info.Step
	newInfo.MaxStep = info.MaxStep
	newInfo.Round = info.Round
	newInfo.PlayerId = info.PlayerId
	newInfo.DeckSize = info.DeckSize
	newInfo.Points = info.Points

	newInfo.UsedCards = make([]Card, len(info.UsedCards), cap(info.UsedCards))
	copy(newInfo.UsedCards, info.UsedCards)

	newInfo.TableCards = map[CardColor]Card{}
	for color, card := range info.TableCards {
		newInfo.TableCards[color] = card
	}

	newInfo.PlayerCards = make([][]Card, len(info.PlayerCards), len(info.PlayerCards))
	for i := 0; i < len(newInfo.PlayerCards); i++ {
		newInfo.PlayerCards[i] = make([]Card, len(info.PlayerCards[i]), len(info.PlayerCards[i]))
		copy(newInfo.PlayerCards[i], info.PlayerCards[i])
	}

	newInfo.PlayerCardsInfo = make([][]Card, len(info.PlayerCardsInfo), len(info.PlayerCardsInfo))
	for i := 0; i < len(newInfo.PlayerCardsInfo); i++ {
		newInfo.PlayerCardsInfo[i] = make([]Card, len(info.PlayerCardsInfo[i]), len(info.PlayerCardsInfo[i]))
		copy(newInfo.PlayerCardsInfo[i], info.PlayerCardsInfo[i])
	}

	newInfo.BlueTokens = info.BlueTokens
	newInfo.RedTokens = info.RedTokens
	return newInfo
}

func (info *PlayerGameInfo) IsGameOver() bool {
	if info.RedTokens == MaxRedTokens || info.MaxStep != 0 && info.Step >= info.MaxStep {
		return true
	}

	for _, card := range info.TableCards {
		if card.Value != Five {
			return false
		}
	}
	return true
}

func (info *PlayerGameInfo) IncreasePosition() {
	info.Step++
	info.CurrentPostion++
	if info.CurrentPostion/len(info.PlayerCards) == 1 {
		info.CurrentPostion = 0
		info.Round++
	}
}
