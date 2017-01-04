package game

type PlayerState struct {
	PlayerId       int      `json:"player_id"`
	PlayerPosition int      `json:"pos"`
	PlayersCards   [][]Card `json:"players_cards"`
}

func NewPlayerState(cards [][]Card, playerPosition int, playerId int) PlayerState {
	this := new(PlayerState)
	this.PlayerPosition = playerPosition
	this.PlayerId = playerId

	copyCards := [][]Card{}
	for i := 0; i < len(cards); i++ {
		copyPlayerCards := make([]Card, len(cards[i]))
		copy(copyPlayerCards, cards[i])
		for j := 0; j < len(cards[playerPosition]); j++ {
			copyPlayerCards[j].SetKnown(i != playerPosition)
		}
		copyCards = append(copyCards, copyPlayerCards)
	}
	this.PlayersCards = copyCards
	return *this
}
