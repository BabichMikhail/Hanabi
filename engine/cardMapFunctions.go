package engine

import "github.com/astaxie/beego"

func GetCardValue(value CardValue) string {
	return map[CardValue]string{
		NoneValue: "Unknown Value",
		One:       "1",
		Two:       "2",
		Three:     "3",
		Four:      "4",
		Five:      "5",
	}[value]
}

func GetCardColor(color CardColor) string {
	return map[CardColor]string{
		NoneColor: "Unknown Color",
		Red:       "Red",
		Blue:      "Blue",
		Green:     "Green",
		Gold:      "Gold",
		Black:     "Black",
	}[color]
}

func RegisterFunction() {
	beego.AddFuncMap("cardValue", GetCardValue)
	beego.AddFuncMap("cardColor", GetCardColor)
}
