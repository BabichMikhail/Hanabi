package game

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
		Yellow:    "Gold",
		Orange:    "Black",
	}[color]
}

func getCardUrlByValueAndColor(color CardColor, value CardValue) string {
	name := "/static/img/" + map[CardColor]string{
		NoneColor: "unknown",
		Red:       "red",
		Blue:      "blue",
		Green:     "green",
		Yellow:    "yellow",
		Orange:    "orange",
	}[color]

	name += "_" + map[CardValue]string{
		NoneValue: "unknown",
		One:       "one",
		Two:       "two",
		Three:     "three",
		Four:      "four",
		Five:      "five",
	}[value] + ".png"
	return name
}

func RegisterFunction() {
	beego.AddFuncMap("cardValue", GetCardValue)
	beego.AddFuncMap("cardColor", GetCardColor)
	beego.AddFuncMap("getCardUrl", getCardUrlByValueAndColor)
}
