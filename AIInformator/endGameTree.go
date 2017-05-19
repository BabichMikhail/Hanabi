package informator

/*import (
	"fmt"

	game "github.com/BabichMikhail/Hanabi/game"
)

type ResultGameState struct {
	Probability float64
	Action      game.Action
	State       *game.GameState
}

type PreviewNode struct {
	ResultStates []ResultGameState
}

type PreviewTree struct {
	Nodes  map[game.Action]PreviewNode
	Min    int
	Max    int
	Medium float64
	Count  int
}

func (info *Informator) A() {

}

func (info *Informator) PreviewAction(playerInfo *game.PlayerGameInfo) error {
	if playerInfo.DeckSize > 0 {
		panic("Not implemented")
	}

	currentState := info.getCurrentState().Copy()
	fmt.Println(currentState.RedTokens)

	pos := currentState.CurrentPosition
	informationColorNodes := PreviewNode{}
	informationValueNodes := PreviewNode{}
	for i := 0; i < len(currentState.PlayerStates); i++ {
		if i == pos {
			continue
		}

		for j := 0; j < len(currentState.PlayerStates[i].PlayerCards); j++ {
			infoColorAction := game.NewAction(game.TypeActionInformationColor, i, int(currentState.PlayerStates[i].PlayerCards[j].Color))
			infoColorState := currentState.Copy()
			infoColorState.ApplyAction(infoColorAction)
			resultInfoColor := ResultGameState{
				Probability: 0.0,
				Action:      infoColorAction,
				State:       infoColorState,
			}
			informationColorNodes.ResultStates = append(informationColorNodes.ResultStates, resultInfoColor)

			infoValueAction := game.NewAction(game.TypeActionInformationValue, i, int(currentState.PlayerStates[i].PlayerCards[j].Value))
			infoValueState := currentState.Copy()
			infoValueState.ApplyAction(infoValueAction)
			resultInfoValue := ResultGameState{
				Probability: 0.0,
				Action:      infoValueAction,
				State:       infoValueState,
			}
			informationValueNodes.ResultStates = append(informationValueNodes.ResultStates, resultInfoValue)
		}
	}
	return nil
}*/
