package game

import "fmt"

type PlayerGameInfo struct {
	MyTurn           bool               `json:"my_turn"`
	PlayerCount      int                `json:"player_count"`
	Position         int                `json:"pos"`
	Step             int                `json:"step"`
	Round            int                `json:"round"`
	PlayerId         int                `json:"player_id"`
	DeckSize         int                `json:"deck_size"`
	UsedCards        []Card             `json:"used_cards"`
	TableCards       map[CardColor]Card `json:"table_cards"`
	PlayersCards     [][]Card           `json:"players_cards"`
	PlayersCardsInfo [][]Card           `json:"players_cards_info"`
	BlueTokens       int                `json:"blue_tokens"`
	RedTokens        int                `json:"red_tokens"`
}

func (this *Game) GetPlayerGameInfo(playerId int) PlayerGameInfo {
	state := &this.CurrentState
	var playerState *PlayerState
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == playerId {
			playerState = &state.PlayerStates[i]
		}
	}

	playersCardsInfo := [][]Card{}
	fmt.Println(state.PlayerCount)
	for i := 0; i < state.PlayerCount; i++ {
		cards := []Card{}
		playerCards := state.PlayerStates[i].PlayersCards[i]
		for j := 0; j < len(playerCards); j++ {
			card := playerCards[j].Copy()
			if !card.KnownColor {
				card.Color = NoneColor
			}
			if !card.KnownValue {
				card.Value = NoneValue
			}
			cards = append(cards, card)
		}
		playersCardsInfo = append(playersCardsInfo, cards)
	}

	cards := playerState.PlayersCards[playerState.PlayerPosition]
	for j := 0; j < len(cards); j++ {
		card := &cards[j]
		if !card.KnownColor {
			card.Color = NoneColor
		}
		if !card.KnownValue {
			card.Value = NoneValue
		}
	}

	return PlayerGameInfo{
		MyTurn:           state.CurrentPosition == playerState.PlayerPosition,
		PlayerCount:      this.PlayerCount,
		Position:         playerState.PlayerPosition,
		Step:             state.Step,
		Round:            state.Round,
		PlayerId:         playerState.PlayerId,
		DeckSize:         len(state.Deck),
		UsedCards:        state.UsedCards,
		TableCards:       state.TableCards,
		PlayersCards:     playerState.PlayersCards,
		PlayersCardsInfo: playersCardsInfo,
		BlueTokens:       state.BlueTokens,
		RedTokens:        MaxRedTokens - state.RedTokens,
	}
}

func NewPlayerInfo(game *Game, playerId int) PlayerGameInfo {
	return game.GetPlayerGameInfo(playerId)
}
