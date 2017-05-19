package ai

import (
	"fmt"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI8 struct {
	BaseAI
	CardCodes map[game.HashValue]game.ColorValue
	Knowledge []int
}

func NewAI8(baseAI *BaseAI) *AI8 {
	ai := new(AI8)
	ai.BaseAI = *baseAI
	ai.CardCodes = map[game.HashValue]game.ColorValue{}
	for _, color := range game.ColorsTable {
		for _, value := range game.ValuesTable {
			ai.CardCodes[game.HashColorValue(color, value)] = game.ColorValue{Color: color, Value: value}
		}
	}
	return ai
}

func (ai *AI8) GetClueActions(info *game.PlayerGameInfo, myPos, step int) []*game.Action {
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

	lowestValue := 5
	completeColors := map[game.CardColor]struct{}{}
	for color, card := range info.TableCards {
		if int(card.Value) < lowestValue {
			lowestValue = int(card.Value)
		}
		if card.Value == 5 {
			completeColors[color] = struct{}{}
		}
	}

	for j := len(actions) - 1; j >= 0 && len(actions) > 25; j-- {
		action := actions[j]
		if action.ActionType == game.TypeActionInformationValue && action.Value <= lowestValue {
			actions = append(actions[:j], actions[j+1:]...)
		}
	}

	for j := len(actions) - 1; j >= 0 && len(actions) > 25; j-- {
		action := actions[j]
		if _, ok := completeColors[game.CardColor(action.Value)]; action.ActionType == game.TypeActionInformationColor && ok {
			actions = append(actions[:j], actions[j+1:]...)
		}
	}

	return actions
}

func (ai *AI8) ClueOnStep(info *game.PlayerGameInfo, chooseCardPos func(int, int) int) *game.Action {
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

func (ai *AI8) Clue(info *game.PlayerGameInfo) *game.Action {
	myPos := info.Position
	return ai.ClueOnStep(info, func(pos, i int) int {
		cardPos := ai.Knowledge[pos]
		if pos == myPos {
			panic("Bad clue")
		}
		ai.Knowledge[pos] = Min(ai.Knowledge[pos]+1, len(info.PlayerCards[pos]))
		return Min(cardPos, len(info.PlayerCards[pos])-1)
	})
}

func (ai *AI8) DecodeClueOnStep(info *game.PlayerGameInfo, firstPos int, step int, chooseCardPos func(int, int) int) {
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

func (ai *AI8) DecodeClue(info *game.PlayerGameInfo) {
	myPos := info.Position
	ai.DecodeClueOnStep(info, info.CurrentPosition, info.Step, func(pos, j int) int {
		if myPos == info.CurrentPosition {
			panic("Bad cluer position")
		}
		cardPos := ai.Knowledge[pos]
		if cardPos < len(info.PlayerCards[pos]) {
			ai.Knowledge[pos]++
		}
		return Min(cardPos, len(info.PlayerCards[pos])-1)
	})
}

func (ai *AI8) DecodeClues(info *game.PlayerGameInfo) {
	myPos := info.Position
	myCards := info.PlayerCards[myPos]

	lastStep := Max(-1, info.Step-info.PlayerCount)
	for i := Max(-1, info.Step-info.PlayerCount) + 1; i < info.Step && ai.Knowledge[myPos] < len(myCards); i++ {
		lastStep = i
		cardPos := ai.Knowledge[myPos]
		if cardPos == len(myCards) {
			cardPos--
		}
		card := &myCards[cardPos]
		action := &ai.History[i]
		if action.IsInfoAction() {
			myCluedInfo := ai.Informator.GetPlayerState(i)
			ai.DecodeClue(&myCluedInfo)
			newCard := myCluedInfo.PlayerCards[myPos][cardPos]
			newCard.CheckVisible()
			card.SetColor(newCard.Color)
			card.SetValue(newCard.Value)
		} else {
			pos := action.PlayerPosition
			if pos == info.Position {
				panic("Bad action")
			}
			if action.Value < ai.Knowledge[pos] {
				ai.Knowledge[pos]--
			}
			if ai.Knowledge[pos] < 0 {
				panic("Bad knowledge")
			}
		}
	}

	for i := lastStep + 1; i < info.Step; i++ {
		action := ai.History[i]
		if action.IsInfoAction() {
			for j := 0; j < len(ai.Knowledge); j++ {
				if j == myPos {
					continue
				}
				cluerPos := (info.CurrentPosition - (info.Step - i) + 100*info.PlayerCount) % info.PlayerCount
				if j == cluerPos {
					continue
				}
				ai.Knowledge[j] = Min(ai.Knowledge[j]+1, len(info.PlayerCards[j]))
			}
		} else {
			pos := action.PlayerPosition
			if pos == info.Position {
				panic("Bad action")
			}
			if action.Value < ai.Knowledge[pos] {
				ai.Knowledge[pos]--
			}
			if ai.Knowledge[pos] < 0 {
				panic("Bad knowledge")
			}
		}
	}
}

func (ai *AI8) TryDiscard(info *game.PlayerGameInfo, pos ...int) *game.Action {
	myPos := info.CurrentPosition
	if len(pos) == 1 {
		myPos = pos[0]
	}
	myCards := info.PlayerCards[myPos][:Min(len(info.PlayerCards[myPos]), ai.Knowledge[myPos])]
	for i, card := range myCards { /* useless */
		if info.TableCards[card.Color].Value >= card.Value {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
	}

	cardCounts := map[game.ColorValue]int{}
	for i, card := range myCards { /* duplicate */
		colorValue := game.ColorValue{Color: card.Color, Value: card.Value}
		count := cardCounts[colorValue]
		if count == 1 {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
		cardCounts[colorValue]++
	}

	for i, card := range myCards { /* not last card */
		unplayedCards := info.GetUnplayedCards()
		if unplayedCards[game.ColorValue{Color: card.Color, Value: card.Value}] > 1 {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
	}
	return nil
}

func (ai *AI8) Discard(info *game.PlayerGameInfo) *game.Action {
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

func (ai *AI8) InitializeCache() {
	ai.Knowledge = make([]int, ai.PlayerInfo.PlayerCount)
	myFirstInfo := ai.Informator.GetPlayerState(0)
	ai.Informator.SetCache([]interface{}{ai.Knowledge, myFirstInfo.Copy()})
}

func (ai *AI8) LoadInformation(info *game.PlayerGameInfo) {
	myPos := info.Position
	myCachedData := ai.Informator.GetCache().([]interface{})
	ai.Knowledge = myCachedData[0].([]int)
	myOldInfo := myCachedData[1].(*game.PlayerGameInfo)
	if myOldInfo == nil || myOldInfo.Step != Max(0, info.Step-info.PlayerCount) {
		panic("Bad Cache")
	}

	if info.Step >= info.PlayerCount {
		myCards := info.PlayerCards[myPos][:ai.Knowledge[myPos]]
		var myLastAction *game.Action
		if myOldInfo.CurrentPosition == myOldInfo.Position {
			myLastAction = &ai.History[myOldInfo.Step]
		}

		myOldCards := myOldInfo.PlayerCards[myPos]
		if myLastAction == nil || myLastAction.IsInfoAction() {
			for i := 0; i < len(myCards); i++ {
				card := &myCards[i]
				card.SetColor(myOldCards[i].Color)
				card.SetValue(myOldCards[i].Value)
			}
		} else {
			if myPos != myLastAction.PlayerPosition {
				panic("It's not my action")
			}
			for i := 0; i < ai.Knowledge[myPos]; i++ {
				card := &myCards[i]
				idx := i
				if i >= myLastAction.Value {
					idx++
				}
				oldCard := &myOldCards[idx]
				oldCard.CheckVisible()
				card.SetColor(oldCard.Color)
				card.SetValue(oldCard.Value)
			}
		}
	}
	ai.DecodeClues(info)

	/* @todo I don't know cards[i] but i know cards[i + 1] */
	/*myCards := info.PlayerCards[myPos]
	for i := info.PlayerCards[myPos]; i < len(myCards); i++ {
		card := &myCards[i]
		if card.KnownColor && card.KnownValue {
			ai.Knowledge[myPos]++
		}
	}

	if info.DeckSize == 1 && ai.Knowledge[myPos] == len(myCards) {
		info.SetVariantsCount(false, false)
		for colorValue, _ := range info.VariantsCount {
			color, value := colorValue.Color, colorValue.Value
			info.Deck[0].SetColor(color)
			info.Deck[0].SetValue(value)
		}
	} else if info.DeckSize == 0 && ai.Knowledge[myPos] == len(myCards)-1 {
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
		ai.Knowledge[myPos]++
	}*/

	ai.Informator.SetCache([]interface{}{ai.Knowledge, info.Copy()})
}

func (ai *AI8) GetProgress(info *game.PlayerGameInfo) map[game.CardColor]game.CardValue {
	progress := map[game.CardColor]game.CardValue{}
	for color, card := range info.TableCards {
		progress[color] = card.Value
	}
	return progress
}

func (ai *AI8) GetPoints(progress map[game.CardColor]game.CardValue) int {
	points := 0
	for _, value := range progress {
		points += int(value)
	}
	return points
}

func (ai *AI8) GetPlayableUnknownCards(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) (int, int) {
	cards := info.PlayerCards[pos][ai.Knowledge[pos]:]
	unplayedCards := info.GetUnplayedCards()
	topIdx := -1
	topRel := 1.0
	untopIdx := -1
	untopRel := 0.0
	for i, card := range cards {
		sumCount := 0
		playedCount := 0
		for unplayedCard, count := range unplayedCards {
			color, value := unplayedCard.Color, unplayedCard.Value
			if _, ok := card.ProbabilityColors[color]; !ok {
				continue
			}
			if _, ok := card.ProbabilityValues[value]; !ok {
				continue
			}
			if progress[color]+1 == value {
				playedCount += count
			}
			sumCount += count
		}

		newRel := float64(playedCount) / float64(sumCount)
		if untopIdx == -1 || newRel < untopRel {
			untopIdx = i
			untopRel = newRel
		}
		if topIdx == -1 || newRel > topRel {
			topIdx = i
			topRel = newRel
		}
	}
	return topIdx, untopIdx
}

func (ai *AI8) HardDiscard(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) *game.Action {
	_, topIdx := ai.GetPlayableUnknownCards(info, progress, pos)
	return game.NewAction(game.TypeActionPlaying, pos, topIdx)
}

func (ai *AI8) RiskyPlay(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) *game.Action {
	topIdx, _ := ai.GetPlayableUnknownCards(info, progress, pos)
	return game.NewAction(game.TypeActionPlaying, pos, topIdx)
}

func (ai *AI8) DiscardHighest(info *game.PlayerGameInfo, pos int) *game.Action {
	highestIdx := 0
	cards := info.PlayerCards[pos]
	for i := 0; i < len(cards); i++ {
		if cards[i].Value > cards[highestIdx].Value {
			highestIdx = i
		}
	}
	return game.NewAction(game.TypeActionDiscard, pos, highestIdx)
}

func (ai *AI8) GetHardAction(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos, deep int) *game.Action {
	var action *game.Action
	myCards := info.PlayerCards[pos]
	if ai.Knowledge[pos] < len(myCards) {
		if info.RedTokens < 2 {
			if deep == ai.GetMaxDeep(info) {
				fmt.Println("Risky  play")
			}
			action = ai.RiskyPlay(info, progress, pos)
		} else {
			if deep == ai.GetMaxDeep(info) {
				fmt.Println("Hard discard")
			}
			action = ai.HardDiscard(info, progress, pos)
		}
	} else {
		if deep == ai.GetMaxDeep(info) {
			fmt.Println("Critical discard")
		}
		action = ai.DiscardHighest(info, pos)
	}
	return action
}

func (ai *AI8) GetBestAction(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos, deep int) (*game.Action, int) {
	if deep == 1 {
		myCards := info.PlayerCards[pos]
		for i := 0; i < Min(ai.Knowledge[pos], len(myCards)); i++ {
			card := &info.PlayerCards[pos][i]
			if card.IsCardPlayable(progress) {
				progress[card.Color]++
				points := ai.GetPoints(progress)
				progress[card.Color]--
				return game.NewAction(game.TypeActionPlaying, pos, i), points
			}
		}
		var action *game.Action
		if info.BlueTokens == 0 {
			action = ai.TryDiscard(info)
			if action == nil {
				action = ai.GetHardAction(info, progress, pos, deep)
			}
		} else {
			action = game.NewAction(game.TypeActionInformationColor, (pos+1)%info.PlayerCount, 1)
		}
		return action, ai.GetPoints(progress)
	}

	nextPos := (pos + 1) % info.PlayerCount
	topPoints := -1
	var topAction *game.Action

	for i := 0; i < Min(ai.Knowledge[pos], len(info.PlayerCards[pos])); i++ {
		card := &info.PlayerCards[pos][i]
		if card.IsCardPlayable(progress) {
			progress[card.Color]++
			dblueTokens := 0
			if card.Value == 5 {
				dblueTokens++
				info.BlueTokens++
			}
			_, points := ai.GetBestAction(info, progress, nextPos, deep-1)
			progress[card.Color]--
			info.BlueTokens -= dblueTokens
			if points > topPoints {
				topPoints = points
				topAction = game.NewAction(game.TypeActionPlaying, pos, i)
			}
		}
	}

	if info.BlueTokens > 0 {
		info.BlueTokens--
		for i := 0; i < len(ai.Knowledge); i++ {
			if i != pos {
				ai.Knowledge[i]++
			}
		}

		_, points := ai.GetBestAction(info, progress, nextPos, deep-1)
		if points > topPoints {
			topPoints = points
			/* Some Clue */
			topAction = game.NewAction(game.TypeActionInformationColor, nextPos, 1)
		}
		for i := 0; i < len(ai.Knowledge); i++ {
			if i != pos {
				ai.Knowledge[i]--
			}
		}
		info.BlueTokens++
	}

	if info.BlueTokens < game.MaxBlueTokens {
		info.BlueTokens++
		_, points := ai.GetBestAction(info, progress, nextPos, deep-1)
		if points > topPoints {
			discard := ai.TryDiscard(info)
			if discard != nil {
				topPoints = points
				topAction = discard
			}
		}
		info.BlueTokens--
	}

	if topAction == nil { // hard action
		topAction = ai.GetHardAction(info, progress, pos, deep)
		if topAction.ActionType == game.TypeActionDiscard {
			info.BlueTokens++
		} else {
			info.RedTokens++
		}

		_, points := ai.GetBestAction(info, progress, pos, deep-1)
		if topAction.ActionType == game.TypeActionDiscard {
			info.BlueTokens--
		} else {
			info.RedTokens--
		}

		if points != 25 {
			points = -1
		}
		topPoints = points

	}

	return topAction, topPoints
}

func (ai *AI8) GetMaxDeep(info *game.PlayerGameInfo) int {
	deep := info.PlayerCount
	if info.DeckSize == 0 {
		deep = info.MaxStep - info.Step
	}
	if deep == 0 {
		panic("Bad deep")
	}
	return deep
}

func (ai *AI8) FindBestAction(info *game.PlayerGameInfo) *game.Action {
	myPos := info.CurrentPosition
	progress := ai.GetProgress(info)
	action, _ := ai.GetBestAction(info, progress, myPos, ai.GetMaxDeep(info))
	if action.IsInfoAction() {
		action = ai.Clue(info)
	} else if action.Value < ai.Knowledge[myPos] {
		ai.Knowledge[myPos]--
	}
	return action
}

func (ai *AI8) GetAction() *game.Action {
	info := &ai.PlayerInfo

	if info.Step < info.PlayerCount {
		ai.InitializeCache()
	}
	ai.LoadInformation(info)
	return ai.FindBestAction(info)
}
