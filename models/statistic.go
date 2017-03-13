package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
)

// @todo add Stat to registration of models
type Stat struct {
	Id          int       `orm:"auto"`
	PlayerCount int       `orm:"column(player_count)"`
	AITypesJSON string    `orm:"column(ai_types)"`
	AITypes     []int     `orm:"-"`
	GameCount   int       `orm:"column(count)"`
	Points      float64   `orm:"column(points);null"`
	Ready       time.Time `orm:"column(ready_at);null"`
	Created     time.Time `orm:"column(created_at)"`
}

func (stat *Stat) TableName() string {
	return "stats"
}

func NewStat(aiTypes []int, points float64, count int) int {
	o := orm.NewOrm()
	b, err := json.Marshal(aiTypes)
	if err != nil {
		panic(err)
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.InsertInto("stats", "ai_types", "count", "created").
		Values(string(b), "?", "CURRENT_TIMESTAMP").
		String()
	if res, err := o.Raw(sql, count).Exec(); err == nil {
		id64, _ := res.LastInsertId()
		return int(id64)
	} else {
		panic(err)
	}
}

func ReadyStat(id int, points float64) {
	o := orm.NewOrm()
	var stat Stat
	_, err := o.QueryTable(stat).Filter("id", id).Update(orm.Params{
		"points": points,
		"ready":  "CURRENT_TIMESTAMP",
	})
	if err != nil {
		panic(err)
	}
}

func ReadStats() (stats []Stat) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("player_count", "ai_types", "count", "points", "ready_at", "created_at").
		From("stats").
		Where("points IS NOT NULL")
	_, err := o.Raw(qb.String()).QueryRows(&stats)
	if err != nil {
		return []Stat{}
	}

	for _, stat := range stats {
		err := json.Unmarshal([]byte(stat.AITypesJSON), &stat.AITypes)
		if err != nil {
			return []Stat{}
		}
	}

	return
}
