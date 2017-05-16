package ai

import (
	"github.com/BabichMikhail/Hanabi/game"
)

/*
Original:
	https://github.com/chikinn/hanabi/blob/master/players/hat_player.py
	https://github.com/chikinn/hanabi/blob/master/doc_hat_player.md
*/

type HatPlayerRecord struct {
	ImClued           bool
	FirstClued        int
	LastClued         int
	ClueValue         int
	NextPlayerActions []int
	/* pair (target, value) */
	LaterClues       [][]int
	LastCluedAnyClue int

	FutureProgress  map[game.CardColor]game.CardValue
	FutureDeckSize  int
	FutureHints     int
	FutureDiscarded []game.ColorValue

	CardsDrawn   int
	WillBePlayed []game.ColorValue
	HintChanges  []int

	MinFutureHints int
	MaxFutureHints int
}

type Hat struct {
	BaseAI
	PlayerRecords []*HatPlayerRecord
	MyRecord      *HatPlayerRecord
}

func NewAIHat(baseAI *BaseAI) *Hat {
	ai := new(Hat)
	ai.BaseAI = *baseAI
	return ai
}

func (rec *HatPlayerRecord) ResetMemory() {
	rec.ImClued = false
	rec.FirstClued = -1
	rec.LastClued = -1
	rec.ClueValue = -1
	rec.NextPlayerActions = []int{}
	rec.LaterClues = [][]int{}
}

func (ai *Hat) HatIsCardPlayable(card game.Card, progress map[game.CardColor]game.CardValue) bool {
	card.CheckVisible()
	return progress[card.Color]+1 == card.Value
}

func (ai *Hat) GetPlays(cards []game.Card, progress map[game.CardColor]game.CardValue) []game.Card {
	result := []game.Card{}
	for _, card := range cards {
		card.CheckVisible()
		if ai.HatIsCardPlayable(card, progress) {
			result = append(result, card)
		}
	}
	return result
}

func (ai *Hat) FindLowestIdx(cards []game.Card) int {
	minIdx := -1
	for i, card := range cards {
		card.CheckVisible()
		if minIdx == -1 || cards[minIdx].Value > card.Value {
			minIdx = i
		}
	}
	return minIdx
}

func (ai *Hat) CountUnplayedCards(progress map[game.CardColor]game.CardValue) int {
	count := 0
	for _, value := range progress {
		count += 5 - int(value)
	}
	return count
}

func (ai *Hat) GetPlayedCards(cards []game.Card, progress map[game.CardColor]game.CardValue) []game.Card {
	result := []game.Card{}
	for _, card := range cards {
		card.CheckVisible()
		if progress[card.Color] >= card.Value {
			result = append(result, card)
		}
	}
	return result
}

func (ai *Hat) GetDuplicateCards(cards []game.Card) []game.Card {
	colorValueCount := map[game.ColorValue]int{}
	for _, card := range cards {
		colorValueCount[game.ColorValue{Color: card.Color, Value: card.Value}]++
	}

	result := []game.Card{}
	for _, card := range cards {
		card.CheckVisible()
		if colorValueCount[game.ColorValue{Color: card.Color, Value: card.Value}] > 1 {
			result = append(result, card)
		}
	}
	return result
}

// @todo type game.Cards
func (ai *Hat) CardIndex(cards []game.Card, card game.Card) int {
	for i, card := range cards {
		card.CheckVisible()
		if card.Value == card.Value && card.Color == card.Color {
			return i
		}
	}
	panic("Bad cards")
}

func (ai *Hat) GetDiscardPile() []game.ColorValue {
	info := &ai.PlayerInfo
	result := make([]game.ColorValue, len(info.UsedCards))
	for i, card := range info.UsedCards {
		result[i].Color = card.Color
		result[i].Value = card.Value
	}
	for _, card := range info.TableCards {
		result = append(result, game.ColorValue{Color: card.Color, Value: card.Value})
	}
	return result
}

func (ai *Hat) GetNonvisibleCardsIdxs(cards []game.Card, names []game.ColorValue) []int {
	idxs := []int{}
	for i := 0; i < len(cards); i++ {
		ok := true
		for j := 0; j < len(names); j++ {
			if cards[i].Value == names[j].Value && cards[i].Color == names[j].Color {
				ok = false
				break
			}
		}
		if ok {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func (ai *Hat) GetProgress(playerInfo *game.PlayerGameInfo) map[game.CardColor]game.CardValue {
	progress := map[game.CardColor]game.CardValue{}
	for color, card := range playerInfo.TableCards {
		progress[color] = card.Value
	}
	return progress
}

func (ai *Hat) InterpretClue(me int) {
	info := &ai.PlayerInfo
	rec := ai.PlayerRecords[me]
	if !rec.ImClued || rec.FirstClued != info.CurrentPosition {
		return
	}

	n := len(info.PlayerCards)
	rec.NextPlayerActions = []int{}
	ai.InitializeFuturePrediction(rec, 0)
	ai.InitializeFutureClue(rec)
	i := rec.LastClued
	for i != me {
		x := ai.StandardPlay(info.PlayerCards[i], i, rec.WillBePlayed, ai.GetProgress(info), info.DeckSize, info.BlueTokens)
		rec.NextPlayerActions = append(rec.NextPlayerActions, x)
		rec.ClueValue = (rec.ClueValue - x + 9) % 9
		ai.FinalizeFutureAction(rec, x, i, me)
		i = (i - 1 + n) % n
	}
}

func (rec *HatPlayerRecord) ThinkOutOfTurn(me, player int, action *game.Action, n int) {
	if rec.LastCluedAnyClue == (player+n-1)%n {
		rec.LastCluedAnyClue = player
	}
	if rec.ImClued && rec.IsBetween(player, rec.FirstClued, rec.LastClued) {
		rec.ClueValue = (rec.ClueValue - rec.ActionToNumber(action) + 9) % 9
	}
	if action.ActionType != game.TypeActionInformationColor && action.ActionType != game.TypeActionInformationValue {
		return
	}

	target, value := action.PlayerPosition, action.Value
	if rec.ClueToNumber(value) == -1 {
		return
	}
	if rec.ImClued {
		rec.LastCluedAnyClue = target
		rec.LaterClues = append(rec.LaterClues, []int{target, rec.ClueToNumber(value)})
		return
	}
	rec.FirstClued = (rec.LastCluedAnyClue + 1) % n
	rec.LastCluedAnyClue = target
	if !rec.IsBetween(me, rec.FirstClued, target) {
		return
	}
	rec.ImClued = true
	rec.LastClued = target
	rec.ClueValue = rec.ClueToNumber(value)
	rec.LaterClues = [][]int{}
}

func (ai *Hat) StandardPlay(cards []game.Card, me int, dontPlay []game.ColorValue, progress map[game.CardColor]game.CardValue, decksize, hints int) int {
	info := ai.PlayerInfo
	if info.Position == me {
		panic("Stardard call on Me")
	}

	playableCards := ai.GetPlays(cards, progress)
	for i := len(playableCards) - 1; i >= 0; i-- {
		playableCards[i].CheckVisible()
		for j := 0; j < len(dontPlay); j++ {
			if playableCards[i].Color == dontPlay[j].Color && playableCards[i].Value == dontPlay[j].Value {
				playableCards = append(playableCards[:i], playableCards[i+1:]...)
				break
			}
		}
	}

	if len(playableCards) > 0 {
		idx := ai.FindLowestIdx(playableCards)
		for i := 0; i < len(cards); i++ {
			if playableCards[idx].Value == cards[i].Value && playableCards[idx].Color == cards[i].Color {
				return i
			}
		}
		panic("Missing Lowest")
	}

	if hints > 6 {
		return 8
	}

	if decksize-(hints-1)/3 < ai.CountUnplayedCards(progress) {
		return 8
	}

	x := ai.EasyDiscards(cards, dontPlay, progress)
	if x > 0 {
		return x
	}
	return 8
}

func (ai *Hat) ModifiedPlay(cards []game.Card, hinter int, player int, dontPlay []game.ColorValue, progress map[game.CardColor]game.CardValue, decksize, hints int) int {
	info := ai.PlayerInfo
	if info.CurrentPosition != hinter {
		panic("Magic")
	}
	if info.CurrentPosition == player {
		panic("Magic")
	}
	x := ai.StandardPlay(cards, player, dontPlay, progress, decksize, hints)
	rec := ai.MyRecord
	if x < 4 {
		if rec.MinFutureHints <= 0 && cards[x].Value != 5 {
			for i, card := range cards {
				card.CheckVisible()
				if ai.HatIsCardPlayable(card, progress) && card.Value == 5 {
					return i
				}
			}
		}
	}

	if rec.MaxFutureHints >= 8 {
		return 8
	}

	if x == 8 && (rec.MinFutureHints <= 0 || (rec.CardsDrawn == 0 && rec.FutureDeckSize == info.DeckSize)) {
		y := ai.EasyDiscards(cards, dontPlay, progress)
		if y > 0 {
			return y
		}

		y = ai.HardDiscards(cards, dontPlay, progress)
		if y > 0 {
			return y
		}
	}

	return x
}

func (ai *Hat) EasyDiscards(cards []game.Card, dontPlay []game.ColorValue, progress map[game.CardColor]game.CardValue) int {
	discardCards := ai.GetPlayedCards(cards, progress)
	if len(discardCards) > 0 {
		return ai.CardIndex(cards, discardCards[0]) + 4
	}

	discardCards = ai.GetDuplicateCards(cards)
	if len(discardCards) > 0 {
		return ai.CardIndex(cards, discardCards[0]) + 4
	}

	for i := 0; i < len(cards); i++ {
		for j := 0; j < len(dontPlay); j++ {
			if cards[i].Color == dontPlay[j].Color && cards[i].Value == dontPlay[j].Value {
				return i + 4
			}
		}
	}

	return 0
}

func (ai *Hat) HardDiscards(cards []game.Card, dontPlay []game.ColorValue, progress map[game.CardColor]game.CardValue) int {
	rec := ai.MyRecord
	discardsCardsIdxs := ai.GetNonvisibleCardsIdxs(cards, append(ai.GetDiscardPile(), rec.FutureDiscarded...))
	for i := len(discardsCardsIdxs) - 1; i >= 0; i-- {
		if cards[discardsCardsIdxs[i]].Value == 5 {
			discardsCardsIdxs = append(discardsCardsIdxs[:i], discardsCardsIdxs[i+1:]...)
		}
	}
	if len(discardsCardsIdxs) > 0 {
		highest := discardsCardsIdxs[0]
		for i := 1; i < len(discardsCardsIdxs); i++ {
			if cards[discardsCardsIdxs[i]].Value > cards[highest].Value {
				return highest + 4
			}
		}
	}
	return 0
}

func (rec *HatPlayerRecord) IsBetween(x, begin, end int) bool {
	return begin <= x && x <= end || end < begin && begin <= x || x <= end && end < begin
}

func (ai *Hat) NumberToAction(n int) *game.Action {
	myPos := ai.PlayerInfo.CurrentPosition
	if n < 4 {
		return game.NewAction(game.TypeActionPlaying, myPos, n)
	} else if n < 8 {
		return game.NewAction(game.TypeActionDiscard, myPos, n-4)
	}

	return nil
}

func (rec *HatPlayerRecord) ActionToNumber(play *game.Action) int {
	if play.ActionType == game.TypeActionInformationColor || play.ActionType == game.TypeActionInformationValue {
		return 8
	}
	d := 0
	if play.ActionType == game.TypeActionDiscard {
		d = 4
	}
	return play.Value + d
}

func (rec *HatPlayerRecord) ClueToNumber(clue int) int {
	if clue == 5 {
		return -1
	}
	if clue == 0 {
		panic("Magic")
	}
	if clue >= 1 && clue <= 4 {
		return clue - 1
	}
	number := clue - int(game.ColorsTable[0]) - 1
	if number >= 9 {
		panic("Magic")
	}
	return number
}

func (ai *Hat) NumberToClue(n int) int {
	if n < 4 { // value
		return n + 1
	}
	clue := n + int(game.ColorsTable[0]) + 1
	return clue // color + 5
}

func (ai *Hat) ExecuteAction(action *game.Action) *game.Action {
	info := &ai.PlayerInfo
	me := info.CurrentPosition
	rec := ai.MyRecord
	rec.ResetMemory()
	for i := 1; i < len(info.PlayerCards); i++ {
		pos := (me + i) % len(info.PlayerCards)
		if pos == me {
			panic("Bad position")
		}
		if action.ActionType == game.TypeActionInformationColor {
			action.Value += 5
		}
		ai.PlayerRecords[pos].ThinkOutOfTurn(pos, me, action, len(info.PlayerCards))
		if action.ActionType == game.TypeActionInformationColor {
			action.Value -= 5
		}
	}
	if action.ActionType == game.TypeActionDiscard && info.BlueTokens == game.MaxBlueTokens {
		panic("Magic")
	}

	if action.ActionType == game.TypeActionInformationColor || action.ActionType == game.TypeActionInformationValue {
		return action
	} else {
		return action
	}
}

func (ai *Hat) InitializeFuturePrediction(rec *HatPlayerRecord, penalty int) {
	info := &ai.PlayerInfo
	rec.FutureProgress = map[game.CardColor]game.CardValue{}
	for color, card := range info.TableCards {
		rec.FutureProgress[color] = card.Value
	}
	rec.FutureDeckSize = info.DeckSize
	rec.FutureHints = info.BlueTokens - penalty
	rec.FutureDiscarded = []game.ColorValue{}
}

func (ai *Hat) InitializeFutureClue(rec *HatPlayerRecord) {
	rec.CardsDrawn = 0
	rec.WillBePlayed = []game.ColorValue{}
	rec.HintChanges = []int{}
}

func (ai *Hat) FinalizeFutureAction(rec *HatPlayerRecord, action, i, me int) {
	if i == me {
		panic("me is i")
	}
	if action < 8 {
		rec.CardsDrawn++
		info := &ai.PlayerInfo
		cards := info.PlayerCards[i]
		if action < 4 {
			card := &cards[action]
			card.CheckVisible()
			rec.WillBePlayed = append(rec.WillBePlayed,
				game.ColorValue{Color: card.Color, Value: card.Value})
			if cards[action].Value == 5 {
				rec.HintChanges = append(rec.HintChanges, 1)
			}
		} else {
			rec.HintChanges = append(rec.HintChanges, 2)
			card := &cards[action-4]
			card.CheckVisible()
			rec.FutureDiscarded = append(rec.FutureDiscarded,
				game.ColorValue{Color: card.Color, Value: card.Value})
		}
	} else {
		rec.HintChanges = append(rec.HintChanges, -1)
	}
}

func (ai *Hat) FinalizeFutureClue(rec *HatPlayerRecord) {
	for _, p := range rec.WillBePlayed {
		rec.FutureProgress[p.Color]++
	}
	if result := rec.FutureDeckSize - rec.CardsDrawn; result > 0 {
		rec.FutureDeckSize = result
	} else {
		rec.FutureDeckSize = 0
	}
	ai.CountHints(rec)
}

func (ai *Hat) CountHints(rec *HatPlayerRecord) {
	rec.MinFutureHints = rec.FutureHints
	rec.MaxFutureHints = rec.FutureHints
	for i := len(rec.HintChanges) - 1; i >= 0; i-- {
		if i == 2 {
			if rec.FutureHints == 8 {
				rec.MaxFutureHints = 9
			} else {
				rec.FutureHints++
				if rec.FutureHints > rec.MaxFutureHints {
					rec.MaxFutureHints = rec.FutureHints
				}
			}
		} else {
			rec.FutureHints += i
			if rec.FutureHints > rec.MaxFutureHints {
				rec.MaxFutureHints = rec.FutureHints
			}
			if rec.FutureHints < rec.MinFutureHints {
				rec.MinFutureHints = rec.FutureHints
			}
			if rec.FutureHints < 0 {
				rec.FutureHints = 1
			}
			if rec.FutureHints > 8 {
				rec.FutureHints = 8
			}
		}
	}
}

func (ai *Hat) GetAction() *game.Action {
	info := &ai.PlayerInfo
	me := info.CurrentPosition
	rec := ai.MyRecord
	n := len(info.PlayerCards)
	for i := 0; i < n; i++ {
		ai.InterpretClue(i)
	}

	if rec.ImClued {
		myAction := ai.NumberToAction(rec.ClueValue)
		if myAction != nil {
			if myAction.ActionType == game.TypeActionDiscard && info.BlueTokens > 1 &&
				rec.LastCluedAnyClue == me &&
				len(ai.GetPlays(info.PlayerCards[(me+1)%n], ai.GetProgress(info))) > 0 {

			} else if myAction.ActionType == game.TypeActionPlaying ||
				(myAction.ActionType == game.TypeActionDiscard && info.BlueTokens < game.MaxBlueTokens) {
				return ai.ExecuteAction(myAction)
			}
		}
	} else {
		if rec.LastCluedAnyClue != (me-1+n)%n {
			panic("Bad last clued any clue")
		}
	}

	if info.BlueTokens == 0 {
		return ai.ExecuteAction(game.NewAction(game.TypeActionDiscard, me, 3))
	}

	if rec.LastCluedAnyClue == (me-1+n)%n {
		rec.LastCluedAnyClue = me
	}
	ai.InitializeFuturePrediction(rec, 1)
	ai.InitializeFutureClue(rec)
	i := rec.LastClued
	for _, x := range rec.NextPlayerActions {
		ai.FinalizeFutureAction(rec, x, i, me)
		i = (i - 1 + n) % n
	}
	ai.FinalizeFutureClue(rec)

	firstTarget := (rec.LastClued + 1) % n
	for _, targetPair := range rec.LaterClues {
		lastTarget, value := targetPair[0], targetPair[1]
		ai.InitializeFutureClue(rec)
		i = lastTarget
		for i != firstTarget {
			x := ai.StandardPlay(info.PlayerCards[i], i, rec.WillBePlayed,
				rec.FutureProgress, rec.FutureDeckSize, rec.FutureHints)
			ai.FinalizeFutureAction(rec, x, i, me)
			value = (value - x + 9) % 9
			i = (i - 1 + n) % n
		}
		ai.FinalizeFutureAction(rec, value, i, me)
		ai.FinalizeFutureClue(rec)
		firstTarget = (lastTarget + 1) % n
	}

	clueNumber := 0
	ai.InitializeFutureClue(rec)
	target := (me - 1 + n) % n
	i = target
	for i != (rec.LastCluedAnyClue+1)%n {
		x := ai.StandardPlay(info.PlayerCards[i], i, rec.WillBePlayed,
			rec.FutureProgress, rec.FutureDeckSize, rec.FutureHints)
		clueNumber = (clueNumber + x) % 9
		ai.FinalizeFutureAction(rec, x, i, me)
		i = (i - 1 + n) % n
	}

	ai.CountHints(rec)
	if i == me {
		panic("Magic")
	}
	x := ai.ModifiedPlay(info.PlayerCards[i], me, i, rec.WillBePlayed,
		rec.FutureProgress, rec.FutureDeckSize, rec.FutureHints)
	clueNumber = (clueNumber + x) % 9
	clue := ai.NumberToClue(clueNumber)
	var myAction *game.Action
	if clueNumber < 4 {
		myAction = game.NewAction(game.TypeActionInformationValue, target, clue)
	} else {
		myAction = game.NewAction(game.TypeActionInformationColor, target, clue-5)
	}
	rec.LastCluedAnyClue = target
	return ai.ExecuteAction(myAction)
}
