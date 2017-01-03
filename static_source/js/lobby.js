function lobbyHandler() {
    this.Create = function () {
        $.ajax({
            type: "POST",
            url: "/api/lobby/create",
            data: {
                playersCount: $("input[name=playersCount]").val()
            },
        }).done(function(data) {
            console.log(data)
            if (data.status == "OK") {
                let gameId = data.game.Id
                let html = `<tr id="game-` +  gameId + `">
                    <td>` + data.game.Owner + `</td>
                    <td><a href="Game.LoadUserList(` + gameId + `)">1</a></td>
                    <td>
                        <button type="button" onClick="Game.Leave(` + gameId + `);">` + (data.currentUserId == data.game.OwnerId ? `Leave` : `Join`) + `</button>
                    </td>
                </tr>`
                $("#games").append(html)
            }
        }).fail(function(data) {
            alert("CREATE FAIL")
        })
    }

    this.Leave = function(id) {
        $.ajax({
            type: "POST",
            url: "/api/lobby/leave/" + id,
            data: {}
        }).done(function(data) {
            console.log(data)
            if (data.status == "OK") {
                if (data.action == "delete") {
                    $("#game-" + id).remove()
                } else {
                    location.reload()
                }
            }
        }).fail(function(data) {
            alert("LEAVE FAIL")
        })
    }

    this.OpenGame = function(URL) {
        location.reload()
        var win = window.open(URL, '_blank')
        win.focus()
    }

    this.Join = function(id) {
        $.ajax({
            type: "POST",
            url: "/api/lobby/join/" + id,
            data: {}
        }).done(function(data) {
            console.log(data)
            if (data.status == "OK") {
                if (data.game_status == "active") {
                    setTimeout(Lobby.OpenGame, 1000, data.URL)
                }
            }
        }).fail(function(data) {
            alert("JOIN FAIL")
        })
    }

    this.Statuses = null

    this.Update = function() {
        $.ajax({
            type: "GET",
            url: "/api/lobby/status",
            data: {}
        }).done(function(data) {
            console.log(data)
            console.log(Lobby.Statuses)
            if (data.game != null) {
                for (var i = 0; i < data.game.length; ++i) {
                    if (Lobby.Statuses[data.game[i].game_id] == null) {
                        Lobby.Statuses[data.game[i].game_id] = data.game[i]
                        continue
                    }
                    let oldStatus = Lobby.Statuses[data.game[i].game_id].status_code
                    let newStatus = data.game[i].status_code
                    if (oldStatus != newStatus) {
                        Lobby.Statuses[data.game[i].game_id].status_code = newStatus
                        setTimeout(Lobby.OpenGame, 1000, data.game[i].URL)
                    }
                }
            }
            setTimeout(Lobby.Update, 5000)
        }).fail(function(data) {
            alert("UPDATE FAIL")
        })
    }

    this.LoadUserList = function(id) {
        $.ajax({
            type: "GET",
            url: "/api/lobby/users/"  + id,
            data: {}
        }).done(function(data) {
            console.log(data)
            if (data.status == "OK") {
                let html = "<ul>"
                for (let i = 0; i < data.players.length; ++i) {
                    html += "<li>" + data.players[i].nick_name + "</li>";
                }
                html += "</ul>"
                $("#users-" + id).html(html)
            }
        }).fail(function(data) {
            alert("GET USERS FAIL")
        })
    }

    this.Init = function () {
        $.ajax({
            type: "GET",
            url: "/api/lobby/status",
            data: {}
        }).done(function(data) {
            console.log(data)
            Lobby.Statuses = []
            if (data.game != null) {
                for (var i = 0; i < data.game.length; ++i) {
                    Lobby.Statuses[data.game[i].game_id] = data.game[i]
                }
            }
            setTimeout(Lobby.Update, 5000)
        }).fail(function(data) {
            alert("INIT FAIL")
        })
    }

    this.Init()

    return this
}

if (window.location.pathname == "/games") {
    window.Lobby = new lobbyHandler()
}
