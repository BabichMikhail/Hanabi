package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/BabichMikhail/Hanabi/game"
	_ "github.com/BabichMikhail/Hanabi/routers"
	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestCardColors
func TestCardColors(t *testing.T) {
	Convey("Test game colors", t, func() {
		Convey("Colors should be 5 and NoneColor", func() {
			So(len(game.Colors), ShouldEqual, 6)
		})

		Convey("All colors must be different", func() {
			colors := map[game.CardColor]bool{}
			for _, color := range game.Colors {
				colors[color] = true
			}
			So(len(colors), ShouldEqual, 6)
		})

		Convey("Colors contain NoneColor", func() {
			contain := false
			for _, color := range game.Colors {
				contain = contain || (color == game.NoneColor)
			}
			So(contain, ShouldBeTrue)
		})
	})
}

// TestCardValues
func TestCardValues(t *testing.T) {
	Convey("Test game colors", t, func() {
		Convey("Values should be 5 and NoneValue", func() {
			So(len(game.Values), ShouldEqual, 6)
		})

		Convey("All values must be different", func() {
			values := map[game.CardValue]bool{}
			for _, value := range game.Values {
				values[value] = true
			}
			So(len(values), ShouldEqual, 6)
		})

		Convey("Values contain NoneValue", func() {
			contain := false
			for _, value := range game.Values {
				contain = contain || (value == game.NoneValue)
			}
			So(contain, ShouldBeTrue)
		})
	})
}

//TestHashValues
func TestHashValues(t *testing.T) {
	Convey("Test hashvalues", t, func() {
		values := map[game.HashValue]int{}
		Convey("Hashvalue must be 0 <= val < 25", func() {
			ok := true
			for _, color := range game.Colors {
				for _, value := range game.Values {
					val := game.HashColorValue(color, value)
					values[val]++
					if color != game.NoneColor && value != game.NoneValue {
						ok = ok && 0 <= val && val < 25
					}
				}
			}
			So(ok, ShouldBeTrue)

		})

		Convey("Hashvalue must be unique", func() {
			ok := true
			for _, count := range values {
				if count > 1 {
					ok = false
				}
			}
			So(ok, ShouldBeTrue)
		})
	})
}
