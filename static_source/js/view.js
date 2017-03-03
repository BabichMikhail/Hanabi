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

    this.CalculateGame = function(step) {
        if (typeof View.games[step] != 'undefined') {
            return
        }
        if (step > 0 && typeof View.games[step - 1] == 'undefined') {
            this.CalculateGame(step - 1)
        }

        // game_i + action_i = game_(i + 1)
        let action = View.actions[step - 1]
        let game = jQuery.extend(true, {}, View.games[step - 1]);
        if (action.type == View.actionTypes["infoValue"]) {
            for (let i = 0; i < game.playerStates[action.pos].playerCards.length; ++i) {
                if (game.playerStates[action.pos].playerCards[i].value == action.value) {
                    game.playerStates[action.pos].playerCards[i].knownValue = true
                }
            }
            --game.blueTokens
        } else if (action.type == View.actionTypes["infoColor"]) {
            for (let i = 0; i < game.playerStates[action.pos].playerCards.length; ++i) {
                if (game.playerStates[action.pos].playerCards[i].value == action.value) {
                    game.playerStates[action.pos].playerCards[i].knownValue = true
                }
            }
            --game.blueTokens
        } else if (action.type == View.actionTypes["discard"]) {
            let oldCard = game.playerStates[action.pos].playerCards[action.value]
            game.playerStates[action.pos].playerCards.splice(action.value, 1)
            if (game.deck.length > 0) {
                let newCard = game.deck.shift()
                game.playerStates[action.pos].playerCards.push(newCard)
            }
            ++game.blueTokens
            game.usedCards.push(oldCard)
        } else if (action.type == View.actionTypes["play"]) {
            let oldCard = game.playerStates[action.pos].playerCards[action.value]
            game.playerStates[action.pos].playerCards.splice(action.value, 1)
            if (game.deck.length > 0) {
                let newCard = game.deck.shift()
                game.playerStates[action.pos].playerCards.push(newCard)
            }
            if (game.tableCards[oldCard.color].value + 1 == oldCard.value) {
                game.tableCards[oldCard.color] = oldCard
                if (oldCard.value == 5 && game.blueTokens < View.maxBlueTokens) {
                    ++game.blueTokens
                }
            } else {
                game.usedCards.push(oldCard)
                --game.redTokens
            }
        }
        View.currentStep = step
        View.games[step] = game
    }

    this.MakeNextGame = function() {
        if (View.currentStep >= View.actions.length) {
            return
        }
        View.currentStep++
        View.MakeTable()
    }

    this.MakePrevGame = function() {
        if (View.currentStep <= 0) {
            return
        }
        View.currentStep--
        View.MakeTable()
    }

    this.MakeGame = function(step) {
        if (step < 0 || step > View.actions.length) {
            return
        }
        View.currentStep = step
        View.MakeTable()
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
        if (typeof View.games[View.currentStep] == 'undefined') {
            View.CalculateGame(View.currentStep)
        }
        game = View.games[View.currentStep]

        let htmlTable = ``
        let htmlPlayers = []
        for (let i = 0; i < game.playerStates.length; ++i) {
            cards = game.playerStates[i].playerCards
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
            htmlPlayers[i] = View.players[game.playerStates[i].playerId] +
                `<ul id="basic-cards-` + i + `" name="basic-cards">` + cardsBasicHtml + `</ul>`
            htmlPlayers[i] += `<ul id="additional-cards-` + i + `" name="additional-cards" class="invisible">
                ` + cardsAdditionalHtml + `</ul>`
            htmlPlayers[i] = `<div id="player-pos-` + i + `" class="col-md-12">
                <div id="player-cards-` + i + `" class="col-md-12">` + htmlPlayers[i] + `
                </div>
            </div>`
        }

        let tableHtml = `<div id="table-pos" class="col-md-12" style="text-align: center">
            <div class="col-md-12">
                <ul>
                    <li id="table-blue-cards" class="list-inline" style="display: inline-block">
                        <img src="` + View.GetCardUrlByCardIgnoreKnown(game.tableCards[View.colors["blue"]]) + `" class="myCard">
                    </li>
                    <li id="table-green-cards" class="list-inline" style="display: inline-block">
                        <img src="` + View.GetCardUrlByCardIgnoreKnown(game.tableCards[View.colors["green"]]) + `" class="myCard">
                    </li>
                    <li id="table-orange-cards" class="list-inline" style="display: inline-block">
                        <img src="` + View.GetCardUrlByCardIgnoreKnown(game.tableCards[View.colors["orange"]]) + `" class="myCard">
                    </li>
                    <li id="table-red-cards" class="list-inline" style="display: inline-block">
                        <img src="` + View.GetCardUrlByCardIgnoreKnown(game.tableCards[View.colors["red"]]) + `" class="myCard">
                    </li>
                    <li id="table-yellow-cards" class="list-inline" style="display: inline-block">
                        <img src="` + View.GetCardUrlByCardIgnoreKnown(game.tableCards[View.colors["yellow"]]) + `" class="myCard">
                    </li>
                </ul>
            </div>
            <div class="col-md-12">
                <ul>
                    <li class="list-inline" style="display: inline-block">
                        <div class="col-md-12" style="padding: 0px" height="120px">
                            <img src="/static/img/deck.png" width="76" height="106" style="position:relative;">
                            <img src="/static/img/number_` + Math.floor(game.deck.length / 10) + `.png" width="30" style="position:absolute; left:6px; top:30px">
                            <img src="/static/img/number_` + game.deck.length % 10 + `.png" width="30" style="position:absolute; left:31px; top:30px">
                        </div>
                    </li>
                    <li class="list-inline" style="display: inline-block">
                        <img src="/static/img/token_red.png" style="position: relative; top: -30px" width="50" height="50">
                        <img src="/static/img/number_` + (View.maxRedTokens - game.redTokens) + `.png" style="position: relative; top: -30px" width="30">
                    </li>
                    <li class="list-inline" style="display: inline-block">
                        <img src="/static/img/token_blue.png" style="position: relative; top: -30px" width="50" height="50">
                        <img src="/static/img/number_` + game.blueTokens + `.png" style="position: relative; top: -30px" width="30">
                    </li>
                </ul>
            </div>
        </div>`

        let html = ""
        if (game.playerStates.length == 2) {
            html += `
                <div class="col-md-12" style="text-align:center" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-12" id="table">` + tableHtml + `</div>
                <div class="col-md-12" style="text-align:center" id="player-0">` + htmlPlayers[0] + `</div>`
        } else if (game.playerStates.length == 3) {
            html += `
                <div class="col-md-6" style="text-align:center" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-6" style="text-align:center" id="player-2">` + htmlPlayers[2] + `</div>
                <div class="col-md-12" id="table">` + tableHtml + `</div>
                <div class="col-md-12" style="text-align:center" id="player-0">` + htmlPlayers[0] + `</div>`
        } else if (game.playerStates.length == 4) {
            html += `
                <div class="col-md-12" style="text-align:center" id="player-2">` + htmlPlayers[2] + `</div>
                <div class="col-md-4" style="text-align:center" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-4" id="table">` + tableHtml + `</div>
                <div class="col-md-4" style="text-align:center" id="player-3">` + htmlPlayers[3] + `</div>
                <div class="col-md-12" style="text-align:center" id="player-0">` + htmlPlayers[0] + `</div>`
        } else if (game.playerStates.length == 5) {
            html += `
                <div class="col-md-2"></div>
                <div class="col-md-4" style="text-align:center" id="player-2">` + htmlPlayers[2] + `</div>
                <div class="col-md-4" style="text-align:center" id="player-3">` + htmlPlayers[3] + `</div>
                <div class="col-md-2"></div>
                <div class="col-md-4" style="text-align:center" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-4" id="table">` + tableHtml + `</div>
                <div class="col-md-4" style="text-align:center" id="player-4">` + htmlPlayers[4] + `</div>
                <div class="col-md-12" style="text-align:center" id="player-0">` + htmlPlayers[0] + `</div>`
        }
        $("#users").html(html)
    }

    return this
}
