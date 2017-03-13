function gameHandler() {
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
            alert("INIT GAME FAIL")
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
            alert("FAIL PLAY CARD #" + cardPosition)
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
            alert("FAIL DISCARD CARD #" + cardPosition)
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
            alert("FAIL INFO ABOUT CARD VALUE")
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
            alert("FAIL INFO ABOUT CARD COLOR")
        })
    }

    this.LoadGameInfo = function() {
        $.ajax({
            type: "GET",
            url: "/api/games/info",
            data: {
                game_id: Game.id,
            },
        }).done(function(data) {
            Game.playerCount = data.player_count
            let count = Game.playerCount
            Game.myPos = data.player_position
            let offset = Game.myPos
            let html = ""
            if (Game.playerCount == 2) {
                html += `
                    <div class="col-md-12" style="text-align:center" id="player-` + ((offset + 1) % count) + `"></div>
                    <div class="col-md-12" id="table"></div>
                    <div class="col-md-12" style="text-align:center" id="player-` + offset + `"></div>`
                $("#game-table").append(html)
                $("#table").append($("#table-pos").detach())
                $("#player-0").append($("#player-pos-0").detach())
                $("#player-1").append($("#player-pos-1").detach())
            } else if (Game.playerCount == 3) {
                html += `
                    <div class="col-md-6" style="text-align:center" id="player-` + ((offset + 1) % count) + `"></div>
                    <div class="col-md-6" style="text-align:center" id="player-` + ((offset + 2) % count) + `"></div>
                    <div class="col-md-12" id="table"></div>
                    <div class="col-md-12" style="text-align:center" id="player-` + offset + `"></div>`
                $("#game-table").append(html)
                $("#table").append($("#table-pos").detach())
                $("#player-0").append($("#player-pos-0").detach())
                $("#player-1").append($("#player-pos-1").detach())
                $("#player-2").append($("#player-pos-2").detach())
            } else if (Game.playerCount == 4) {
                html += `
                    <div class="col-md-12" style="text-align:center" id="player-` + ((offset + 2) % count) + `"></div>
                    <div class="col-md-4" style="text-align:center" id="player-` + ((offset + 1) % count) + `"></div>
                    <div class="col-md-4" id="table"></div>
                    <div class="col-md-4" style="text-align:center" id="player-` + ((offset + 3) % count) + `"></div>
                    <div class="col-md-12" style="text-align:center" id="player-` + offset + `"></div>`
                $("#game-table").append(html)
                $("#table").append($("#table-pos").detach())
                $("#player-0").append($("#player-pos-0").detach())
                $("#player-1").append($("#player-pos-1").detach())
                $("#player-2").append($("#player-pos-2").detach())
                $("#player-3").append($("#player-pos-3").detach())
            } else if (Game.playerCount == 5) {
                html += `
                    <div class="col-md-2"></div>
                    <div class="col-md-4" style="text-align:center" id="player-` + ((offset + 2) % count) + `"></div>
                    <div class="col-md-4" style="text-align:center" id="player-` + ((offset + 3) % count) + `"></div>
                    <div class="col-md-2"></div>
                    <div class="col-md-4" style="text-align:center" id="player-` + ((offset + 1) % count) + `"></div>
                    <div class="col-md-4" id="table"></div>
                    <div class="col-md-4" style="text-align:center" id="player-` + ((offset + 4) % count) + `"></div>
                    <div class="col-md-12" style="text-align:center" id="player-` + offset + `"></div>`
                $("#game-table").append(html)
                $("#table").append($("#table-pos").detach())
                $("#player-0").append($("#player-pos-0").detach())
                $("#player-1").append($("#player-pos-1").detach())
                $("#player-2").append($("#player-pos-2").detach())
                $("#player-3").append($("#player-pos-3").detach())
            }
        }).fail(function(data) {
            alert("FAIL LOAD GAME INFO")
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
                location.reload()
            }
            setTimeout(Game.CheckStep, 10000)
        })
    }

    this.ChangeCardsVisible = function() {
        let basicCards = $('ul[name=basic-cards]')
        let additionalCards = $('ul[name=additional-cards]')
        if (basicCards.hasClass('invisible')) {
            basicCards.removeClass('invisible')
            additionalCards.addClass('invisible')
        } else {
            basicCards.addClass('invisible')
            additionalCards.removeClass('invisible')
        }
    }

    this.Cards = this.Init()

    let meta = $("meta[name=step]")
    this.CurrentStep = meta.attr("step")
    this.id = window.location.pathname.match(/\d+/)[0]

    return this
}
