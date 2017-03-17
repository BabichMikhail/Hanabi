package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"

	ai "github.com/BabichMikhail/Hanabi/AI"
	stats "github.com/BabichMikhail/Hanabi/statistic"
)

type Stat struct {
	Id          int       `orm:"auto" json:"id"`
	PlayerCount int       `orm:"column(player_count)" json:"player_count"`
	AITypesJSON string    `orm:"column(ai_types)" json:"-"`
	AITypes     []int     `orm:"-" json:"ai_types"`
	AINames     []string  `orm:"-" json:"ai_names"`
	GameCount   int       `orm:"column(count)" json:"count"`
	Points      float64   `orm:"column(points);null" json:"points"`
	Ready       time.Time `orm:"column(ready_at);null" json:"-"`
	ReadyStr    string    `orm:"-" json:"ready_at"`
	Created     time.Time `orm:"column(created_at)" json:"-"`
	CreatedStr  string    `orm:"-" json:"created_at"`
	ExecTime    int       `orm:"-" json:"execution_time"`
}

func (stat *Stat) TableName() string {
	return "stats"
}

func NewStat(aiTypes []int, count int) {
	o := orm.NewOrm()
	b, err := json.Marshal(aiTypes)
	if err != nil {
		panic(err)
	}

	var statId int
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.InsertInto("stats", "ai_types", "count", "player_count", "created_at").
		Values("?", "?", "?", "CURRENT_TIMESTAMP").
		String()
	if res, err := o.Raw(sql, string(b), count, len(aiTypes)).Exec(); err == nil {
		id64, _ := res.LastInsertId()
		statId = int(id64)
	} else {
		panic(err)
	}
	StartStat(statId, aiTypes, count)
}

func StartStat(id int, aiTypes []int, count int) {
	stat := stats.RunGames(aiTypes, count)
	ReadyStat(id, stat.Medium)
}

func ReadyStat(id int, points float64) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Update("stats").
		Set("ready_at = CURRENT_TIMESTAMP", "points = "+strconv.FormatFloat(points, 'E', -1, 64)).
		Where("id = ?").
		String()
	_, err := o.Raw(sql, id).Exec()
	if err != nil {
		panic(err)
	}
}

func ReadStats() (stats []Stat) {
	stats = []Stat{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id", "player_count", "ai_types", "count", "points", "ready_at", "created_at").
		From("stats").
		Where("points IS NOT NULL")
	_, err := o.Raw(qb.String()).QueryRows(&stats)
	if err != nil {
		return []Stat{}
	}

	for idx, stat := range stats {
		err := json.Unmarshal([]byte(stat.AITypesJSON), &stats[idx].AITypes)
		if err != nil {
			return []Stat{}
		}

		stats[idx].AINames = make([]string, len(stats[idx].AITypes), cap(stats[idx].AITypes))
		for i, aiType := range stats[idx].AITypes {
			stats[idx].AINames[i] = ai.AINames[aiType]
		}
		stats[idx].ReadyStr = stat.Ready.Format("15:04:05 02.01.2006")
		stats[idx].CreatedStr = stat.Created.Format("15:04:05 02.01.2006")
		stats[idx].ExecTime = int((stat.Ready.UnixNano() - stat.Created.UnixNano()) / 1000000000)
	}

	return
}
