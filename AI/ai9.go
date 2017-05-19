package ai

import (
	"fmt"
	"math"
	"sort"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI9 struct {
	BaseAI
	CardCodes map[game.HashValue]game.ColorValue
	Knowledge [][]CardKnowledge
}

func NewAI9(baseAI *BaseAI) *AI9 {
	ai := new(AI9)
	ai.BaseAI = *baseAI
	ai.CardCodes = map[game.HashValue]game.ColorValue{}
	for _, color := range game.ColorsTable {
		for _, value := range game.ValuesTable {
			ai.CardCodes[game.HashColorValue(color, value)] = game.ColorValue{Color: color, Value: value}
		}
	}
	return ai
}

func (ai *AI9) FilterClueActions() {

}

func (ai *AI9) GetClueActions(info *game.PlayerGameInfo, myPos, step int) []*game.Action {
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

	/*lowestValue := 5
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
	}*/
	fmt.Println(actions)
	return actions
}

type Variants []*Variant
type Variant struct {
	Count int
	Pos   int
}

func (v Variants) Len() int {
	return len(v)
}

func (v Variants) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Variants) Less(i, j int) bool {
	return v[i].Count > v[j].Count || v[i].Count == v[j].Count && v[i].Pos < v[j].Pos
}

// Useless
// Playable
// Playable + 1
// Playable + 2...

type CardKnowledge map[game.ColorValue]struct{}

func (ai *AI9) NewCardKnowledge() CardKnowledge {
	knowledge := CardKnowledge{}
	for _, color := range game.ColorsTable {
		for _, value := range game.ValuesTable {
			knowledge[game.ColorValue{Color: color, Value: value}] = struct{}{}
		}
	}
	return knowledge
}

func (ai *AI9) GetVariants(info *game.PlayerGameInfo, variants CardKnowledge, f func(game.CardValue, game.CardValue) bool) CardKnowledge {
	result := CardKnowledge{}
	for cv, _ := range variants {
		if f(info.TableCards[cv.Color].Value, cv.Value) {
			result[cv] = struct{}{}
		}
	}
	return result
}

func (ai *AI9) GetUselessVariants(info *game.PlayerGameInfo, knowledge CardKnowledge) CardKnowledge {
	return ai.GetVariants(info, knowledge, func(tableValue game.CardValue, variantValue game.CardValue) bool {
		return tableValue >= variantValue
	})
}

func (ai *AI9) GetPlayableVariants(info *game.PlayerGameInfo, knowledge CardKnowledge) CardKnowledge {
	return ai.GetVariants(info, knowledge, func(tableValue game.CardValue, variantValue game.CardValue) bool {
		return tableValue+1 == variantValue
	})
}

func (ai *AI9) GetUnplayableVariants(info *game.PlayerGameInfo, knowledge CardKnowledge) CardKnowledge {
	return ai.GetVariants(info, knowledge, func(tableValue game.CardValue, variantValue game.CardValue) bool {
		return tableValue+1 != variantValue
	})
}

func (ai *AI9) GetSoonPlayableVariants(info *game.PlayerGameInfo, knowledge CardKnowledge) CardKnowledge {
	return ai.GetVariants(info, knowledge, func(tableValue game.CardValue, variantValue game.CardValue) bool {
		return tableValue+2 == variantValue
	})
}

func (ai *AI9) GetAnyOnePlayableVariants(info *game.PlayerGameInfo, knowledge CardKnowledge) CardKnowledge {
	return ai.GetVariants(info, knowledge, func(tableValue game.CardValue, variantValue game.CardValue) bool {
		return tableValue+2 < variantValue
	})
}

func (kn CardKnowledge) RemoveWithColor(color game.CardColor) {
	fmt.Println("REMOVE0")
	for cv, _ := range kn {
		if cv.Color == color {
			delete(kn, cv)
		}
	}
}

func (kn CardKnowledge) RemoveWithoutColor(color game.CardColor) {
	fmt.Println("REMOVE1")
	for cv, _ := range kn {
		if cv.Color != color {
			delete(kn, cv)
		}
	}
}

func (kn CardKnowledge) RemoveWithValue(value game.CardValue) {
	fmt.Println("REMOVE2")
	for cv, _ := range kn {
		if cv.Value == value {
			delete(kn, cv)
		}
	}
}

func (kn CardKnowledge) RemoveWithoutValue(value game.CardValue) {
	fmt.Println("REMOVE3")
	for cv, _ := range kn {
		if cv.Value != value {
			delete(kn, cv)
		}
	}
}

func (ai *AI9) MergeVariants(variants1, variants2 CardKnowledge) CardKnowledge {
	result := CardKnowledge{}
	for k, v := range variants1 {
		result[k] = v
	}
	for k, v := range variants2 {
		result[k] = v
	}
	return result
}

func Pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func (ai *AI9) Encode2PlayerClue(info *game.PlayerGameInfo) *game.Action {
	fmt.Println("Clue")
	points := info.GetPoints()
	myPos := info.Position
	cluerPos := info.CurrentPosition
	isEncode := myPos == cluerPos
	/*if points >= 20 { // 2

	} else */if false && points >= 13 { // 1
		//codeN := 10 //3 * 3 /* useless / playable / not playable && not useless */
		code := 0
		for i, cards := range info.PlayerCards {
			if i == cluerPos || i == myPos {
				continue
			}
			variants := make(Variants, len(cards))
			for j, _ := range cards {
				variants[j] = &Variant{
					Count: len(ai.GetPlayableVariants(info, ai.Knowledge[i][j])),
					Pos:   j,
				}
			}
			sort.Sort(variants)

			for j := 0; j < 2; j++ {
				variant := variants[j]
				val1, val2 := cards[variant.Pos].Value, info.TableCards[cards[variant.Pos].Color].Value
				if val1 == val2+1 {
					code += Pow(3, j)
				} else if val1 <= val2 {
					code += 2 * Pow(3, j)
				}
			}
		}
		actions := ai.GetClueActions(info, cluerPos, info.Step)
		return actions[code%len(actions)]
	} else { // 0
		code := 0
		for i, cards := range info.PlayerCards {
			if i == cluerPos || i == myPos {
				continue
			}
			variants := make(Variants, len(cards))
			for j, _ := range cards {
				variants[j] = &Variant{
					Count: len(ai.GetPlayableVariants(info, ai.Knowledge[i][j])),
					Pos:   j,
				}
			}
			sort.Sort(variants)
			for j := 0; j < 3; j++ {
				variant := variants[j]
				if cards[variant.Pos].Value == info.TableCards[cards[variant.Pos].Color].Value+1 {
					code += 1 << uint(j)
				}
			}
		}

		actions := ai.GetClueActions(info, cluerPos, info.Step)
		if isEncode {
			fmt.Println("Return Clue:", cluerPos, actions[code%len(actions)])
			return actions[code%len(actions)]
		}
		for i, action := range actions {
			if action.Equal(&ai.History[info.Step]) {
				code += i
				break
			}
		}
		code = code % len(actions)
		cards := info.PlayerCards[myPos]
		variants := make(Variants, len(cards))
		for j, _ := range cards {
			variants[j] = &Variant{
				Count: len(ai.GetPlayableVariants(info, ai.Knowledge[myPos][j])),
				Pos:   j,
			}
		}
		sort.Sort(variants)
		for j := 0; j < 3; j++ {
			if code&1<<uint(j) != 0 {
				ai.Knowledge[myPos][j] = ai.GetPlayableVariants(info, ai.Knowledge[myPos][j])
			} else {
				ai.Knowledge[myPos][j] = ai.GetUnplayableVariants(info, ai.Knowledge[myPos][j])
			}
		}
		return nil
	}
}

func (ai *AI9) Decode2PlayerClue(oldInfo, info *game.PlayerGameInfo) {
	myPos := info.Position
	fmt.Println("DecodeClue")
	action := ai.History[oldInfo.Step]
	ai.Encode2PlayerClue(oldInfo)
	if action.PlayerPosition == myPos {
		fmt.Println("HELLO WORLD")
		for i, _ := range info.PlayerCards[myPos] {
			card := &info.PlayerCards[myPos][i]
			if action.ActionType == game.TypeActionInformationColor {
				if card.KnownColor && card.Color == game.CardColor(action.Value) {
					ai.Knowledge[myPos][i].RemoveWithoutColor(game.CardColor(action.Value))
				} else {
					ai.Knowledge[myPos][i].RemoveWithColor(game.CardColor(action.Value))
				}
			} else {
				if card.KnownValue && card.Value == game.CardValue(action.Value) {
					ai.Knowledge[myPos][i].RemoveWithoutValue(game.CardValue(action.Value))
				} else {
					ai.Knowledge[myPos][i].RemoveWithValue(game.CardValue(action.Value))
				}
			}
		}
	} else {
		fmt.Println(action)
		panic(action)
	}
	fmt.Println(ai.Knowledge[myPos])
}

func (ai *AI9) Encode3PlayerClue(info *game.PlayerGameInfo) *game.Action {
	return nil
}

func (ai *AI9) Decode3PlayerClue(oldInfo, info *game.PlayerGameInfo) {

}

func (ai *AI9) Encode4PlayerClue(info *game.PlayerGameInfo) *game.Action {
	return nil
}

func (ai *AI9) Decode4PlayerClue(oldInfo, info *game.PlayerGameInfo) {

}

func (ai *AI9) Encode5PlayerClue(info *game.PlayerGameInfo) *game.Action {
	return nil
}

func (ai *AI9) Decode5PlayerClue(oldInfo, info *game.PlayerGameInfo) {

}

func (ai *AI9) ClueOnStep(info *game.PlayerGameInfo, chooseCardPos func(int, int) int) *game.Action {
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

func (ai *AI9) Clue(info *game.PlayerGameInfo) *game.Action {
	return map[int](func(*game.PlayerGameInfo) *game.Action){
		2: ai.Encode2PlayerClue,
		3: ai.Encode3PlayerClue,
		4: ai.Encode4PlayerClue,
		5: ai.Encode5PlayerClue,
	}[info.PlayerCount](info)
}

func (ai *AI9) NewCard(info *game.PlayerGameInfo, pos int) CardKnowledge {
	// @todo
	return ai.NewCardKnowledge()
}

func (ai *AI9) DecodeClues(info *game.PlayerGameInfo) {
	for i := Max(-1, info.Step-info.PlayerCount) + 1; i < info.Step; i++ {
		action := &ai.History[i]
		oldInfo := ai.Informator.GetPlayerState(i)
		if action.IsInfoAction() {
			ai.DecodeClue(&oldInfo, info)
		} else {
			pos := action.PlayerPosition
			val := action.Value
			if pos == info.Position {
				panic("Bad action")
			}
			knowledge := ai.Knowledge[pos]
			knowledge = append(knowledge[:val], knowledge[val+1:]...)
			if len(knowledge) < len(info.PlayerCards[pos]) {
				knowledge = append(knowledge, ai.NewCard(&oldInfo, pos))
			}
			ai.Knowledge[pos] = knowledge
		}
	}
}

func (ai *AI9) DecodeClue(oldInfo, info *game.PlayerGameInfo) {
	map[int](func(*game.PlayerGameInfo, *game.PlayerGameInfo)){
		2: ai.Decode2PlayerClue,
		3: ai.Decode3PlayerClue,
		4: ai.Decode4PlayerClue,
		5: ai.Decode5PlayerClue,
	}[info.PlayerCount](oldInfo, info)
}

func (ai *AI9) TryDiscard(info *game.PlayerGameInfo, pos ...int) *game.Action {
	myPos := info.CurrentPosition
	if len(pos) == 1 {
		myPos = pos[0]
	}
	myCards := info.PlayerCards[myPos]
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

func (ai *AI9) Discard(info *game.PlayerGameInfo) *game.Action {
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

func (ai *AI9) InitializeCache() {
	myFirstInfo := ai.Informator.GetPlayerState(0)
	ai.Knowledge = make([][]CardKnowledge, myFirstInfo.PlayerCount)
	for i := 0; i < myFirstInfo.PlayerCount; i++ {
		ai.Knowledge[i] = make([]CardKnowledge, len(myFirstInfo.PlayerCards[i]))
		for j := 0; j < len(myFirstInfo.PlayerCards[i]); j++ {
			ai.Knowledge[i][j] = ai.NewCardKnowledge()
		}
	}
	ai.Informator.SetCache([]interface{}{ai.Knowledge, myFirstInfo.Copy()})
}

func (ai *AI9) LoadInformation(info *game.PlayerGameInfo) {
	//myPos := info.Position
	myCachedData := ai.Informator.GetCache().([]interface{})
	ai.Knowledge = myCachedData[0].([][]CardKnowledge)
	myOldInfo := myCachedData[1].(*game.PlayerGameInfo)
	if myOldInfo == nil || myOldInfo.Step != Max(0, info.Step-info.PlayerCount) {
		panic("Bad Cache")
	}

	/*if info.Step >= info.PlayerCount {
		myCards := info.PlayerCards[myPos]
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
		}
	}*/
	ai.DecodeClues(info)
	ai.Informator.SetCache([]interface{}{ai.Knowledge, info.Copy()})
}

func (ai *AI9) GetProgress(info *game.PlayerGameInfo) map[game.CardColor]game.CardValue {
	progress := map[game.CardColor]game.CardValue{}
	for color, card := range info.TableCards {
		progress[color] = card.Value
	}
	return progress
}

func (ai *AI9) GetPoints(progress map[game.CardColor]game.CardValue) int {
	points := 0
	for _, value := range progress {
		points += int(value)
	}
	return points
}

func (ai *AI9) GetPlayableUnknownCards(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) (int, int) {
	cards := info.PlayerCards[pos]
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

func (ai *AI9) HardDiscard(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) *game.Action {
	_, topIdx := ai.GetPlayableUnknownCards(info, progress, pos)
	return game.NewAction(game.TypeActionPlaying, pos, topIdx)
}

func (ai *AI9) RiskyPlay(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) *game.Action {
	topIdx, _ := ai.GetPlayableUnknownCards(info, progress, pos)
	return game.NewAction(game.TypeActionPlaying, pos, topIdx)
}

func (ai *AI9) DiscardHighest(info *game.PlayerGameInfo, pos int) *game.Action {
	highestIdx := 0
	cards := info.PlayerCards[pos]
	for i := 0; i < len(cards); i++ {
		if cards[i].Value > cards[highestIdx].Value {
			highestIdx = i
		}
	}
	return game.NewAction(game.TypeActionDiscard, pos, highestIdx)
}

func (ai *AI9) GetHardAction(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos, deep int) *game.Action {
	var action *game.Action
	//myCards := info.PlayerCards[pos]
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
	if deep == ai.GetMaxDeep(info) {
		fmt.Println("Critical discard")
	}
	action = ai.DiscardHighest(info, pos)
	return action
}

func (ai *AI9) GetBestAction(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos, deep int) (*game.Action, int) {
	for i := 0; i < len(info.PlayerCards[pos]); i++ {
		variants := ai.Knowledge[pos][i]
		isPlayable := true
		for cv, _ := range variants {
			isPlayable = isPlayable && progress[cv.Color]+1 == cv.Value
		}
		if isPlayable {
			return game.NewAction(game.TypeActionPlaying, pos, i), 0
		}
	}
	if info.BlueTokens > 0 {
		return ai.Clue(info), 0
	}
	return game.NewAction(game.TypeActionDiscard, pos, 0), 0
}

func (ai *AI9) GetMaxDeep(info *game.PlayerGameInfo) int {
	deep := info.PlayerCount
	if info.DeckSize == 0 {
		deep = info.MaxStep - info.Step
	}
	if deep == 0 {
		panic("Bad deep")
	}
	return deep
}

func (ai *AI9) FindBestAction(info *game.PlayerGameInfo) *game.Action {
	myPos := info.CurrentPosition
	progress := ai.GetProgress(info)
	action, _ := ai.GetBestAction(info, progress, myPos, ai.GetMaxDeep(info))
	fmt.Println(action)
	return action
}

func (ai *AI9) GetAction() *game.Action {
	info := &ai.PlayerInfo
	fmt.Println("Step:", info.Step)
	if info.Step < info.PlayerCount {
		ai.InitializeCache()
	}
	ai.LoadInformation(info)
	return ai.FindBestAction(info)
}
