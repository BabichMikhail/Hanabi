function gameHandler() {

    this.Init = function () {
        $.ajax({
            type: "GET",
            url: "/api/games/cards",
            data: {}
        }).done(function(data) {
            console.log(data)
            Game.Cards = data
        }).fail(function(data) {
            alert("INIT GAME FAIL")
        })
    }

    this.Init()

    return this
}

//if (window.location.pathname == "/games/room/*") {
    window.Game = new gameHandler()
//}
