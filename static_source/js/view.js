function viewHandler() {
    this.games = []
    this.stepCount = 0
    this.currentStep = -1;

    this.NewCard = function(color, knownColor, value, knownValue) {
        return {
            color: color,
            knownColor: knownColor,
            value: value,
            knownValue: knownValue,
        }
    }

    this.NewAction = function(type, playerPosition, value) {
        return {
            type: type,
            pos: playerPosition,
            value: value,
        }
    }

    this.MakeGame = function(step) {
        if (typeof View.games[step - 1] == 'undefined') {
            this.MakeGame(step - 1)
        }
        // @todo
    }

    this.cardUrls = []

    this.AddCardUrl = function(url, color, value) {
        if (typeof View.cardUrls[color] == 'undefined') {
            View.cardUrls[color] = []
        }
        View.cardUrls[color][value] = url
    }

    this.GetCardUrlByCard = function(card) {
        let color = card.knownColor ? card.color : View.noneColor
        let value = card.knownValue ? card.value : View.noneValue
        return View.cardUrls[color][value]
    }

    this.GetCardUrlByCardIgnoreKnown = function(card) {
        return View.cardUrls[card.color][card.value]
    }

    this.MakeTable = function() {
        if (typeof View.games[View.currentStep] == 'undefined') {c
            View.MakeGame(View.currentStep)
        }
        game = View.games[View.currentStep]

        let htmlTable = ``
        let htmlPlayers = []
        for (let i = 0; i < game.playerStates.length; ++i) {
            cards = game.playerStates[i].playerCards[i]
            cardsBasicHtml = ``
            cardsAdditionalHtml = ``
            for (let j = 0; j < cards.length; ++j) {
                cardsBasicHtml += `<li class="list-inline" style="display: inline-block">
                    <img class="myCard" src="` + View.GetCardUrlByCardIgnoreKnown(cards[j]) + `">
                </li>`
                cardsAdditionalHtml += `<li class="list-inline" style="display: inline-block">
                    <img class="myCard" src="` + View.GetCardUrlByCard(cards[j]) + `">
                </li>`
            }
            htmlPlayers[i] = `<ul id="basic-cards-` + i + `" name="basic-cards">
                ` + cardsBasicHtml + `</ul>`
            htmlPlayers[i] += `<ul id="additional-cards-` + i + `" name="additional-cards" class="invisible">
                ` + cardsAdditionalHtml + `</ul>`
            htmlPlayers[i] = `<div id="player-pos-` + i + `" class="col-md-12">
                <div id="player-cards-` + i + `" class="col-md-12">` + htmlPlayers[i] + `
                </div>
            </div>`
        }

        tmpHtml = ``
        for (let i = 0; i < htmlPlayers.length; ++i) {
            tmpHtml += htmlPlayers[i]
        }
        $("#users").html(tmpHtml)
    }

    return this
}
