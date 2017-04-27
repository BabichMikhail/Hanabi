package models

import (
	"github.com/BabichMikhail/Hanabi/game"
	"github.com/astaxie/beego/orm"
)

type QData struct {
	Id           int     `orm:"auto"`
	PlayersCount int     `orm:"column(players_count)"`
	BlueTokens   int     `orm:"column(blue_tokens)"`
	RedTokens    int     `orm:"column(red_tokens)"`
	DeckSize     int     `orm:"column(deck_size)"`
	StepLeft     int     `orm:"column(step_left)"`
	InfoCount    int     `orm:"column(info_count)"`
	Points       int     `orm:"column(points)"`
	Count        int     `orm:"column(count);default(0)"`
	Result       float64 `orm:"column(result);default(25)"`
}

func (q *QData) TableName() string {
	return "qdata"
}

func (q *QData) TableIndex() [][]string {
	return [][]string{
		[]string{"result"},
	}
}

func (q *QData) TableUnique() [][]string {
	return [][]string{
		[]string{"players_count", "blue_tokens", "red_tokens", "deck_size", "step_left", "info_count", "points"},
	}
}

func getStepLeft(deckSize, step, maxStep int) int {
	stepLeft := 0
	if deckSize == 0 {
		stepLeft = maxStep - step + 1
	}
	return stepLeft
}

func getInfoCount(playerCardsInfo [][]game.Card) int {
	infoCount := 0
	for _, cards := range playerCardsInfo {
		for _, card := range cards {
			if card.KnownColor {
				infoCount++
			}
			if card.KnownValue {
				infoCount++
			}
		}
	}
	return infoCount
}

func qRead(info *game.PlayerGameInfo) (*QData, error) {
	stepLeft := getStepLeft(info.DeckSize, info.Step, info.MaxStep)
	infoCount := getInfoCount(info.PlayerCardsInfo)

	qdata := QData{
		PlayersCount: len(info.PlayerCards),
		BlueTokens:   info.BlueTokens,
		RedTokens:    info.RedTokens,
		DeckSize:     info.DeckSize,
		StepLeft:     stepLeft,
		InfoCount:    infoCount,
		Points:       info.GetPoints(),
		Count:        0,
		Result:       25,
	}
	o := orm.NewOrm()
	_, _, err := o.ReadOrCreate(&qdata, "players_count", "blue_tokens", "red_tokens", "deck_size", "step_left", "info_count", "points")
	return &qdata, err
}

func QRead(info *game.PlayerGameInfo) float64 {
	qdata, err := qRead(info)
	if err == nil {
		return qdata.Result
	}
	return 25.0
}

func qUpdate(qdata *QData, result float64) {
	sql := `UPDATE qdata SET count = count + 1, result = (count*result + ?) / (count + 1)`
	o := orm.NewOrm()
	_, err := o.Raw(sql, result).Exec()
	if err != nil {
		panic(err)
	}
}

func QUpdate(gameState *game.GameState, result float64) {
	info := gameState.GetPlayerGameInfoByPos(gameState.CurrentPosition, game.InfoTypeUsually)
	qdata, err := qRead(&info)
	if err != nil {
		return
	}
	qUpdate(qdata, result)
}
