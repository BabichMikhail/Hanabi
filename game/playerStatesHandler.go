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
	for i := 0; i < len(state.PlayerCards); i++ {
		state.PlayerCards[i] = cards[playerPosition][i].Copy()
	}
	return state
}

func (state *PlayerState) Copy() PlayerState {
	playerCards := make([]Card, len(state.PlayerCards))
	for i := 0; i < len(playerCards); i++ {
		playerCards[i] = state.PlayerCards[i].Copy()
	}
	return PlayerState{
		PlayerId:       state.PlayerId,
		PlayerPosition: state.PlayerPosition,
		PlayerCards:    playerCards,
	}
}
