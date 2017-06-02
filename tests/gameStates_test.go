package test

import (
	"fmt"
	"testing"

	"github.com/BabichMikhail/Hanabi/game"
	_ "github.com/BabichMikhail/Hanabi/routers"
	. "github.com/smartystreets/goconvey/convey"
)

// TestGame
func TestGame(t *testing.T) {
	Convey("Game Init", t, func() {
		pseudoIds := []int{1, 2, 3, 4, 5}
		g := game.NewGame(pseudoIds, game.Type_NormalGame)
		Convey("Game is created", func() {
			So(g, ShouldNotBeNil)
		})

		Convey("Check init state", func() {
			Convey("Init state have 3 red tokens", func() {
				So(g.InitState.RedTokens, ShouldEqual, game.MaxRedTokens)
			})

			Convey("Init state have 8 blue tokens", func() {
				So(g.InitState.BlueTokens, ShouldEqual, game.MaxBlueTokens)
			})

			Convey(fmt.Sprintf("Game must have %d players states", len(pseudoIds)), func() {
				So(len(g.InitState.PlayerStates), ShouldEqual, len(pseudoIds))
			})

			Convey("Step must be equal 0", func() {
				So(g.InitState.Step, ShouldEqual, 0)
			})
		})
	})
}
