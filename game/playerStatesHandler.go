package game

type PlayerState struct {
	PlayerId       int    `json:"player_id"`
	PlayerPosition int    `json:"pos"`
	PlayerCards    []Card `json:"player_cards"`
}

func NewPlayerState(cards [][]Card, playerPosition int, playerId int) *PlayerState {
	state := new(PlayerState)
	state.PlayerPosition = playerPosition
	state.PlayerId = playerId
	state.PlayerCards = make([]Card, len(cards[playerPosition]))
	copy(state.PlayerCards, cards[playerPosition])
	return state
}

func (state *PlayerState) Copy() PlayerState {
	playerCards := make([]Card, len(state.PlayerCards))
	copy(playerCards, state.PlayerCards)
	return PlayerState{
		PlayerId:       state.PlayerId,
		PlayerPosition: state.PlayerPosition,
		PlayerCards:    playerCards,
	}
}
