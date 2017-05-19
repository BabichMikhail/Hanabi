package test

/*import (
	"fmt"
	"testing"

	ai "github.com/BabichMikhail/Hanabi/AI"
	info "github.com/BabichMikhail/Hanabi/AIInformator"
	"github.com/BabichMikhail/Hanabi/game"
	_ "github.com/BabichMikhail/Hanabi/routers"
	. "github.com/smartystreets/goconvey/convey"
)

func PrintPlayerInfo(playerInfo *game.PlayerGameInfo, name string) {
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println(name)
	fmt.Println("Table:")
	for _, color := range game.Colors {
		card := playerInfo.TableCards[color]
		fmt.Print(card.Value, " ", game.GetCardColor(color), "\t|\t")
	}
	fmt.Println()
	fmt.Println("Cards:")
	for i := 0; i < len(playerInfo.PlayerCards); i++ {
		fmt.Print("Player #", i, ":\t")
		for j := 0; j < len(playerInfo.PlayerCards[i]); j++ {
			card := &playerInfo.PlayerCards[i][j]
			fmt.Print(card.Value, " ", game.GetCardColor(card.Color), "\t|\t")
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------------------------------------")
}

// TestSetAvailableInformation
func TestsetAvailableInformation(t *testing.T) {
	Convey("Test ai.SetAvailableInformation", t, func() {
		Convey("TESst 1", func() {
			ids := []int{
				1, 2, 3, 4, 5,
			}
			g := game.NewGame(ids, 1002)
			informator := info.NewInformator(g.CurrentState, g.Actions)

			fmt.Println("-------------------------------------------------------------------------")
			fmt.Println("Initial GameState")
			fmt.Println("Table:")
			currentState := g.CurrentState
			playerStates := currentState.PlayerStates
			for _, color := range game.Colors {
				card := currentState.TableCards[color]
				fmt.Print(card.Value, " ", game.GetCardColor(color), "\t|\t")
			}
			fmt.Println()
			fmt.Println("Cards:")
			for i := 0; i < len(playerStates); i++ {
				fmt.Print("Player #", i, ":\t")
				for j := 0; j < len(playerStates[i].PlayerCards); j++ {
					card := &playerStates[i].PlayerCards[j]
					fmt.Print(card.Value, " ", game.GetCardColor(card.Color), "\t|\t")
				}
				fmt.Println()
			}
			fmt.Println("-------------------------------------------------------------------------")

			AI := informator.NextAI(ai.Type_AIRandom).(*ai.AIRandom)
			AI.setAvailableInformation()
			playerInfo := &AI.PlayerInfo
			PrintPlayerInfo(playerInfo, "First PlayerInfo")
			pos := playerInfo.CurrentPosition
			for i := 0; i < len(playerInfo.PlayerCards[pos]); i++ {
				t.Log(fmt.Sprintln(playerInfo.PlayerCards[pos][i]))
			}
			So(pos, ShouldEqual, 0)

			informator.ApplyAction(game.NewAction(game.TypeActionInformationValue, 1, 1))
			AI = informator.NextAI(ai.Type_AIRandom).(*ai.AIRandom)
			playerInfo = &AI.PlayerInfo
			AI.setAvailableInformation()
			PrintPlayerInfo(playerInfo, "Second PlayerInfo")

			fmt.Println()

		})
	})
}
*/
