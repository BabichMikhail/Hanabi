package game

import "fmt"

type PlayerGameInfo struct {
	MyTurn          bool               `json:"my_turn"`
	CurrentPosition int                `json:"current_position"`
	PlayerCount     int                `json:"player_count"`
	Position        int                `json:"pos"`
	Step            int                `json:"step"`
	MaxStep         int                `json:"max_step"`
	Round           int                `json:"round"`
	PlayerId        int                `json:"player_id"`
	DeckSize        int                `json:"deck_size"`
	Deck            []Card             `json:"deck"`
	UsedCards       []Card             `json:"used_cards"`
	TableCards      map[CardColor]Card `json:"table_cards"`
	PlayerCards     [][]Card           `json:"player_cards"`
	PlayerCardsInfo [][]Card           `json:"player_cards_info"`
	BlueTokens      int                `json:"blue_tokens"`
	RedTokens       int                `json:"red_tokens"`
	Points          int                `json:"points"`
	VariantsCount   map[ColorValue]int `json:"-"`
	GameType        int                `json:"game_type"`
	HashKey         *string            `json:"-"`
	InfoIsSetted    bool               `json:"-"`
}

func (game *Game) GetPlayerGameInfo(playerId int, infoType int) PlayerGameInfo {
	return game.CurrentState.GetPlayerGameInfo(playerId, infoType)
}

func (state *GameState) GetPlayerGameInfoByPos(playerPosition int, infoType int) PlayerGameInfo {
	playerCardsInfo := [][]Card{}
	playerCards := [][]Card{}
	for i := 0; i < len(state.PlayerStates); i++ {
		cards1 := make([]Card, len(state.PlayerStates[i].PlayerCards))
		cards2 := make([]Card, len(state.PlayerStates[i].PlayerCards))
		for j := 0; j < len(cards1); j++ {
			cards1[j] = state.PlayerStates[i].PlayerCards[j].Copy()
			cards2[j] = state.PlayerStates[i].PlayerCards[j].Copy()
		}
		playerCardsInfo = append(playerCardsInfo, cards1)
		playerCards = append(playerCards, cards2)
	}

	for i := 0; i < len(playerCards[playerPosition]); i++ {
		card := &playerCards[playerPosition][i]
		if infoType == InfoTypeUsually {
			if !card.KnownColor {
				card.Color = NoneColor
			}
			if !card.KnownValue {
				card.Value = NoneValue
			}
		} else {
			card.KnownColor = true
			card.KnownValue = true
		}
		if card.KnownColor {
			card.ProbabilityColors = map[CardColor]float64{card.Color: 1.0}
		}
		if card.KnownValue {
			card.ProbabilityValues = map[CardValue]float64{card.Value: 1.0}
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
				card.Color = NoneColor
			} else {
				card.ProbabilityColors = map[CardColor]float64{card.Color: 1.0}
			}
			if !card.KnownValue {
				card.Value = NoneValue
			} else {
				card.ProbabilityValues = map[CardValue]float64{card.Value: 1.0}
			}
		}
	}

	deckCopy := make([]Card, len(state.Deck), len(state.Deck))
	for i := 0; i < len(deckCopy); i++ {
		deckCopy[i] = state.Deck[i].Copy()
	}
	for i := 0; i < len(deckCopy); i++ {
		card := &deckCopy[i]
		if infoType != InfoTypeFullCheat {
			card.Color = NoneColor
			card.Value = NoneValue
		} else {
			card.KnownColor = true
			card.KnownValue = true
		}
	}

	return PlayerGameInfo{
		MyTurn:          state.CurrentPosition == playerPosition,
		CurrentPosition: state.CurrentPosition,
		PlayerCount:     len(state.PlayerStates),
		Position:        playerPosition,
		Step:            state.Step,
		MaxStep:         state.MaxStep,
		Round:           state.Round,
		PlayerId:        state.PlayerStates[playerPosition].PlayerId,
		DeckSize:        len(state.Deck),
		Deck:            deckCopy,
		UsedCards:       state.UsedCards,
		TableCards:      state.TableCards,
		PlayerCards:     playerCards,
		PlayerCardsInfo: playerCardsInfo,
		BlueTokens:      state.BlueTokens,
		RedTokens:       MaxRedTokens - state.RedTokens,
		Points:          0,
		GameType:        state.GameType,
		HashKey:         nil,
		InfoIsSetted:    false,
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

const (
	InfoTypeUsually = iota
	InfoTypeCheat
	InfoTypeFullCheat
)

func (state *GameState) GetPlayerGameInfo(playerId int, infoType int) PlayerGameInfo {
	var playerPosition int
	for i := 0; i < len(state.PlayerStates); i++ {
		if state.PlayerStates[i].PlayerId == playerId {
			playerPosition = i
		}
	}
	return state.GetPlayerGameInfoByPos(playerPosition, infoType)
}

func NewPlayerInfo(game *Game, playerId int, infoType int) PlayerGameInfo {
	return game.GetPlayerGameInfo(playerId, infoType)
}

func (info PlayerGameInfo) String() string {
	result := ""
	for pos, cards := range info.PlayerCards {
		result += "[ "
		for i := 0; i < len(cards); i++ {
			result += cards[i].String()
		}
		result += " ]"
		if pos == info.CurrentPosition {
			result += " *"
		}
		result += "\n"
	}
	result += fmt.Sprintf("step %d blue %d; red %d\n", info.Step, info.BlueTokens, info.RedTokens)
	result += "[ "
	for _, color := range ColorsTable {
		result += info.TableCards[color].String()
	}
	result += " ]"
	return result
}

func (info *PlayerGameInfo) Copy() *PlayerGameInfo {
	newInfo := new(PlayerGameInfo)

	newInfo.MyTurn = info.MyTurn
	newInfo.CurrentPosition = info.CurrentPosition
	newInfo.PlayerCount = info.PlayerCount
	newInfo.Position = info.Position
	newInfo.Step = info.Step
	newInfo.MaxStep = info.MaxStep
	newInfo.Round = info.Round
	newInfo.PlayerId = info.PlayerId
	newInfo.GameType = info.GameType
	newInfo.VariantsCount = map[ColorValue]int{}
	for k, v := range info.VariantsCount {
		newInfo.VariantsCount[k] = v
	}
	newInfo.Deck = make([]Card, len(info.Deck), len(info.Deck))
	for i := 0; i < len(info.Deck); i++ {
		newInfo.Deck[i] = info.Deck[i].Copy()
	}
	newInfo.DeckSize = info.DeckSize
	newInfo.Points = info.Points

	newInfo.UsedCards = make([]Card, len(info.UsedCards), cap(info.UsedCards))
	for i := 0; i < len(newInfo.UsedCards); i++ {
		newInfo.UsedCards[i] = info.UsedCards[i].Copy()
	}

	newInfo.TableCards = map[CardColor]Card{}
	for color, card := range info.TableCards {
		newInfo.TableCards[color] = card.Copy()
	}

	newInfo.PlayerCards = make([][]Card, len(info.PlayerCards), len(info.PlayerCards))
	for i := 0; i < len(newInfo.PlayerCards); i++ {
		newInfo.PlayerCards[i] = make([]Card, len(info.PlayerCards[i]), len(info.PlayerCards[i]))
		for j := 0; j < len(newInfo.PlayerCards[i]); j++ {
			newInfo.PlayerCards[i][j] = info.PlayerCards[i][j].Copy()
		}
	}

	newInfo.PlayerCardsInfo = make([][]Card, len(info.PlayerCardsInfo), len(info.PlayerCardsInfo))
	for i := 0; i < len(newInfo.PlayerCardsInfo); i++ {
		newInfo.PlayerCardsInfo[i] = make([]Card, len(info.PlayerCardsInfo[i]), len(info.PlayerCardsInfo[i]))
		for j := 0; j < len(newInfo.PlayerCardsInfo[i]); j++ {
			newInfo.PlayerCardsInfo[i][j] = info.PlayerCardsInfo[i][j].Copy()
		}
	}

	newInfo.BlueTokens = info.BlueTokens
	newInfo.RedTokens = info.RedTokens
	newInfo.HashKey = info.HashKey
	newInfo.InfoIsSetted = info.InfoIsSetted
	return newInfo
}

func (info *PlayerGameInfo) IsGameOver() bool {
	if info.RedTokens == MaxRedTokens || info.MaxStep != 0 && info.Step >= info.MaxStep {
		return true
	}

	if info.GameType == Type_InfinityGame {
		for i := 0; i < len(info.PlayerCards); i++ {
			if len(info.PlayerCards[i]) == 0 && info.BlueTokens == 0 {
				return true
			}
		}

		gameOver := true
		for i := 0; i < len(info.PlayerCards); i++ {
			if len(info.PlayerCards[i]) > 0 {
				gameOver = false
			}
		}

		if gameOver {
			return gameOver
		}
	}

	for _, card := range info.TableCards {
		if card.Value != Five {
			return false
		}
	}

	return true
}

func (info *PlayerGameInfo) GetDefaultDeck() map[ColorValue]int {
	cards := map[ColorValue]int{}
	for _, value := range ValuesTable {
		for _, color := range ColorsTable {
			if value == 1 {
				cards[ColorValue{Color: color, Value: value}] = 3
			}
			if value >= 2 && value <= 4 {
				cards[ColorValue{Color: color, Value: value}] = 2
			}
			if value == 5 {
				cards[ColorValue{Color: color, Value: value}] = 1
			}
		}
	}
	return cards
}

func (info *PlayerGameInfo) GetUnplayedCards() map[ColorValue]int {
	cards := info.GetDefaultDeck()
	for _, card := range info.UsedCards {
		cards[ColorValue{Color: card.Color, Value: card.Value}]--
	}
	for color, card := range info.TableCards {
		for value := CardValue(1); value < card.Value; value++ {
			cards[ColorValue{Color: color, Value: value}]--
		}
	}
	return cards
}

func (info *PlayerGameInfo) IsCardPlayable(card *Card) bool {
	card.CheckVisible()
	return info.TableCards[card.Color].Value+1 == card.Value
}

func (info *PlayerGameInfo) GetPlayableCardPositions() []int {
	result := []int{}
	pos := info.CurrentPosition
	for i := 0; i < len(info.PlayerCards[pos]); i++ {
		if info.IsCardPlayable(&info.PlayerCards[pos][i]) {
			result = append(result, i)
		}
	}
	return result
}

func (info *PlayerGameInfo) IncreasePosition() {
	info.Step++
	info.CurrentPosition++
	if info.CurrentPosition/len(info.PlayerCards) == 1 {
		info.CurrentPosition = 0
		info.Round++
	}
}
