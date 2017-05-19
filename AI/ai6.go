package ai

import (
	"fmt"

	"github.com/BabichMikhail/Hanabi/game"
)

type AI6 struct {
	BaseAI
}

func NewAI6(baseAI *BaseAI) *AI6 {
	ai := new(AI6)
	ai.BaseAI = *baseAI
	return ai
}

/*
5 Players
Cluer = CurrentPlayer
NP = CurrentPlayer + 3
NNP = CurrentPlayer + 4

20 InfoTypes:
1) NP Value 1
2) NP Value 2
3) NP Value 3
4) NP Value 4
5) NP Value 5
6) NP Color 1
7) NP Color 2
8) NP Color 3
9) NP Color 4
10) NP Color 5
11) NNP Value 1
12) NNP Value 2
13) NNP Value 3
14) NNP Value 4
15) NNP Value 5
16) NNP Color 1
17) NNP Color 2
18) NNP Color 3
19) NNP Color 4
20) NNP Color 5

Drop 8 most useless and use 12 for play cards on pos for NP & NNP:
1)	0	-1
2)	1	-1
3)	2	-1
4)	3	-1
5)	0	0
6)	1	0
7)	2	0
8)	3	0
9)	0	1
10)	1	1
11)	2	1
12)	3	1

For this play actions clue to NP about most left card. (13-0, 14-1, 15-2, 16-3)
13) 0   2
14) 1	2
15) 2 	2
16) 3	2

For this play actions clue to NNP about most left card. (17-0, 18-1, 15-2, 16-3)
17) 0 	3
18) 1	3
19) 2 	3
20) 3	3

If no useful clue that do useless clue (from 8/20) or discard card
*/

func (ai *AI6) DefaultInfoActions(info *game.PlayerGameInfo) ([]*game.Action, []*game.Action) {
	values := game.Values[1:]
	colors := game.ColorsTable
	myPos := info.CurrentPosition
	NP := ai.NextPlayer(myPos, 3)
	NNP := ai.NextPlayer(myPos, 4)
	actions := []*game.Action{}
	uselessActions := []*game.Action{}
	for _, value := range values {
		actions = append(actions, game.NewAction(game.TypeActionInformationValue, NP, int(value)))
	}
	for _, color := range colors {
		actions = append(actions, game.NewAction(game.TypeActionInformationColor, NP, int(color)))
	}
	for _, value := range values {
		actions = append(actions, game.NewAction(game.TypeActionInformationValue, NNP, int(value)))
	}
	for _, color := range colors {
		actions = append(actions, game.NewAction(game.TypeActionInformationColor, NNP, int(color)))
	}

	iterationCount := 0
	for len(actions) > 12 {
		iterationCount++
		for idx, action := range actions {
			isUseful := false
			for pos, card := range info.PlayerCards[action.PlayerPosition] {
				playerPos := action.PlayerPosition
				cardInfo := info.PlayerCards[playerPos][pos]
				if action.ActionType == game.TypeActionInformationColor {
					if card.Color == game.CardColor(action.Value) && !cardInfo.KnownColor {
						isUseful = true
						break
					}
				} else {
					if card.Value == game.CardValue(action.Value) && !cardInfo.KnownValue {
						isUseful = true
						break
					}
				}
			}
			if !isUseful {
				//uselessActions = append(uselessActions, actions[idx])
				actions = append(actions[:idx], actions[idx+1:]...)
				break
			}
		}
		if iterationCount == 20 {
			uselessActions = append(uselessActions, actions[12:]...)
			actions = actions[:12]
		}
	}
	return actions, uselessActions
}

func (ai *BaseAI) NextPlayer(pos, offset int) int {
	playersCount := len(ai.PlayerInfo.PlayerCards)
	return (pos + offset + playersCount) % playersCount
}

func (ai *AI6) DecodeClue1(info *game.PlayerGameInfo, action *game.Action) int {
	myPos := info.CurrentPosition
	cluerPos := ai.NextPlayer(myPos, -1)
	NP := myPos
	NPP := ai.NextPlayer(myPos, 1)
	cluerInfo := ai.Informator.GetPlayerState(info.Step - 1)
	cluerInfo.PlayerCards[cluerPos] = cluerInfo.PlayerCardsInfo[cluerPos]
	actions, _ := ai.DefaultInfoActions(&cluerInfo)
	clueIdx := -1
	for idx, clueAction := range actions {
		if clueAction.Equal(action) {
			clueIdx = idx
		}
	}

	if clueIdx == -1 {
		if action.PlayerPosition == NPP || action.PlayerPosition == NP {
			cluerInfo = *info
			cards := cluerInfo.PlayerCards[action.PlayerPosition]
			fmt.Println(cards, NP)
			for i := 0; i < len(cards); i++ {
				card := &cards[i]
				if action.ActionType == game.TypeActionInformationColor {
					if card.KnownColor && card.Color == game.CardColor(action.Value) {
						//fmt.Println(i)
						//panic("ABC")
						return i
					}
				} else {
					if card.KnownValue && card.Value == game.CardValue(action.Value) {
						//fmt.Println(i)
						//panic("ABC")
						return i
					}
				}
			}
		}
		return -1
	}

	cardPos := map[int]int{
		0:  0,
		1:  1,
		2:  2,
		3:  3,
		4:  0,
		5:  1,
		6:  2,
		7:  3,
		8:  0,
		9:  1,
		10: 2,
		11: 3,
		/*12: 0,
		13: 1,
		14: 2,
		15: 3,
		16: 0,
		17: 1,
		18: 2,
		19: 3,*/
	}[clueIdx]
	return cardPos
}

func (ai *AI6) DecodeClue2(info *game.PlayerGameInfo, action *game.Action) int {
	myPos := info.CurrentPosition
	cluerPos := ai.NextPlayer(myPos, -2)
	NP := ai.NextPlayer(myPos, -1)
	NPP := myPos
	cluerInfo := ai.Informator.GetPlayerState(info.Step - 2)
	cluerInfo.PlayerCards[cluerPos] = cluerInfo.PlayerCardsInfo[cluerPos]
	actions, _ := ai.DefaultInfoActions(&cluerInfo)
	clueIdx := -1
	for idx, clueAction := range actions {
		if clueAction.Equal(action) {
			clueIdx = idx
		}
	}

	if clueIdx == -1 {
		if action.PlayerPosition == NPP {
			return 3
		} else if action.PlayerPosition == NP {
			return 2
		} else {
			return -1
		}
	}
	cardPos := map[int]int{
		0:  -1,
		1:  -1,
		2:  -1,
		3:  -1,
		4:  0,
		5:  0,
		6:  0,
		7:  0,
		8:  1,
		9:  1,
		10: 1,
		11: 1,
		/*12: 2,
		13: 2,
		14: 2,
		15: 2,
		16: 3,
		17: 3,
		18: 3,
		19: 3,*/
	}[clueIdx]
	return cardPos
}

func (ai *AI6) GetAction() *game.Action {
	info := &ai.PlayerInfo
	fmt.Println("Step:", info.Step, info.RedTokens)
	info.SetProbabilities(false, false)
	myPos := info.CurrentPosition

	for idx, card := range info.PlayerCards[myPos] {
		if card.KnownColor && card.KnownValue && info.TableCards[card.Color].Value+1 == card.Value {
			fmt.Println(game.NewAction(game.TypeActionPlaying, myPos, idx))
			return game.NewAction(game.TypeActionPlaying, myPos, idx)
		}
	}

	if info.DeckSize < 3 {
		if info.RedTokens < 1 {
			for idx, card := range info.PlayerCards[myPos] {
				for colorValue, prob := range card.ProbabilityCard {
					color, value := game.ColorValueByHashColorValue(colorValue)
					if prob > 0.6 && info.TableCards[color].Value+1 == value {
						fmt.Println(game.NewAction(game.TypeActionPlaying, myPos, idx))
						return game.NewAction(game.TypeActionPlaying, myPos, idx)
					}
				}
			}
		}

		if info.RedTokens < 2 {
			for idx, card := range info.PlayerCards[myPos] {
				for colorValue, prob := range card.ProbabilityCard {
					color, value := game.ColorValueByHashColorValue(colorValue)
					if prob > 0.8 && info.TableCards[color].Value+1 == value {
						fmt.Println(game.NewAction(game.TypeActionPlaying, myPos, idx))
						return game.NewAction(game.TypeActionPlaying, myPos, idx)
					}
				}
			}
		}
	}

	if len(info.PlayerCards) != 5 {
		panic("Not implemented")
	}

	if len(ai.History) > 0 {
		action := ai.History[len(ai.History)-1]
		isInformationAction := action.ActionType == game.TypeActionInformationColor || action.ActionType == game.TypeActionInformationValue
		if isInformationAction {
			cardPos := ai.DecodeClue1(info, &action)
			fmt.Println("Decode1:", cardPos)
			if cardPos != -1 && ai.isCardPlayable(info.PlayerCards[myPos][cardPos]) {
				fmt.Println(myPos, game.NewAction(game.TypeActionPlaying, myPos, cardPos))
				return game.NewAction(game.TypeActionPlaying, myPos, cardPos)
			}
		}
	}

	if len(ai.History) > 1 {
		action1 := ai.History[len(ai.History)-2]
		action2 := ai.History[len(ai.History)-1]
		isInformationAction := action1.ActionType == game.TypeActionInformationColor || action1.ActionType == game.TypeActionInformationValue
		isActionPlay := action2.ActionType == game.TypeActionPlaying
		if isInformationAction && isActionPlay {
			cardPos := ai.DecodeClue2(info, &action1)
			fmt.Println("Decode2:", cardPos)
			if cardPos != -1 && ai.isCardPlayable(info.PlayerCards[myPos][cardPos]) {
				fmt.Println(myPos, game.NewAction(game.TypeActionPlaying, myPos, cardPos))
				return game.NewAction(game.TypeActionPlaying, myPos, cardPos)
			}
		}
	}

	NP := ai.NextPlayer(myPos, 1)
	NNP := ai.NextPlayer(myPos, 2)
	clueVariants := [][]int{}
	var clueActions, uselessActions []*game.Action
	if info.BlueTokens > 0 {
		clueActions, uselessActions = ai.DefaultInfoActions(info)
		for cardPosFirst, cardFirst := range info.PlayerCards[NP] {
			if !ai.isCardPlayable(cardFirst) {
				continue
			}
			copyInfo := *info.Copy()
			preview, err := copyInfo.PreviewActionInformationColor(NP, cardFirst.Color)
			if err != nil {
				panic(err)
			}
			if len(preview.Results) == 0 {
				//fmt.Println(len(preview.Results))
				panic("Bad preview result")
			}
			result := preview.Results[0]
			newInfo := result.Info
			if newInfo.IsGameOver() {
				break
			}
			preview, err = newInfo.PreviewActionPlaying(cardPosFirst)
			if err != nil {
				panic(err)
			}
			result = preview.Results[0]
			for cardPosSecond, cardSecond := range copyInfo.PlayerCards[NNP] {
				info = info.Copy()
				if !cardSecond.KnownColor || !cardSecond.KnownValue {
					panic("DSADSA")
				}
				ai.PlayerInfo = *result.Info
				if ai.isCardPlayable(cardSecond) {
					clueVariants = append(clueVariants, []int{cardPosFirst, cardPosSecond})
				} else {
					clueVariants = append(clueVariants, []int{cardPosFirst, -1})
				}
				ai.PlayerInfo = *info
			}
		}
	}

	if len(clueVariants) > 0 {
		results := [][]int{
			[]int{0, -1, 0},
			[]int{1, -1, 1},
			[]int{2, -1, 2},
			[]int{3, -1, 3},
			[]int{0, 0, 4},
			[]int{1, 0, 5},
			[]int{2, 0, 6},
			[]int{3, 0, 7},
			[]int{0, 1, 8},
			[]int{1, 1, 9},
			[]int{2, 1, 10},
			[]int{3, 1, 11},
		}

		if len(clueActions) != 12 {
			panic("Bad len(clueActions)")
		}

		topActionIdx := -1
		maxUsefulCount := -1
		f := func(result []int, v []int) {
			g := func(action *game.Action, pos int) int {
				cards := info.PlayerCards[pos]
				playerInfo := ai.Informator.GetPlayerState(info.Step)
				playerInfo.PlayerCards[pos] = playerInfo.PlayerCardsInfo[pos]
				playerInfo.SetProbabilities(false, false)
				cardsInfo := playerInfo.PlayerCardsInfo[pos]
				usefulCount := 0
				for i := 0; i < len(cards); i++ {
					if action.ActionType == game.TypeActionInformationColor {
						if cards[i].Color == game.CardColor(action.Value) && !cardsInfo[i].KnownColor {
							if ai.isCardMayBeUsefull(cards[i]) {
								usefulCount++
							}
						}
					} else {
						if cards[i].Value == game.CardValue(action.Value) && !cardsInfo[i].KnownValue {
							if ai.isCardMayBeUsefull(cards[i]) {
								usefulCount++
							}
						}
					}
				}
				return usefulCount
			}

			if result[0] == v[0] && result[1] == v[1] {
				action := clueActions[result[2]]
				pos := action.PlayerPosition
				if usefulCount := g(action, pos); usefulCount > maxUsefulCount {
					maxUsefulCount = usefulCount
					topActionIdx = result[2]
				}
			} else if v[1] >= 2 {
				var pos int
				if v[1] == 2 {
					pos = NP
				} else {
					pos = NNP
				}
				cards := info.PlayerCards[pos]
				card := &cards[v[0]]
				actions := []*game.Action{
					game.NewAction(game.TypeActionInformationColor, pos, int(card.Color)),
					game.NewAction(game.TypeActionInformationValue, pos, int(card.Value)),
				}

				for _, action := range actions {
					actionOK := true
					for i := 0; i < len(cards); i++ {
						if action.ActionType == game.TypeActionInformationColor {
							if cards[i].Color == card.Color && i < v[0] {
								actionOK = false
								break
							}
						} else {
							if cards[i].Value == card.Value && i < v[0] {
								actionOK = false
								break
							}
						}
					}
					if !actionOK {
						continue
					}
					if usefulCount := g(action, pos); usefulCount > maxUsefulCount {
						results = append(results, []int{v[0], v[1], len(results)})
						clueActions = append(clueActions, action)
						maxUsefulCount = usefulCount
						topActionIdx = len(results) - 1
					}
				}
			}
		}

		for j := 0; j < len(clueVariants); j++ {
			v := clueVariants[j]
			if v[1] != -1 {
				for i := 0; i < len(results); i++ {
					f(results[i], v)
				}
			}
		}

		if topActionIdx != -1 {
			fmt.Println(clueActions[topActionIdx], results[topActionIdx][0], results[topActionIdx][1])
			return clueActions[topActionIdx]
		}

		for j := 0; j < len(clueVariants); j++ {
			for i := 0; i < len(results); i++ {
				f(results[i], clueVariants[j])
			}
		}

		var a, b int
		for i := 0; i < len(clueVariants); i++ {
			v := clueVariants[i]
			if v[0] == results[topActionIdx][0] && results[topActionIdx][1] == v[1] {
				a, b = v[0], v[1]
			}
		}

		fmt.Println(clueActions[topActionIdx], a, b)
		return clueActions[topActionIdx]
	}

	if info.BlueTokens > 0 && len(uselessActions) > 0 {
		fmt.Println("//", uselessActions[0], -1)
		//return uselessActions[0]
	}

	if info.BlueTokens > 5 {
		for i := 0; i < len(info.PlayerCards[NP]); i++ {
			copyInfo := info.Copy()
			copyInfo.PlayerCards[NP] = copyInfo.PlayerCardsInfo[NP]
			copyInfo.InfoIsSetted = false
			copyInfo.SetProbabilities(false, false)
			card := &copyInfo.PlayerCards[NP][i]

			var previewColor, previewValue *game.ResultPreviewPlayerInformations
			var err error
			countValue := -1
			copyInfo1 := info.Copy()
			if !card.KnownValue {
				//actionValue = game.NewAction(game.TypeActionInformationValue, NP, int(card.Value))
				previewValue, err = copyInfo1.PreviewActionInformationValue(NP, card.Value)
				if err != nil {
					panic(err)
				}
				for k := 0; k < len(previewValue.Results); k++ {
					previewInfo := previewValue.Results[k].Info
					previewCard := previewInfo.PlayerCards[NP][k]
					ai.PlayerInfo = *previewValue.Results[k].Info
					isPlayable := ai.isCardPlayable(previewCard)
					ai.PlayerInfo = *copyInfo1
					if !isPlayable {
						countValue = 0
						for j := 0; j < len(info.PlayerCardsInfo[NP]); j++ {
							pcard := info.PlayerCards[NP][j]
							cardInfo := info.PlayerCardsInfo[NP][j]
							if pcard.Value == card.Value && !cardInfo.KnownValue {
								countValue++
							}
						}
					}
				}
			}

			countColor := -1
			copyInfo2 := info.Copy()
			if !card.KnownColor {
				//actionColor = game.NewAction(game.TypeActionInformationColor, NP, int(card.Color))
				previewColor, err = copyInfo2.PreviewActionInformationColor(NP, card.Color)
				if err != nil {
					panic(err)
				}
				for k := 0; k < len(previewColor.Results); k++ {
					previewInfo := previewColor.Results[k].Info
					previewCard := previewInfo.PlayerCards[NP][k]
					ai.PlayerInfo = *previewColor.Results[k].Info
					isPlayable := ai.isCardPlayable(previewCard)
					ai.PlayerInfo = *copyInfo2
					//fmt.Println(card)
					if !isPlayable {
						countColor = 0
						for j := 0; j < len(info.PlayerCardsInfo[NP]); j++ {
							pcard := info.PlayerCards[NP][j]
							cardInfo := info.PlayerCardsInfo[NP][j]
							if pcard.Color == card.Color && !cardInfo.KnownColor {
								countColor++
							}
						}
					}
				}
			}

			//fmt.Println("Counts:", countValue, countColor, card)
			if countValue == -1 && countColor == -1 {
				continue
			}

			if countValue >= countColor {
				//panic("ABC")
				//panic("QWEERTY2")
				fmt.Println(game.NewAction(game.TypeActionInformationValue, NP, int(card.Value)), -1, -1)
				return game.NewAction(game.TypeActionInformationValue, NP, int(card.Value))
			} else {
				//panic("ABC")
				//panic("QWEERTY1")
				fmt.Println(game.NewAction(game.TypeActionInformationColor, NP, int(card.Color)), -1, -1)
				return game.NewAction(game.TypeActionInformationColor, NP, int(card.Color))
			}
		}
	}

	if info.BlueTokens < 6 {
		for idx, card := range info.PlayerCards[myPos] {
			if !ai.isCardMayBeUsefull(card) {
				fmt.Println(game.NewAction(game.TypeActionDiscard, myPos, idx))
				return game.NewAction(game.TypeActionDiscard, myPos, idx)
			}
		}
		fmt.Println(game.NewAction(game.TypeActionDiscard, myPos, 0))
		return game.NewAction(game.TypeActionDiscard, myPos, 0)
	}

	if info.BlueTokens > 0 && len(uselessActions) > 0 {
		return uselessActions[0]
	}

	if info.RedTokens == 2 && info.BlueTokens < game.MaxBlueTokens {
		fmt.Println(game.NewAction(game.TypeActionDiscard, myPos, 0), "BadDiscard")
		return game.NewAction(game.TypeActionDiscard, myPos, 0)
	} else {
		if info.Step == 0 {
			fmt.Println(info)
			panic("ABCWQE")
		}
		fmt.Println(game.NewAction(game.TypeActionPlaying, myPos, 0), "BadPlay")
		return game.NewAction(game.TypeActionPlaying, myPos, 0)
	}
}
