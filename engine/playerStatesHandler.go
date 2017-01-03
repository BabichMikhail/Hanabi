package engine

type ColorInfo struct {
	Color CardColor `json:"color"`
	Count int       `json:"count"`
}

type ValueInfo struct {
	Value CardColor `json:"value"`
	Count int       `json:"count"`
}

type PlayerState struct {
	PlayerId       int         `json:"player_id"`
	PlayerPosition int         `json:"pos"`
	PlayersCards   [][]Card    `json:"players_cards"`
	ColorInfo      []ColorInfo `json:"color_info"`
	ValueInfo      []ValueInfo `json:"value_info"`
}

func NewPlayerState(cards [][]Card, playerPosition int, playerId int) PlayerState {
	this := new(PlayerState)
	this.PlayerPosition = playerPosition
	this.ColorInfo = []ColorInfo{}
	this.ValueInfo = []ValueInfo{}
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
