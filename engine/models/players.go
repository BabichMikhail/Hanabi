package enginemodels

type Player struct {
	Id       int    `orm:"column(id)" json:"id"`
	NickName string `orm:"column(nick_name)" json:"nick_name"`
}
