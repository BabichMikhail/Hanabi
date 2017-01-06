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
                game_id: Game.Id,
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
                game_id: Game.Id,
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
                game_id: Game.Id,
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
                game_id: Game.Id,
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


    this.Cards = this.Init()
    this.Id = window.location.pathname.match(/\d+/)[0]
    return this
}

console.log(window.location.pathname.match(/\/games\/room\/.*/))
if (window.location.pathname.match(/\/games\/room\/.*/)) {
    window.Game = new gameHandler()
}
