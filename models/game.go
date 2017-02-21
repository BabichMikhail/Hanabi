package models

import "github.com/astaxie/beego/orm"

func GetGamePlayers(gameIds []int) map[int]([]LobbyPlayer) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	playersMap := map[int]([]LobbyPlayer){}
	if len(gameIds) > 0 {
		inString := IntSliceToString(gameIds)
		qb.Select("p.user_id", "u.nick_name", "p.game_id").
			From("players p").
			InnerJoin("user u").
			On("u.id = p.user_id").
			Where("p.game_id").In(inString)
		sql := qb.String()
		var splayers []struct {
			UserId   int    `orm:"column(user_id)"`
			NickName string `orm:"column(nick_name)"`
			GameId   int    `orm:"column(game_id)"`
		}
		o.Raw(sql).QueryRows(&splayers)

		for _, v := range splayers {
			playersMap[v.GameId] = append(playersMap[v.GameId], LobbyPlayer{
				Id:       v.UserId,
				NickName: v.NickName,
			})
		}
	}
	return playersMap
}
