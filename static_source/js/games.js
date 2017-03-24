function gameHandler() {
    this.showBasicCards = [
        1, 1, 1, 1, 1,
    ]

    this.Init = function () {
        var Cards = []
        $.ajax({
            type: "GET",
            url: "/api/games/cards",
            data: {}
        }).done(function(data) {
            console.log(data)
            Cards = data
        }).fail(function(data) {
            console.log("init game fail")
        })
        return Cards
    }

    this.PlayCard = function (cardPosition) {
        $.ajax({
            type: "POST",
            url: "/api/games/action/play",
            data: {
                game_id: Game.id,
                card_position: cardPosition
            }
        }).done(function(data) {
            console.log(data)
            location.reload()
        }).fail(function(data) {
            console.log("fail play card #" + cardPosition)
        })
    }

    this.DiscardCard = function (cardPosition) {
        $.ajax({
            type: "POST",
            url: "/api/games/action/discard",
            data: {
                game_id: Game.id,
                card_position: cardPosition,
            },
        }).done(function(data) {
            console.log(data)
            location.reload()
        }).fail(function(data) {
            console.log("fail discard card #" + cardPosition)
        })
    }

    this.InfoCardValue = function(playerPosition, cardValue) {
        $.ajax({
            type: "POST",
            url: "/api/games/action/info/value",
            data: {
                game_id: Game.id,
                player_position: playerPosition,
                card_value: cardValue,
            },
        }).done(function(data) {
            console.log(data)
            location.reload()
        }).fail(function(data) {
            console.log("fail info about card value")
        })
    }

    this.InfoCardColor = function(playerPosition, cardColor) {
        $.ajax({
            type: "POST",
            url: "/api/games/action/info/color",
            data: {
                game_id: Game.id,
                player_position: playerPosition,
                card_color: cardColor,
            },
        }).done(function(data) {
            console.log(data)
            location.reload()
        }).fail(function(data) {
            console.log("fail info about card color")
        })
    }

    this.GetCardUrlByCard = function(card) {
        let color = card.knownColor ? card.color : Game.noneColor
        let value = card.knownValue ? card.value : Game.noneValue
        return Game.cardUrls[color][value]
    }

    this.GetCardUrlByCardIgnoreKnown = function(card) {
        return Game.cardUrls[card.color][card.value]
    }

    this.cardUrls = []

    this.AddCardUrl = function(url, color, value) {
        if (typeof Game.cardUrls[color] == 'undefined') {
            Game.cardUrls[color] = []
        }
        Game.cardUrls[color][value] = url
    }

    this.LoadGameInfo = function() {
        $.ajax({
            type: "GET",
            url: "/api/games/info",
            data: {
                game_id: Game.id,
            },
        }).done(function(data) {
            let gameInfo = data.data
            let playersHtml = []
            for (let i = 0; i < gameInfo.player_count; ++i) {
                playersHtml[i] = `<ul style="margin:0px">` +
                    (gameInfo.current_position == i ? `<i class="fa fa-hourglass fa-1" aria-hidden="true"></i>` : ``) +
                    `<li class="list-inline" style="display:inline-block; margin:1px">` + Game.players[i] + `</li>` +
                    `<li class="list-inline" style="display:inline-block; margin:1px">` +
                        `<button class="btn btn-info btn-sm" onclick="">info</button>` +
                    `</li>` +
                `</ul>`
                playersHtml[i] += `<ul id="basic-cards-` + i + `" style="margin:0px">`
                for (let j = 0; j < gameInfo.player_cards.length; ++j) {
                    let cards = gameInfo.player_cards[i]
                    playersHtml[i] = `<ul style="margin:0px">` +
                        (i == gameInfo.current_position ? `<i class="fa fa-hourglass fa-1" aria-hidden="true"></i>` : ``) +
                        `<li class="list-inline" style="display:inline-block; margin:1px">` + Game.nickNames[i] + `</li>` +
                        (gameInfo.pos != i ? `<li class="list-inline" style="display:inline-block; margin:1px">` +
                            `<button class="btn btn-info btn-sm" onClick="Game.ChangeCardsVisible(` + i + `)">info</button>` +
                        `</li>` : ``) +
                    `</ul>`

                    let cardsBasicHtml = ``
                    let actionsHtml = [ ``, ``, ``, ``, `` ]
                    for (let j = 0; j < cards.length; ++j) {
                        if (gameInfo.my_turn && gameInfo.pos == i) {
                            actionsHtml[j] = `<div style="font-size:14px"><a href="javascript:Game.PlayCard(` + j + `)">` +
                                `<i class="fa fa-play"> play</i>` +
                            `</a></div>` +
                            `<div style="font-size:14px"><a href="javascript:Game.DiscardCard(` + j + `)">` +
                                `<i class="fa fa-shopping-basket"> discard</i>` +
                            `</a></div>`
                        } else if (gameInfo.my_turn && gameInfo.pos != i) {
                            actionsHtml[j] = `<div style="font-size:14px"><a href="javascript:Game.InfoCardValue(` + j + `, ` + cards[j].value + `)">` +
                                `<i class="fa fa-info"> value</i>` +
                            `</a></div>` +
                            `<div style="font-size:14px"><a href="javascript:Game.InfoCardColor(` + j + `, ` + cards[j].color + `)">` +
                                `<i class="fa fa-info"> color</i>` +
                            `</a></div>`
                        }
                        cardsBasicHtml += `<li class="list-inline" style="display: inline-block; margin:1px">` +
                            `<img class="my-card" src="` + Game.GetCardUrlByCardIgnoreKnown(cards[j]) + `">` + actionsHtml[j] +
                        `</li>`
                    }
                    playersHtml[i] +=
                        `<ul id="basic-cards-` + i + `" ` + (Game.showBasicCards[i] ? `` : `class="invisible"`) + ` style="margin:0px">` +
                            cardsBasicHtml +
                        `</ul>`

                    if (gameInfo.pos != i) {
                        let cardsInfo = gameInfo.player_cards_info[i]
                        let cardsAdditionalHtml = ``
                        for (let j = 0; j < cards.length; ++j) {
                            cardsAdditionalHtml += `<li class="list-inline" style="display: inline-block; margin:1px">` +
                                `<img class="my-card" src="` + Game.GetCardUrlByCard(cardsInfo[j]) + `">` +  actionsHtml[j] +
                            `</li>`
                        }
                        playersHtml[i] +=
                            `<ul id="additional-cards-` + i + `" ` + (!Game.showBasicCards[i] ? `` : `class="invisible"`) + ` style="margin:0px">` +
                                cardsAdditionalHtml +
                            `</ul>`
                    }
                }
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
                    <img src="` + Game.GetCardUrlByCardIgnoreKnown(gameInfo.table_cards[Game.colors[tableColors[i]]]) + `" class="my-card">
                </li>`
            }

            let tableHtml = `<ul style="margin:0px">` + tableCardsHtml + `</ul>
                <ul>
                    <li class="list-inline" style="display: inline-block">
                        <div class="col-md-12" style="padding: 0px; height: 120px">
                            <img src="/static/img/deck.png" style="position:relative; width: 76px; height: 106px">
                            <img src="/static/img/number_` + Math.floor(gameInfo.deck_size / 10) + `.png" style="position:absolute; left:6px; top:30px; width: 30px">
                            <img src="/static/img/number_` + (gameInfo.deck_size % 10) + `.png" style="position:absolute; left:31px; top:30px; width: 30px">
                         </div>
                    </li>
                    <li class="list-inline" style="display: inline-block">
                        <img src="/static/img/token_red.png" class="game-token">
                        <img src="/static/img/number_` + gameInfo.red_tokens + `.png" style="position: relative; top: -60px; width: 30px">
                    </li>
                    <li class="list-inline" style="display: inline-block">
                        <img src="/static/img/token_blue.png" class="game-token">
                        <img src="/static/img/number_` + gameInfo.blue_tokens + `.png" style="position: relative; top: -60px; width: 30px">
                    </li>
                </ul>`

            let html = ""
            let count = gameInfo.player_cards.length
            let offset = gameInfo.pos
            if (count == 2) {
                html += `
                    <div class="col-md-12 game-player" id="player-1">` + playersHtml[(offset + 1) % count] + `</div>
                    <div class="col-md-12 game-table">` + tableHtml + `</div>
                    <div class="col-md-12 game-player" id="player-0">` + playersHtml[(offset + 0) % count] + `</div>`
            } else if (count == 3) {
                html += `
                    <div class="col-md-6 game-player" id="player-1">` + playersHtml[(offset + 1) % count] + `</div>
                    <div class="col-md-6 game-player" id="player-2">` + playersHtml[(offset + 2) % count] + `</div>
                    <div class="col-md-12 game-table">` + tableHtml + `</div>
                    <div class="col-md-12 game-player" id="player-0">` + playersHtml[(offset + 0) % count] + `</div>`
            } else if (count == 4) {
                html += `
                    <div class="col-md-12 game-player" id="player-2">` + playersHtml[(offset + 2) % count] + `</div>
                    <div class="col-md-4 game-player" id="player-1">` + playersHtml[(offset + 1) % count] + `</div>
                    <div class="col-md-4 game-table">` + tableHtml + `</div>
                    <div class="col-md-4 game-player" id="player-3">` + playersHtml[(offset + 3) % count] + `</div>
                    <div class="col-md-12 game-player" id="player-0">` + playersHtml[(offset + 0) % count] + `</div>`
            } else if (count == 5) {
                html += `
                    <div class="col-md-2"></div>
                    <div class="col-md-4 game-player" id="player-2">` + playersHtml[(offset + 2) % count] + `</div>
                    <div class="col-md-4 game-player" id="player-3">` + playersHtml[(offset + 3) % count] + `</div>
                    <div class="col-md-2"></div>
                    <div class="col-md-4 game-player" id="player-1">` + playersHtml[(offset + 1) % count] + `</div>
                    <div class="col-md-4 game-table" style="">` + tableHtml + `</div>
                    <div class="col-md-4 game-player" id="player-4">` + playersHtml[(offset + 4) % count] + `</div>
                    <div class="col-md-12 game-player"id="player-0">` + playersHtml[(offset + 0) % count] + `</div>`
            }
            $("#game").html(html)
        }).fail(function(data) {
            console.log("fail load game info")
        })
    }

    this.CheckStep = function() {
        $.ajax({
            type: "GET",
            url: "/api/games/step",
            data: {
                game_id: Game.id,
            },
        }).done(function(data) {
            let meta = $("meta[name=step]")
            if (data.step != meta.attr("step")) {
                //location.reload()
                Game.LoadGameInfo
            }
            setTimeout(Game.CheckStep, 10000)
        })
    }

    this.ChangeCardsVisible = function(pos) {
        let basicCard = $('ul[id="basic-cards-' + pos + '"]')
        let additionalCard = $('ul[id="additional-cards-' + pos + '"]')
        if (basicCard.hasClass('invisible')) {
            Game.showBasicCards[pos] = 1
            basicCard.removeClass('invisible')
            additionalCard.addClass('invisible')
        } else {
            Game.showBasicCards[pos] = 0
            basicCard.addClass('invisible')
            additionalCard.removeClass('invisible')
        }
    }

    this.Cards = this.Init()

    let meta = $("meta[name=step]")
    this.CurrentStep = +meta.attr("step")
    this.id = window.location.pathname.match(/\d+/)[0]

    return this
}
