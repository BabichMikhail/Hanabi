<script type="text/javascript">
    $(function() {
        Game = new gameHandler()

        Game.noneColor = {{ $.NoneColor }}
        Game.noneValue = {{ $.NoneValue }}
        Game.maxRedTokens = {{ $.MaxRedTokens }}
        Game.maxBlueTokens = {{ $.MaxRedTokens }}
        Game.actionTypes = {
            infoColor: {{ index $.ActionTypes "infoColor" }},
            infoValue: {{ index $.ActionTypes "infoValue" }},
            discard: {{ index $.ActionTypes "discard" }},
            play: {{ index $.ActionTypes "play" }},
        }

        {{ range $idx, $url := $.CardUrls }}
        Game.AddCardUrl({{ $url.Url }}, {{ $url.Color }}, {{ $url.Value }})
        {{ end }}

        Game.players = {
            {{ range $idx, $player := $.Players }}
            {{ $player.Id }}: {{ $player.NickName }},
            {{ end }}
        }

        Game.nickNames = [
            {{ range $idx, $nickName := $.NickNames }}
            {{ $nickName }},
            {{ end }}
        ]

        Game.colors = {
            {{ range $colorName, $color := $.TableColors }}
            {{ $colorName }}: {{ $color }},
            {{ end }}
        }

        console.log(Game)
        Game.LoadGameInfo()
        Game.CheckStep()
    })
</script>
