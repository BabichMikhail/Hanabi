package ai

import (
	"fmt"
	"math"
	"sort"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI9 struct {
	BaseAI
	Knowledge [][]CardKnowledge
	Deep      int
}

func NewAI9(baseAI *BaseAI) *AI9 {
	ai := new(AI9)
	ai.BaseAI = *baseAI
	ai.Deep = ai.PlayerInfo.PlayerCount
	return ai
}

/* Clues */

// @todo use filter for actions
func (ai *AI9) FilterClueActions(info *game.PlayerGameInfo, actions []*game.Action, limit int) []*game.Action {
	if len(actions) == limit {
		return actions
	}

	for j := len(actions) - 1; j >= 0; j-- {
		for i := Max(0, info.Step-info.PlayerCount); len(actions) > limit && i < info.Step; i++ {
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

	for j := len(actions) - 1; j >= 0 && len(actions) > limit; j-- {
		action := actions[j]
		if action.ActionType == game.TypeActionInformationValue && action.Value <= lowestValue {
			actions = append(actions[:j], actions[j+1:]...)
		}
	}

	for j := len(actions) - 1; j >= 0 && len(actions) > limit; j-- {
		action := actions[j]
		if _, ok := completeColors[game.CardColor(action.Value)]; action.ActionType == game.TypeActionInformationColor && ok {
			actions = append(actions[:j], actions[j+1:]...)
		}
	}

	return actions
}

func (ai *AI9) GetClueActions(info *game.PlayerGameInfo, myPos int) []*game.Action {
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

	return actions
}

type Variants []*Variant
type Variant struct {
	MaxCount   int
	SumCount   int
	Counts     []int
	Pos        int
	OKIdx      int
	N          int
	Knowledges []CardKnowledge
}

func (v Variant) GetCardKnowledge() CardKnowledge {
	if len(v.Knowledges[v.OKIdx]) == 0 {
		panic("Bad choice")
	}
	return v.Knowledges[v.OKIdx]
}

func (v Variant) String() string {
	return fmt.Sprintf("(%d %d %d %d)", v.Pos, v.Counts[0], len(v.Counts), len(v.Knowledges))
}
func (v Variants) Len() int {
	return len(v)
}

func (v Variants) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Variants) Less(i, j int) bool {
	return v[i].Counts[0] > v[j].Counts[0] ||
		v[i].Counts[0] == v[j].Counts[0] && len(v[i].Counts) > len(v[j].Counts) ||
		v[i].Counts[0] == v[j].Counts[0] && len(v[i].Counts) == len(v[j].Counts) && v[i].Pos < v[j].Pos
}

type CardKnowledge map[game.ColorValue]struct{}

func (kn CardKnowledge) Copy() CardKnowledge {
	newKn := CardKnowledge{}
	for k, v := range kn {
		newKn[k] = v
	}
	return newKn
}

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
	for cv, _ := range kn {
		if cv.Color == color {
			delete(kn, cv)
		}
	}
}

func (kn CardKnowledge) RemoveWithoutColor(color game.CardColor) {
	for cv, _ := range kn {
		if cv.Color != color {
			delete(kn, cv)
		}
	}
}

func (kn CardKnowledge) RemoveWithValue(value game.CardValue) {
	for cv, _ := range kn {
		if cv.Value == value {
			delete(kn, cv)
		}
	}
}

func (kn CardKnowledge) RemoveWithoutValue(value game.CardValue) {
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

func (ai *AI9) GetOkIdx(knowledges []CardKnowledge, card *game.Card) int {
	okIdx := -1
	if card.KnownColor && card.KnownValue {
		for i, kn := range knowledges {
			for cv, _ := range kn {
				if cv.Color == card.Color && cv.Value == card.Value {
					okIdx = i
					break
				}
			}
			if okIdx != -1 {
				break
			}
		}
	}
	return okIdx
}

func (ai *AI9) EncodeCard2Variants(info *game.PlayerGameInfo, knowledge CardKnowledge, card *game.Card, pos int) *Variant {
	knowledges := []CardKnowledge{
		ai.GetPlayableVariants(info, knowledge),
		ai.GetUnplayableVariants(info, knowledge),
	}

	if len(knowledges[0]) == 0 || len(knowledges[1]) == 0 {
		var knowledge1, knowledge2 CardKnowledge
		if len(knowledges[0]) == 0 {
			knowledge1 = knowledges[0]
			knowledge2 = knowledges[1]
		} else {
			knowledge1 = knowledges[1]
			knowledge2 = knowledges[0]
		}

		values := sort.IntSlice{}
		for cv, _ := range knowledge2 {
			values = append(values, int(game.HashColorValue(cv.Color, cv.Value)))
		}
		sort.Sort(values)

		for _, hashVal := range values {
			if math.Abs(float64(len(knowledge1)-len(knowledge2))) < 2 {
				break
			}
			color, value := game.ColorValueByHashColorValue(game.HashValue(hashVal))
			cv := game.ColorValue{Color: color, Value: value}
			knowledge1[cv] = struct{}{}
			delete(knowledge2, cv)
		}
	}

	variant := Variant{
		MaxCount:   Max(len(knowledges[0]), len(knowledges[1])),
		SumCount:   len(knowledge),
		Counts:     []int{len(knowledges[0]), len(knowledges[1])},
		Pos:        pos,
		OKIdx:      ai.GetOkIdx(knowledges, card),
		N:          2,
		Knowledges: knowledges,
	}
	return &variant
}

func (ai *AI9) EncodeCard3Variants(info *game.PlayerGameInfo, knowledge CardKnowledge, card *game.Card, pos int) *Variant {
	if len(knowledge) < 3 {
		return ai.EncodeCard2Variants(info, knowledge, card, pos)
	}

	useful := CardKnowledge{}
	playable := CardKnowledge{}
	useless := CardKnowledge{}
	for cv, _ := range knowledge {
		tval := info.TableCards[cv.Color].Value
		cval := cv.Value
		if tval+1 < cval {
			useful[cv] = struct{}{}
		} else if tval+1 == cval {
			playable[cv] = struct{}{}
		} else if tval+1 > cval {
			useless[cv] = struct{}{}
		}
	}

	if len(useless) == 0 && len(useful) > 0 {
		val := 0
		for len(useless) == 0 {
			if val > 5 {
				panic("Bad value")
			}
			for cv, _ := range useful {
				if int(info.TableCards[cv.Color].Value)+val == int(cv.Value) {
					useless[cv] = struct{}{}
					delete(useful, cv)
				}
			}
			val++
		}
	}

	updateKnowledges := func(knowledge1, knowledge2 CardKnowledge) {
		if len(knowledge1) == 0 && len(knowledge2) > 1 {
			val := 1
			for len(knowledge1) == 0 {
				if val > 5 {
					panic("Bad value")
				}
				for _, color := range game.ColorsTable {
					for cv, _ := range knowledge2 {
						if math.Abs(float64(len(knowledge1)-len(knowledge2))) < 2 {
							break
						}
						if cv.Color != color || val != int(cv.Value) {
							continue
						}
						knowledge1[cv] = struct{}{}
						delete(knowledge2, cv)
					}
				}
				val++
			}
		}
	}

	updateKnowledges(useful, playable)
	updateKnowledges(playable, useful)

	knowledges := []CardKnowledge{}
	counts := []int{}

	if len(useless) > 0 {
		knowledges = append(knowledges, useless)
		counts = append(counts, len(useless))
	}

	if len(playable) > 0 {
		knowledges = append(knowledges, playable)
		counts = append(counts, len(playable))
	}

	if len(useful) > 0 {
		knowledges = append(knowledges, useful)
		counts = append(counts, len(useful))
	}

	maxCount := Max(len(useful), Max(len(playable), len(useless)))

	return &Variant{
		MaxCount:   maxCount,
		SumCount:   len(knowledge),
		Counts:     counts,
		Knowledges: knowledges,
		N:          len(knowledges),
		OKIdx:      ai.GetOkIdx(knowledges, card),
		Pos:        pos,
	}
}

func (ai *AI9) EncodeCard20Variants(info *game.PlayerGameInfo, knowledge CardKnowledge, card *game.Card, pos int) *Variant {
	counts := []int{}
	knowledges := []CardKnowledge{}
	var maxCount int

	usefulCount := 0
	values := sort.IntSlice{}
	useless := CardKnowledge{}
	useful := CardKnowledge{}
	for cv, _ := range knowledge {
		if cv.Color == game.NoneColor || cv.Value == game.NoneValue {
			panic("Bad color or value")
		}
		if info.TableCards[cv.Color].Value < cv.Value {
			usefulCount++
			values = append(values, int(game.HashColorValue(cv.Color, cv.Value)))
		} else {
			useless[cv] = struct{}{}
		}

	}
	sort.Sort(values)

	maxUsefulCount := 18
	if len(useless) == 0 {
		maxUsefulCount = 19
	}

	for _, hashVal := range values {
		if len(knowledges) == maxUsefulCount {
			break
		}
		color, value := game.ColorValueByHashColorValue(game.HashValue(hashVal))
		knowledges = append(knowledges, CardKnowledge{game.ColorValue{Color: color, Value: value}: struct{}{}})
		counts = append(counts, 1)
	}

	for i := maxUsefulCount; i < len(values); i++ {
		color, value := game.ColorValueByHashColorValue(game.HashValue(values[i]))
		useful[game.ColorValue{Color: color, Value: value}] = struct{}{}
	}

	maxCount = Max(1, Max(len(useless), len(useful)))

	if len(useless) > 0 {
		knowledges = append(knowledges, useless)
		counts = append(counts, len(useless))
	}

	knowledges = append(knowledges, useful)
	counts = append(counts, len(useful))

	return &Variant{
		MaxCount:   maxCount,
		SumCount:   len(knowledge),
		Counts:     counts,
		Knowledges: knowledges,
		N:          len(knowledges),
		OKIdx:      ai.GetOkIdx(knowledges, card),
		Pos:        pos,
	}
}

func (ai *AI9) EncodeCard25Variants(info *game.PlayerGameInfo, knowledge CardKnowledge, card *game.Card, pos int) *Variant {
	counts := []int{}
	knowledges := []CardKnowledge{}
	useless := CardKnowledge{}

	values := sort.IntSlice{}
	for cv, _ := range knowledge {
		if info.TableCards[cv.Color].Value < cv.Value {
			values = append(values, int(game.HashColorValue(cv.Color, cv.Value)))
		} else {
			useless[cv] = struct{}{}
		}
	}

	sort.Sort(values)
	for i := 0; i < len(values); i++ {
		color, value := game.ColorValueByHashColorValue(game.HashValue(values[i]))
		knowledges = append(knowledges, CardKnowledge{
			game.ColorValue{Color: color, Value: value}: struct{}{},
		})
		counts = append(counts, 1)
	}

	if length := len(useless); length > 0 {
		knowledges = append(knowledges, useless)
		counts = append(counts, length)
	}

	return &Variant{
		MaxCount:   1,
		SumCount:   len(knowledge),
		Counts:     counts,
		Knowledges: knowledges,
		N:          len(knowledges),
		OKIdx:      ai.GetOkIdx(knowledges, card),
		Pos:        pos,
	}
}

func (ai *AI9) EncodeCard(info *game.PlayerGameInfo, knowledge CardKnowledge, card *game.Card, pos, nvariants int) *Variant {
	switch nvariants {
	case 2:
		return ai.EncodeCard2Variants(info, knowledge, card, pos)
	case 3:
		return ai.EncodeCard3Variants(info, knowledge, card, pos)
	case 20:
		return ai.EncodeCard20Variants(info, knowledge, card, pos)
	case 25:
		return ai.EncodeCard25Variants(info, knowledge, card, pos)
	default:
		panic("Not implemented")
	}
}

func (ai *AI9) EncodeCardsNVariants(info *game.PlayerGameInfo, knowledge []CardKnowledge, cards []game.Card, variantsToCard int, updateKnowledge bool) (int, Variants) {
	code := 0
	variants := make(Variants, len(cards))
	for j, _ := range cards {
		card := &cards[j]
		variants[j] = ai.EncodeCard(info, knowledge[j], card, j, variantsToCard)
	}
	sort.Sort(variants)
	n := ai.GetClueCount()
	k := 1
	for j := 0; j < len(variants) && n >= variants[j].N; j++ {
		variant := variants[j]
		n /= variant.N
		if variant.OKIdx == -1 {
			break
		}
		if updateKnowledge {
			knowledge[variant.Pos] = variant.GetCardKnowledge()
		}
		code += variant.OKIdx * k
		k *= variant.N
	}
	// @todo encode/decode if n >= 2

	return code, variants
}

func (ai *AI9) EncodeCards(info *game.PlayerGameInfo, cluerPos, cluedPos int, knowledge []CardKnowledge, cards []game.Card, nvariants int, updateKnowledge bool) (int, Variants) {
	availableCards := info.GetUnplayedCards()
	for _, cardKnowledge := range knowledge {
		for cv, _ := range cardKnowledge {
			if count, ok := availableCards[cv]; !ok || count == 0 {
				delete(cardKnowledge, cv)
			}
		}
	}

	switch nvariants {
	case 10:
		return ai.EncodeCardsNVariants(info, knowledge, cards, 3, updateKnowledge)
	case 20:
		return ai.EncodeCardsNVariants(info, knowledge, cards, 3, updateKnowledge)
	case 30:
		return ai.EncodeCardsNVariants(info, knowledge, cards, 25, updateKnowledge)
	case 40:
		if info.Step >= 45 {
			return ai.EncodeCardsNVariants(info, knowledge, cards, 3, updateKnowledge)
		}
		return ai.EncodeCardsNVariants(info, knowledge, cards, 25, updateKnowledge)
	default:
		panic("Not implemented")
	}
}

func (ai *AI9) GetClueCount() int {
	return (ai.PlayerInfo.PlayerCount - 1) * 10
}

func (ai *AI9) Clue(info *game.PlayerGameInfo) *game.Action {
	myPos := info.Position
	cluerPos := info.CurrentPosition
	if cluerPos != myPos {
		panic("I can't encode clue")
	}

	code := 0
	nvariants := ai.GetClueCount()
	for i, cards := range info.PlayerCards {
		if i == myPos {
			continue
		}
		dcode, _ := ai.EncodeCards(info, myPos, i, ai.Knowledge[i], cards, nvariants, false)
		code += dcode
	}

	actions := ai.GetClueActions(info, cluerPos)
	action := actions[code%len(actions)]
	return action
}

func (ai *AI9) NewCard(info *game.PlayerGameInfo, pos int) CardKnowledge {
	return ai.NewCardKnowledge()
}

func (ai *AI9) DecodeClues(info *game.PlayerGameInfo) {
	for i := Max(0, info.Step-info.PlayerCount); i < info.Step; i++ {
		action := &ai.History[i]
		oldInfo := ai.Informator.GetPlayerState(i)
		if action.IsInfoAction() {
			ai.DecodeClue(&oldInfo)
		} else {
			pos := action.PlayerPosition
			val := action.Value
			knowledge := ai.Knowledge[pos]
			knowledge = append(knowledge[:val], knowledge[val+1:]...)
			if len(knowledge) < len(info.PlayerCards[pos]) {
				knowledge = append(knowledge, ai.NewCard(&oldInfo, pos))
			}
			ai.Knowledge[pos] = knowledge
		}
	}
}

func (ai *AI9) DecodeClue(info *game.PlayerGameInfo) {
	myPos := info.Position
	cluerPos := info.CurrentPosition
	code := 0
	clueCount := ai.GetClueCount()
	for i, cards := range info.PlayerCards {
		if i == cluerPos || i == myPos {
			continue
		}
		dcode, _ := ai.EncodeCards(info, cluerPos, i, ai.Knowledge[i], cards, clueCount, true)
		code += dcode
	}

	actions := ai.GetClueActions(info, cluerPos)
	resultCode := -1
	for i, action := range actions {
		if action.Equal(&ai.History[info.Step]) {
			resultCode = i
			break
		}
	}

	for i := 0; i < len(actions); i++ {
		if (code+i)%len(actions) == resultCode {
			code = i
			break
		}
	}

	action := &ai.History[info.Step]

	if myPos == cluerPos {
		pos := action.PlayerPosition
		cards := info.PlayerCards[pos]
		for i, _ := range cards {
			card := &cards[i]
			if action.ActionType == game.TypeActionInformationColor {
				if color := game.CardColor(action.Value); card.KnownColor && card.Color == color {
					ai.Knowledge[pos][i].RemoveWithoutColor(color)
				} else {
					ai.Knowledge[pos][i].RemoveWithColor(color)
				}
			} else {
				if value := game.CardValue(action.Value); card.KnownValue && card.Value == value {
					ai.Knowledge[pos][i].RemoveWithoutValue(value)
				} else {
					ai.Knowledge[pos][i].RemoveWithValue(value)
				}
			}
		}
		return
	}

	code = code % len(actions)
	myCards := info.PlayerCards[myPos]
	_, variants := ai.EncodeCards(info, cluerPos, myPos, ai.Knowledge[myPos], myCards, clueCount, false)

	sort.Sort(variants)
	k := 1
	high := -1
	for j := 0; j < len(variants); j++ {
		k *= variants[j].N
		if k > clueCount {
			k /= variants[j].N
			break
		}
		high = j
	}
	for j := high; j >= 0; j-- {
		variant := variants[j]
		k /= variant.N
		idx := code / k
		code = code % k
		if variant.OKIdx == -1 {
			variant.OKIdx = idx
		}

		ai.Knowledge[myPos][variant.Pos] = variant.GetCardKnowledge()
	}

	pos := action.PlayerPosition
	cards := info.PlayerCards[pos]
	if action.PlayerPosition == info.Position {
		nextInfo := ai.Informator.GetPlayerState(info.Step + 1)
		cards = nextInfo.PlayerCards[pos]
	}

	for i, _ := range cards {
		card := &cards[i]
		knowledge := ai.Knowledge[pos][i]
		if action.ActionType == game.TypeActionInformationColor {
			if color := game.CardColor(action.Value); card.KnownColor && card.Color == color {
				knowledge.RemoveWithoutColor(color)
			} else {
				knowledge.RemoveWithColor(color)
			}
		} else {
			if value := game.CardValue(action.Value); card.KnownValue && card.Value == value {
				knowledge.RemoveWithoutValue(value)
			} else {
				knowledge.RemoveWithValue(value)
			}
		}
	}
}

/* End Clues */
/* AI Utils */

func (ai *AI9) LoadInformation(info *game.PlayerGameInfo) {
	myCachedData := ai.Informator.GetCache().([]interface{})
	ai.Knowledge = myCachedData[0].([][]CardKnowledge)
	myOldInfo := myCachedData[1].(*game.PlayerGameInfo)
	if myOldInfo == nil || myOldInfo.Step != Max(0, info.Step-info.PlayerCount) {
		panic("Bad Cache")
	}

	ai.DecodeClues(info)
	ai.Informator.SetCache([]interface{}{ai.Knowledge, info.Copy()})
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

/* Action */

func (ai *AI9) TryDiscard(info *game.PlayerGameInfo, pos ...int) *game.Action {
	myPos := info.CurrentPosition
	if len(pos) == 1 {
		myPos = pos[0]
	}
	myCards := info.PlayerCards[myPos]
	for i, card := range myCards { /* useless */
		isUseless := true
		for hashVal, _ := range card.ProbabilityCard {
			color, value := game.ColorValueByHashColorValue(hashVal)
			isUseless = isUseless && info.TableCards[color].Value >= value
		}
		if isUseless {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
	}

	cardCounts := map[game.ColorValue]int{}
	for i, card := range myCards { /* duplicate */
		if !card.KnownColor || !card.KnownValue {
			continue
		}
		colorValue := game.ColorValue{Color: card.Color, Value: card.Value}
		count := cardCounts[colorValue]
		if count == 1 {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
		cardCounts[colorValue] = 1
	}

	for i, card := range myCards { /* not last card */
		unplayedCards := info.GetUnplayedCards()
		isNotLast := true
		for hashVal, _ := range card.ProbabilityCard {
			color, value := game.ColorValueByHashColorValue(hashVal)
			isNotLast = isNotLast && unplayedCards[game.ColorValue{Color: color, Value: value}] > 1
		}
		if isNotLast {
			return game.NewAction(game.TypeActionDiscard, myPos, i)
		}
	}
	return nil
}

func (ai *AI9) DiscardHighest(info *game.PlayerGameInfo, pos int) *game.Action {
	highestValue := 0.0
	cardPos := -1
	for i, card := range info.PlayerCards[pos] {
		sum := 0.0
		for hashVal, prob := range card.ProbabilityCard {
			_, value := game.ColorValueByHashColorValue(hashVal)
			sum += float64(value) * prob
		}

		if sum > highestValue {
			highestValue = sum
			cardPos = i
		}
	}
	return game.NewAction(game.TypeActionDiscard, pos, cardPos)
}

func (ai *AI9) HardDiscard(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) *game.Action {
	cardPos := -1
	bestUseless := 0.0
	cards := info.PlayerCards[pos]
	for i, card := range cards {
		probUseless := 0.0
		for hashVal, prob := range card.ProbabilityCard {
			color, value := game.ColorValueByHashColorValue(hashVal)
			if progress[color] >= value {
				probUseless += prob
			}
		}

		if cardPos == -1 || probUseless > bestUseless {
			cardPos = i
			bestUseless = probUseless
		}
	}
	return game.NewAction(game.TypeActionDiscard, pos, cardPos)
}

func (ai *AI9) RiskyPlay(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) *game.Action {
	cardPos := -1
	bestPlayable := 0.0
	for i, card := range info.PlayerCards[pos] {
		probPlayable := 0.0
		for hashVal, prob := range card.ProbabilityCard {
			color, value := game.ColorValueByHashColorValue(hashVal)
			if progress[color]+1 == value {
				probPlayable += prob
			}
		}
		if cardPos == -1 || probPlayable > bestPlayable {
			cardPos = i
			bestPlayable = probPlayable
		}
	}
	return game.NewAction(game.TypeActionPlaying, pos, cardPos)
}

func (ai *AI9) GetHardAction(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos, deep int) *game.Action {
	var action *game.Action
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
	return action
}

func (ai *AI9) GetMathExpectedValuesForPlayability(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos int) float64 {
	myCards := info.PlayerCards[pos]
	result := 0.0
	for _, card := range myCards {
		result += card.GetPlayability(progress)
	}
	return result
}

func (ai *AI9) GetBestAction(info *game.PlayerGameInfo, progress map[game.CardColor]game.CardValue, pos, deep int) (*game.Action, int) {
	if deep == 1 {
		myCards := info.PlayerCards[pos]
		for i := 0; i < len(myCards); i++ {
			card := &myCards[i]
			if card.IsCardPlayable(progress) {
				points := ai.GetPoints(progress) + 1
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
			if info.RedTokens < 2 {
				bestPlayability := -0.1
				idx := -1
				for i := 0; i < len(myCards); i++ {
					c := &myCards[i]
					criticality := c.GetCriticality(progress, info.VariantsCount)
					playability := c.GetPlayability(progress)
					if playability > bestPlayability && criticality < 0.25 {
						bestPlayability = playability
						idx = i
					}
				}

				if bestPlayability > 0.50 {
					return game.NewAction(game.TypeActionPlaying, pos, idx), ai.GetPoints(progress) + 1
				}
			}

			action = ai.Clue(info)
		}
		return action, ai.GetPoints(progress)
	}

	nextPos := (pos + 1) % info.PlayerCount
	topPoints := -1
	var topAction *game.Action

	for i := 0; i < len(info.PlayerCards[pos]); i++ {
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
		_, points := ai.GetBestAction(info, progress, nextPos, deep-1)
		if points > topPoints {
			topPoints = points
			topAction = ai.Clue(info)
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

func (ai *AI9) GetBestResult() (*game.Action, int) {
	info := &ai.PlayerInfo
	if info.Step < info.PlayerCount {
		ai.InitializeCache()
	}
	ai.LoadInformation(info)
	ai.KnowledgeToMyCards(info)
	ai.SetProbabilities(info)

	myPos := info.CurrentPosition
	progress := ai.GetProgress(info)
	action, points := ai.GetBestAction(info, progress, myPos, ai.Deep)
	return action, points
}

func (ai *AI9) SetCard(info *game.PlayerGameInfo, card *game.Card, cardVariants []game.Variants) *game.CardProbs {
	newCardRef := &game.CardProbs{
		Card:  card,
		Probs: make([]float64, len(cardVariants)),
	}

	count := 0.0
	for color, _ := range card.ProbabilityColors {
		card.ProbabilityColors[color] = 1.0
	}

	for value, _ := range card.ProbabilityValues {
		card.ProbabilityValues[value] = 1.0
	}

	for k := 0; k < len(cardVariants); k++ {
		color, value := cardVariants[k].Color, cardVariants[k].Value
		_, colorOK := card.ProbabilityColors[color]
		_, valueOK := card.ProbabilityValues[value]
		_, knowledgeOK := card.ProbabilityCard[game.HashColorValue(color, value)]
		if colorOK && valueOK && knowledgeOK {
			count++
			newCardRef.Probs[k] = 1.0
		} else {
			newCardRef.Probs[k] = 0.0
		}
	}

	for k := 0; k < len(cardVariants); k++ {
		if newCardRef.Probs[k] == 1.0 {
			newCardRef.Probs[k] /= count
		}
	}
	return newCardRef
}

func (ai *AI9) SetProbabilities(info *game.PlayerGameInfo) {
	cardsRef := game.Cards{}
	info.SetVariantsCount(false, false)
	cardVariants := info.GetCardVariants()
	pos := info.CurrentPosition
	cards := info.PlayerCards[pos]
	for idx, _ := range cards {
		card := &cards[idx]
		if card.KnownColor && card.KnownValue {
			card.ProbabilityCard[game.HashColorValue(card.Color, card.Value)] = 1.0
			continue
		}
		cardsRef = append(cardsRef, *ai.SetCard(info, card, cardVariants))
	}

	for idx, _ := range info.Deck {
		card := &info.Deck[idx]
		cardsRef = append(cardsRef, *info.SetCard(card, cardVariants))
	}
	info.SetProbabilities_ConvergenceOfProbability(cardsRef, cardVariants)
}

func (ai *AI9) KnowledgeToMyCards(info *game.PlayerGameInfo) {
	pos := info.CurrentPosition
	cards := info.PlayerCards[pos]
	for i, _ := range cards {
		card := &cards[i]
		card.EmptyProbabilities()
		cardKnowledge := ai.Knowledge[pos][i]
		var cv game.ColorValue
		for cv, _ = range cardKnowledge {
			color, value := cv.Color, cv.Value
			card.ProbabilityCard[game.HashColorValue(color, value)] = 0.0
		}

		if len(cardKnowledge) == 1 {
			card.SetColor(cv.Color)
			card.SetValue(cv.Value)
		}
	}
}

func (ai *AI9) GetAction() *game.Action {
	action, _ := ai.GetBestResult()
	return action
}
