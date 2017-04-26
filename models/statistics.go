package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"

	ai "github.com/BabichMikhail/Hanabi/AI"
	stats "github.com/BabichMikhail/Hanabi/statistics"
)

type Stat struct {
	Id          int       `orm:"auto" json:"id"`
	PlayerCount int       `orm:"column(player_count)" json:"player_count"`
	AITypesJSON string    `orm:"column(ai_types)" json:"-"`
	AITypes     []int     `orm:"-" json:"ai_types"`
	AINames     []string  `orm:"-" json:"ai_names"`
	GameCount   int       `orm:"column(count)" json:"count"`
	WinCount    int       `orm:"column(win_count);default(0)" json:"win_count"`
	Points      float64   `orm:"column(points);null" json:"points"`
	Dispersion  float64   `orm:"column(dispersion);null" json:"dispersion"`
	Kurtosis    float64   `orm:"column(kurtosis);null" json:"curtosis"`
	Asymmetry   float64   `orm:"column(asymmenty);null" json:"asymmetry"`
	Ready       time.Time `orm:"column(ready_at);null" json:"-"`
	IsMetadata  bool      `orm:"is_metadata" json:"is_metadata"`
	ReadyStr    string    `orm:"-" json:"ready_at"`
	Created     time.Time `orm:"column(created_at)" json:"-"`
	CreatedStr  string    `orm:"-" json:"created_at"`
	ExecTime    int       `orm:"-" json:"execution_time"`
	ReadyPart   int       `orm:"column(ready_part);default(0)" json:"-"`
	ReadyJSON   float64   `orm:"-" json:"ready_part"`
	IsReady     bool      `orm:"-" json:"is_ready"`
}

func (stat *Stat) TableName() string {
	return "stats"
}

func NewStat(aiTypes []int, count int, saveDistrInExcel bool) {
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
	StartStat(statId, aiTypes, count, saveDistrInExcel)
}

func updateReady(id int) func(*int, int) {
	updateFunc := func(readyCount *int, count int) {
		repeats := 0
		time.Sleep(time.Second * 10)
		var g func(*int, int, int)
		f := func(readyCount *int, lastReadyCount int, count int) {
			if *readyCount == count {
				return
			}

			if *readyCount == lastReadyCount {
				repeats++
			} else {
				repeats = 0
			}

			if repeats == 100 {
				return
			}
			lastReadyCount = *readyCount

			o := orm.NewOrm()
			qb, _ := orm.NewQueryBuilder("mysql")
			sql := qb.Update("stats").
				Set(
					"ready_part = " + strconv.Itoa((*readyCount)*10000/count),
				).
				Where("id = ?").
				String()

			_, err := o.Raw(sql, id).Exec()
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Second * 6)
			g(readyCount, lastReadyCount, count)
		}

		g = f
		f(readyCount, -1, count)
	}
	return updateFunc
}

func StartStat(id int, aiTypes []int, count int, saveDistrInExcel bool) {
	ids := make([]int, len(aiTypes), len(aiTypes))
	var err error
	for i := 0; i < len(aiTypes); i++ {
		ids[i], err = GetAIUserId(aiTypes[i], i)
		if err != nil {
			userIds, err := CreateAIUsers(aiTypes[i])
			ids[i] = userIds[i]
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	stat, games := stats.RunGames(aiTypes, ids, count, updateReady(id), QRead, QUpdate)
	ReadyStat(id, &stat, saveDistrInExcel)

	for _, game := range games {
		item, _ := NewGame(ids[0], len(aiTypes), StatusWait, true)
		gameId := item.Id
		for i := 1; i < len(ids); i++ {
			JoinGame(gameId, ids[i])
		}

		if err := UpdateInitGameState(gameId, game.InitState); err != nil {
			panic(err)
		}
		if err := UpdateGameState(gameId, game.CurrentState); err != nil {
			panic(err)
		}
		for _, action := range game.Actions {
			NewAction(gameId, &action)
		}
		SetGameFinishedStatus(gameId)
	}
}

func ReadyStat(id int, stat *stats.Stat, saveDistrInExcel bool) {
	err := stat.SaveToFile("stat_metadata", strconv.Itoa(id), saveDistrInExcel)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	isMetadata := 0
	if err == nil {
		isMetadata = 1
	}

	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Update("stats").
		Set(
			"win_count = "+strconv.Itoa(stat.Wins),
			"ready_at = CURRENT_TIMESTAMP",
			"points = "+strconv.FormatFloat(stat.Medium, 'E', -1, 64),
			"dispersion = "+strconv.FormatFloat(stat.Disp, 'E', -1, 64),
			"kurtosis = "+strconv.FormatFloat(stat.Kurt, 'E', -1, 64),
			"asymmenty = "+strconv.FormatFloat(stat.Asym, 'E', -1, 64),
			"is_metadata = "+strconv.Itoa(isMetadata),
		).
		Where("id = ?").
		String()
	_, err = o.Raw(sql, id).Exec()
	if err != nil {
		panic(err)
	}

}

func ReadStats() (stats []Stat) {
	stats = []Stat{}
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id", "player_count", "ai_types", "count", "points", "ready_at", "created_at", "ready_part").
		From("stats").
		OrderBy("created_at DESC")
	_, err := o.Raw(qb.String()).QueryRows(&stats)
	if err != nil {
		return []Stat{}
	}

	for idx, _ := range stats {
		stat := &stats[idx]
		err := json.Unmarshal([]byte(stat.AITypesJSON), &stat.AITypes)
		if err != nil {
			return []Stat{}
		}

		stat.AINames = make([]string, len(stat.AITypes), cap(stat.AITypes))
		for i, aiType := range stat.AITypes {
			stat.AINames[i] = ai.AINames[aiType]
		}

		stat.CreatedStr = stat.Created.Format("15:04:05 02.01.2006")

		stat.IsReady = stat.Points > 0
		if stat.IsReady {
			stat.ReadyStr = stat.Ready.Format("15:04:05 02.01.2006")
			stat.ExecTime = int((stat.Ready.UnixNano() - stat.Created.UnixNano()) / 1000000000)
		} else {
			stat.ReadyJSON = float64(stat.ReadyPart) / 100
		}
	}

	return
}

func DeleteStat(id int) error {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Delete().From("stats").Where("id = ?")
	_, err := o.Raw(qb.String(), id).Exec()
	return err
}
