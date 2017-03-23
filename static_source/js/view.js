function viewHandler() {
    this.games = []
    this.stepCount = 0
    this.currentStep = -1;
    this.showBasicCards = [
        1, 1, 1, 1, 1,
    ]

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
                if (game.playerStates[action.pos].playerCards[i].color == action.value) {
                    game.playerStates[action.pos].playerCards[i].knownColor = true
                }
            }
            --game.blueTokens
        } else if (action.type == View.actionTypes["discard"]) {
            let oldCard = game.playerStates[action.pos].playerCards[action.value]
            oldCard.knownValue = true
            oldCard.KnownColor = true
            game.playerStates[action.pos].playerCards.splice(action.value, 1)
            if (game.deck.length > 0) {
                let newCard = game.deck.shift()
                game.playerStates[action.pos].playerCards.push(newCard)
            }
            ++game.blueTokens
            game.usedCards.push(oldCard)
        } else if (action.type == View.actionTypes["play"]) {
            let oldCard = game.playerStates[action.pos].playerCards[action.value]
            oldCard.knownValue = true
            oldCard.KnownColor = true
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
                cardsBasicHtml += `<li class="list-inline" style="display: inline-block; margin:1px">
                    <img class="my-card" src="` + View.GetCardUrlByCardIgnoreKnown(cards[j]) + `">
                </li>`
                cardsAdditionalHtml += `<li class="list-inline" style="display: inline-block; margin:1px">
                    <img class="my-card" src="` + View.GetCardUrlByCard(cards[j]) + `">
                </li>`
            }
            htmlPlayers[i] = `<ul style="margin:0px">
                <li class="list-inline" style="display:inline-block; margin:1px">` + View.players[game.playerStates[i].playerId] + `</li>
                <li class="list-inline" style="display:inline-block; margin:1px">
                    <button class="btn btn-info btn-sm" onClick="View.ChangeCardsVisible(` + i + `)">info</button>
                </li>
            </ul>`
            htmlPlayers[i] +=
                `<ul id="basic-cards-` + i + `" ` + (View.showBasicCards[i] ? `` : `class="invisible"`) + ` style="margin:0px">` +
                    cardsBasicHtml +
                `</ul>` +
                `<ul id="additional-cards-` + i + `" ` + (!View.showBasicCards[i] ? `` : `class="invisible"`) + ` style="margin:0px">` +
                    cardsAdditionalHtml +
                `</ul>`
        }

        tableColors = [
            "blue",
            "green",
            "orange",
            "red",
            "yellow",
        ]

        let tableCardsHtml = ``
        for (let i in tableColors) {
            tableCardsHtml += `<li id="table-` + tableColors[i] + `-cards" class="list-inline" style="display: inline-block; margin:1px">
                <img src="` + View.GetCardUrlByCardIgnoreKnown(game.tableCards[View.colors[tableColors[i]]]) + `" class="my-card">
            </li>`
        }

        let tableHtml = `<ul style="margin:0px">` + tableCardsHtml + `</ul>
            <ul>
                <li class="list-inline" style="display: inline-block">
                    <div class="col-md-12" style="padding: 0px; height: 120px">
                        <img src="/static/img/deck.png" style="position:relative; width: 76px; height: 106px">
                        <img src="/static/img/number_` + Math.floor(game.deck.length / 10) + `.png" style="position:absolute; left:6px; top:30px; width: 30px">
                        <img src="/static/img/number_` + game.deck.length % 10 + `.png" style="position:absolute; left:31px; top:30px; width: 30px">
                    </div>
                </li>
                <li class="list-inline" style="display: inline-block">
                    <img src="/static/img/token_red.png" class="game-token">
                    <img src="/static/img/number_` + (View.maxRedTokens - game.redTokens) + `.png" style="position: relative; top: -60px; width: 30px">
                </li>
                <li class="list-inline" style="display: inline-block">
                    <img src="/static/img/token_blue.png" class="game-token">
                    <img src="/static/img/number_` + game.blueTokens + `.png" style="position: relative; top: -60px; width: 30px">
                </li>
            </ul>`

        let html = ""
        if (game.playerStates.length == 2) {
            html += `
                <div class="col-md-12 game-player" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-12 game-table">` + tableHtml + `</div>
                <div class="col-md-12 game-player" id="player-0">` + htmlPlayers[0] + `</div>`
        } else if (game.playerStates.length == 3) {
            html += `
                <div class="col-md-6 game-player" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-6 game-player" id="player-2">` + htmlPlayers[2] + `</div>
                <div class="col-md-12 game-table">` + tableHtml + `</div>
                <div class="col-md-12 game-player" id="player-0">` + htmlPlayers[0] + `</div>`
        } else if (game.playerStates.length == 4) {
            html += `
                <div class="col-md-12 game-player" id="player-2">` + htmlPlayers[2] + `</div>
                <div class="col-md-4 game-player" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-4 game-table">` + tableHtml + `</div>
                <div class="col-md-4 game-player" id="player-3">` + htmlPlayers[3] + `</div>
                <div class="col-md-12 game-player" id="player-0">` + htmlPlayers[0] + `</div>`
        } else if (game.playerStates.length == 5) {
            html += `
                <div class="col-md-2"></div>
                <div class="col-md-4 game-player" id="player-2">` + htmlPlayers[2] + `</div>
                <div class="col-md-4 game-player" id="player-3">` + htmlPlayers[3] + `</div>
                <div class="col-md-2"></div>
                <div class="col-md-4 game-player" id="player-1">` + htmlPlayers[1] + `</div>
                <div class="col-md-4 game-table" style="">` + tableHtml + `</div>
                <div class="col-md-4 game-player" id="player-4">` + htmlPlayers[4] + `</div>
                <div class="col-md-12 game-player"id="player-0">` + htmlPlayers[0] + `</div>`
        }

        html += `<div class="col-md-12" style="text-align:center"></ul>`
        for (let i = 0; i < game.usedCards.length; ++i) {
            html += `<li class="list-inline" style="display: inline-block; margin:1px">` +
                `<img src="` + View.GetCardUrlByCardIgnoreKnown(game.usedCards[i]) + `" class="my-card">` +
            `</li>`
        }
        html += `</ul></div>`

        $("#game").html(html)
    }

    this.ChangeCardsVisible = function(pos) {
        let basicCard = $('ul[id="basic-cards-' + pos + '"]')
        let additionalCard = $('ul[id="additional-cards-' + pos + '"]')
        if (basicCard.hasClass('invisible')) {
            View.showBasicCards[pos] = 1
            basicCard.removeClass('invisible')
            additionalCard.addClass('invisible')
        } else {
            View.showBasicCards[pos] = 0
            basicCard.addClass('invisible')
            additionalCard.removeClass('invisible')
        }
    }

    return this
}
