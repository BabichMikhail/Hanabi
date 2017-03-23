<script type="text/javascript">
    $(function() {
        View = new viewHandler()
        View.stepCount = {{ len $.Actions }}

        View.actions = [
        {{ range $idx, $action := $.Actions }}
            View.NewAction({{ $action.ActionType }}, {{ $action.PlayerPosition }}, {{ $action.Value }}),
        {{ end }}
        ]
        deck = [
        {{ $initState := $.InitState }}
        {{ range $idx, $card := $initState.Deck }}
            View.NewCard({{ $card.Color }}, {{ $card.KnownColor }}, {{ $card.Value }}, {{ $card.KnownValue }}),
        {{ end }}
        ]
        View.currentStep = 0

        let playerStates = [
        {{ range $idx, $state := $initState.PlayerStates }}{
            playerId: {{ $state.PlayerId }},
            playerPosition: {{ $state.PlayerPosition }},
            playerCards: [
            {{ range $i, $card := $state.PlayerCards }}
                View.NewCard({{ $card.Color }}, {{ $card.KnownColor }}, {{ $card.Value }}, {{ $card.KnownValue }}),
            {{ end }}
            ],
        },{{ end }}
        ]
        View.games[0] = {
            blueTokens: {{ $initState.BlueTokens }},
            redTokens: {{ $initState.RedTokens }},
            currentPosition: {{ $initState.CurrentPosition }},
            usedCards: [],
            playerStates: playerStates,
            tableCards: {
                {{ range $colorName, $color := $.TableColors }}
                {{ $card := index $initState.TableCards $color }}
                {{ $card.Color }}: View.NewCard({{ $card.Color }}, {{ $card.KnownColor }}, {{ $card.Value }}, {{ $card.KnownValue }}),
                {{ end }}
            },
            deck: deck,

        }
        View.colors = {
            {{ range $colorName, $color := $.TableColors }}
            {{ $colorName }}: {{ $color }},
            {{ end }}
        }

        View.players = {
            {{ range $idx, $player := $.Players }}
            {{ $player.Id }}: {{ $player.NickName }},
            {{ end }}
        }
        View.noneColor = {{ $.NoneColor }}
        View.noneValue = {{ $.NoneValue }}
        View.maxRedTokens = {{ $.MaxRedTokens}}
        View.maxBlueTokens = {{ $.MaxRedTokens}}
        View.actionTypes = {
            infoValue: {{ index $.ActionTypes "infoColor" }},
            infoColor: {{ index $.ActionTypes "infoValue" }},
            discard: {{ index $.ActionTypes "discard" }},
            play: {{ index $.ActionTypes "play" }},
        }
        {{ range $idx, $url := $.CardUrls }}
        View.AddCardUrl({{ $url.Url }}, {{ $url.Color }}, {{ $url.Value }})
        {{ end }}
        View.MakeTable()
        View.MakeActions()
        console.log(View.games[0])
        console.log(View)
    });

</script>
