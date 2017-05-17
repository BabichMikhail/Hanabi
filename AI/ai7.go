package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

type AI7 struct {
	BaseAI
	CardCodes map[game.HashValue]game.ColorValue
}

func NewAI7(baseAI *BaseAI) *AI7 {
	ai := new(AI7)
	ai.BaseAI = *baseAI
	ai.CardCodes = map[game.HashValue]game.ColorValue{}
	for _, color := range game.ColorsTable {
		for _, value := range game.ValuesTable {
			ai.CardCodes[game.HashColorValue(color, value)] = game.ColorValue{Color: color, Value: value}
		}
	}
	return ai
}

func (ai *AI7) GetClueActions(info *game.PlayerGameInfo, myPos, step int) []*game.Action {
	actions := []*game.Action{}
	for i := 1; i < info.PlayerCount; i++ {
		cluedPos := (myPos + i) % info.PlayerCount
		for _, value := range game.ValuesTable {
			actions = append(actions, game.NewAction(game.TypeActionInformationValue, cluedPos, int(value)))
		}
	}
	for i := 1; i < info.PlayerCount; i++ {
		cluedPos := (myPos + i) % info.PlayerCount
		for _, color := range game.ColorsTable {
			actions = append(actions, game.NewAction(game.TypeActionInformationColor, cluedPos, int(color)))
		}
	}

	for j := len(actions) - 1; j >= 0; j-- {
		for i := Max(0, step-info.PlayerCount); i < step; i++ {
			if actions[j].Equal(&ai.History[i]) {
				actions = append(actions[:j], actions[j+1:]...)
				break
			}
		}
	}

	return actions
}

func (ai *AI7) ClueOnStep(info *game.PlayerGameInfo, chooseCardPos func(int, int) int) *game.Action {
	myPos := info.CurrentPosition
	hashsum := game.HashValue(0)
	for i := 1; i < info.PlayerCount; i++ {
		pos := (myPos + i) % info.PlayerCount
		cardPos := chooseCardPos(pos, i)
		card := &info.PlayerCards[pos][cardPos]
		card.CheckVisible()
		hashsum += game.HashColorValue(card.Color, card.Value)
	}

	actions := ai.GetClueActions(info, myPos, info.Step)
	return actions[int(hashsum)%(len(actions))]
}

func (ai *AI7) Clue(info *game.PlayerGameInfo) *game.Action {
	return ai.ClueOnStep(info, func(pos, i int) int {
		cards := info.PlayerCards[pos]
		return len(cards) - 1
	})
}

func (ai *AI7) FirstRoundClue(info *game.PlayerGameInfo) *game.Action {
	return ai.ClueOnStep(info, func(pos, i int) int {
		return i - 1
	})
}

func (ai *AI7) DecodeClueOnStep(info *game.PlayerGameInfo, firstPos int, step int, chooseCardPos func(int, int) int) {
	/* (x + c) % y = v */
	type EqualParams struct {
		c, y, v int
		cardPos int
	}

	myPos := info.Position
	var params EqualParams
	cluerPos := firstPos
	if cluerPos == myPos {
		return
	}
	for j := 1; j < info.PlayerCount; j++ {
		pos := (cluerPos + j) % info.PlayerCount
		cardPos := chooseCardPos(pos, j)
		if pos == myPos {
			params.cardPos = cardPos
			continue
		}
		card := &info.PlayerCards[pos][cardPos]
		card.CheckVisible()
		params.c += int(game.HashColorValue(card.Color, card.Value))
	}

	actions := ai.GetClueActions(info, cluerPos, step)
	params.y = len(actions)

	for j, action := range actions {
		if ai.History[step].Equal(action) {
			params.v = j
		}
	}

	c, v, y := params.c, params.v, params.y
	count := 0
	for x := 0; x < 25; x++ {
		if (x+c)%y == v {
			if count == 1 {
				panic("Magic")
			}
			count++
			color, value := game.ColorValueByHashColorValue(game.HashValue(x))
			card := &info.PlayerCards[myPos][params.cardPos]
			card.SetColor(color)
			card.SetValue(value)
		}
	}
	if count == 0 {
		panic("Bad solve")
	}
}

func (ai *AI7) DecodeFirstRoundClue(currentInfo *game.PlayerGameInfo) {
	var info *game.PlayerGameInfo
	myPos := currentInfo.Position
	if currentInfo.Step == 5 {
		info = currentInfo
	} else {
		i := ai.Informator.GetPlayerState(5)
		info = &i
		info.PlayerCards[myPos] = currentInfo.PlayerCards[myPos]
	}

	firstPos := (myPos - currentInfo.Step + info.PlayerCount) % info.PlayerCount
	for i := 0; i < info.PlayerCount; i++ {
		cluerPos := (firstPos + i) % info.PlayerCount
		ai.DecodeClueOnStep(info, cluerPos, cluerPos, func(pos, j int) int {
			return j - 1
		})
	}
}

func (ai *AI7) DecodeClue(info *game.PlayerGameInfo) {
	ai.DecodeClueOnStep(info, info.CurrentPosition, info.Step, func(pos, j int) int {
		cards := info.PlayerCards[pos]
		return len(cards) - 1
	})
}

func (ai *AI7) PlayLowest(info *game.PlayerGameInfo) *game.Action {
	lowestValue := 6
	myPos := info.CurrentPosition
	cards := info.PlayerCards[myPos]
	cardPos := -1
	for i, card := range cards {
		card.CheckVisible()
		if ai.isCardPlayable(card) && int(card.Value) < lowestValue {
			lowestValue = int(card.Value)
			cardPos = i
		}
	}
	if cardPos != -1 {
		return game.NewAction(game.TypeActionPlaying, myPos, cardPos)
	}
	return nil
}

func (ai *AI7) TryPlayFive(info *game.PlayerGameInfo) *game.Action {
	myPos := info.CurrentPosition
	cards := info.PlayerCards[myPos]
	for cardPos, card := range cards {
		if card.Value == 5 && ai.isCardPlayable(card) {
			return game.NewAction(game.TypeActionPlaying, myPos, cardPos)
		}
	}
	return nil
}

func (ai *AI7) GetLastClueActionAndNotCluedAction(info *game.PlayerGameInfo) int {
	lastClue := -1
	for i := len(ai.History) - 1; i >= 0; i-- {
		if ai.History[i].IsInfoAction() {
			lastClue = i
			break
		}
	}

	if lastClue == -1 {
		panic("Magic")
	}

	firstNotClue := info.Step
	for i := info.Step - 1; i >= lastClue; i-- {
		if action := ai.History[i]; !action.IsInfoAction() {
			if len(info.PlayerCards[action.PlayerPosition]) < 4 {
				/* If last round that I don't need in clue */
				break
			}
			firstNotClue = i
		}
	}
	if info.Step-firstNotClue > 4 {
		panic("Missing clue")
	}
	return firstNotClue
}

func (ai *AI7) ClueIsUseful(info *game.PlayerGameInfo) bool {
	firstNotClue := ai.GetLastClueActionAndNotCluedAction(info)
	return info.Step-firstNotClue > 0
}

func (ai *AI7) NeedClue(info *game.PlayerGameInfo) bool {
	firstNotClue := ai.GetLastClueActionAndNotCluedAction(info)
	return info.Step-firstNotClue == 4
}

func (ai *AI7) ClueTime(info *game.PlayerGameInfo) int {
	firstNotClue := ai.GetLastClueActionAndNotCluedAction(info)
	return 4 - (info.Step - firstNotClue)
}

func (ai *AI7) TryDiscard(info *game.PlayerGameInfo) *game.Action {
	myPos := info.CurrentPosition
	for i, card := range info.PlayerCards[myPos] {
		if info.TableCards[card.Color].Value >= card.Value {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
	}

	for i, card := range info.PlayerCards[myPos] {
		playedCards := info.GetUnplayedCards()
		if playedCards[game.ColorValue{Color: card.Color, Value: card.Value}] > 1 && !ai.isCardPlayable(card) {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
	}
	return nil
}

func (ai *AI7) Discard(info *game.PlayerGameInfo) *game.Action {
	myPos := info.CurrentPosition
	highestValue := 0
	cardPos := -1
	for i, card := range info.PlayerCards[myPos] {
		if int(card.Value) > highestValue {
			highestValue = int(card.Value)
			cardPos = i
		}
	}
	return game.NewAction(game.TypeActionDiscard, myPos, cardPos)
}

func (ai *AI7) LoadInfo(info *game.PlayerGameInfo) {
	myPos := info.Position
	myOldInfo := ai.Informator.GetCache().(*game.PlayerGameInfo)
	if myOldInfo == nil {
		panic("Bad Cache")
	}

	if myOldInfo.Step != info.Step-info.PlayerCount {
		panic("Bad Cache 2")
	}
	myLastAction := ai.History[myOldInfo.Step]
	myOldCards := myOldInfo.PlayerCards[myPos]
	myCards := info.PlayerCards[myPos]
	if len(myCards) == 3 {
		panic("FWEJIFJEWIF")
	}
	if myLastAction.IsInfoAction() {
		for i := 0; i < len(myCards); i++ {
			card := &myCards[i]
			card.SetColor(myOldCards[i].Color)
			card.SetValue(myOldCards[i].Value)
		}
	} else {
		for i := 0; i < len(myCards)-1; i++ {
			card := &myCards[i]
			idx := i
			if i >= myLastAction.Value {
				idx++
			}
			oldCard := &myOldCards[idx]
			card.SetColor(oldCard.Color)
			card.SetValue(oldCard.Value)
		}

		lastCard := len(myCards) - 1
		card := &myCards[lastCard]
		for i := myOldInfo.Step + 1; i < info.Step; i++ {
			action := &ai.History[i]
			if action.IsInfoAction() {
				myCluedInfo := ai.Informator.GetPlayerState(i)
				ai.DecodeClue(&myCluedInfo)
				newCard := myCluedInfo.PlayerCards[myPos][lastCard]
				newCard.CheckVisible()
				card.SetColor(newCard.Color)
				card.SetValue(newCard.Value)
			}
		}

		if myOldInfo.DeckSize == 1 {
			deckCard := &myOldInfo.Deck[0]
			card.SetColor(deckCard.Color)
			card.SetValue(deckCard.Value)
			card.CheckVisible()
		}
	}

	if info.DeckSize == 1 {
		info.SetVariantsCount(false, false)
		for colorValue, _ := range info.VariantsCount {
			color, value := colorValue.Color, colorValue.Value
			info.Deck[0].SetColor(color)
			info.Deck[0].SetValue(value)
		}
	} else if info.DeckSize == 0 {
		info.SetVariantsCount(false, false)
		countSum := 0
		var colorValue game.ColorValue
		for cv, count := range info.VariantsCount {
			countSum += count
			colorValue = cv
		}
		if countSum > 1 {
			panic("Bad count")
		}
		card := &myCards[len(myCards)-1]
		if countSum == 0 {
			card.CheckVisible()
		}
		card.SetColor(colorValue.Color)
		card.SetValue(colorValue.Value)
		card.CheckVisible()
	}

	for _, card := range myCards {
		card.CheckVisible()
	}
	ai.Informator.SetCache(info.Copy())
}

func (ai *AI7) GetAction() *game.Action {
	info := &ai.PlayerInfo
	myPos := info.CurrentPosition
	if info.Step < info.PlayerCount {
		return ai.FirstRoundClue(info)
	}

	if info.Step >= info.PlayerCount && info.Step < 2*info.PlayerCount {
		ai.DecodeFirstRoundClue(info)
		ai.Informator.SetCache(info.Copy())
		if ai.NeedClue(info) {
			return ai.Clue(info)
		}
		for i, card := range info.PlayerCards[myPos] {
			card.CheckVisible()
			if ai.isCardPlayable(card) {
				return game.NewAction(game.TypeActionPlaying, myPos, i)
			}
		}

		discard := ai.TryDiscard(info)
		if discard != nil {
			return discard
		}

		if info.BlueTokens == 0 {
			panic("Panic I Have No BLUE TOKENS")
		}

		return ai.Clue(info)
	}

	ai.LoadInfo(info)
	if ai.NeedClue(info) {
		return ai.Clue(info)
	}

	clueTime := ai.ClueTime(info)
	if info.BlueTokens == 1 || info.BlueTokens == 0 && clueTime < 2 {
		playFive := ai.TryPlayFive(info)
		if playFive != nil {
			return playFive
		}
		discard := ai.TryDiscard(info)
		if discard != nil {
			return discard
		}
	}

	if info.BlueTokens == 0 && clueTime == 0 {
		return ai.Discard(info)
	}

	play := ai.PlayLowest(info)
	if play != nil {
		return play
	}

	if info.BlueTokens < 3 {
		discard := ai.TryDiscard(info)
		if discard != nil {
			return discard
		}
	}

	if info.BlueTokens > 0 {
		return ai.Clue(info)
	} else {
		return ai.Discard(info)
	}
}
