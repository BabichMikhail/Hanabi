function lobbyHandler() {
    this.Create = function () {
        let count = $("input[name=playersCount]").val()
        count = count < 2 ? 2 : count
        count = count > 5 ? 5 : count
        $.ajax({
            type: "POST",
            url: "/api/lobby/create",
            data: {
                playersCount: count,
            },
        }).done(function(data) {
            console.log(data)
            if (data.status == "success") {
                let gameId = data.game.id
                let html = `<tr id="game-` +  gameId + `">
                    <th scope="row">` + gameId + `</th>
                    <td>` + data.game.owner + `</td>
                    <td>` + data.game.owner + `</td>
                    <td>` + count  + `</td>
                    <td>
                        <button type="button" onClick="Lobby.Leave(` + gameId + `);">Leave</button>
                    </td>
                </tr>`
                $("#games").prepend(html)

                queueHtml = `<div id="queue-game-` +  gameId + `" class="card text-center">
                    <div class="card-block">
                        <h4 class="card-title">Game #` + gameId + `</h4>
                        <p class="card-text">Status: ` + data.game.status + `</p>
                        <p class="card-text">Places: ` + count + `</p>
                        <p class="card-text">Players: ` + data.game.owner + `</p>
                    </div>
                </div>`
                let childLeft = $('#queue-1')
                let childRight = $('#queue-2')
                let elem = childLeft[0].childElementCount > childRight[0].childElementCount ? childRight : childLeft
                elem.append(queueHtml)
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
            if (data.status == "success") {
                if (data.action == "delete") {
                    $("#game-" + id).remove()
                    $("#queue-game-" + id).remove()
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
            if (data.status == "success") {
                if (data.game_status == "active") {
                    setTimeout(Lobby.OpenGame, 1000, data.URL)
                }
            }
            location.reload()
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
            let myGames = data.game
            let allGames = data.allGames
            if (myGames != null) {
                for (var i = 0; i < data.game.length; ++i) {
                    if (Lobby.Statuses[myGames[i].game_id] == null) {
                        Lobby.Statuses[myGames[i].game_id] = myGames[i]
                        continue
                    }
                    let oldStatus = Lobby.Statuses[myGames[i].game_id].status_code
                    let newStatus = myGames[i].status_code
                    if (oldStatus != newStatus) {
                        Lobby.Statuses[myGames[i].game_id].status_code = newStatus
                        setTimeout(Lobby.OpenGame, 1000, myGames[i].URL)
                    }
                }
            }
            setTimeout(Lobby.Update, 10000)
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
            if (data.status == "success") {
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
