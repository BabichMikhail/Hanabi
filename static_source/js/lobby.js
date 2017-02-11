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

    this.LoadUser = function() {
        $.ajax({
            type: "GET",
            url: "/api/users/current",
            data: {}
        }).done(function(data) {
            console.log(data)
            if (data.status == "success") {
                Lobby.User = data.user
            }
        }).fail(function(data) {
            alert("LEAVE FAIL")
        })
    }

    this.OpenGame = function(URL) {
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
            if (data.status == "success" && data.game_status == "active") {
                setTimeout(Lobby.OpenGame, 1000, data.URL)
            }
            location.reload()
        }).fail(function(data) {
            alert("JOIN FAIL")
        })
    }

    this.Statuses = null
    this.State = "lobby"

    this.UpdateLobby = function() {
        $.ajax({
            type: "GET",
            url: "/api/lobby/games/active",
            data: {}
        }).done(function(data) {
            if (typeof Lobby.User == 'undefined') {
                setTimeout(Lobby.Update, 3000)
                return
            }
            console.log(data.games)
            let games = data.games
            if (games != null) {
                for (var i = 0; i < games.length; ++i) {
                    let game = games[i]
                    if (Lobby.Statuses[game.id] == null) {
                        Lobby.Statuses[game.id] = game
                        continue
                    }
                    let oldStatus = Lobby.Statuses[game.id].status_code
                    let newStatus = game.status_code
                    if (oldStatus != newStatus) {
                        Lobby.Statuses[game.id].status_code = newStatus
                        setTimeout(Lobby.OpenGame, 1000, game.URL)
                    }
                }
            }

            let html = ``
            for (let i = 0; i < games.length; ++i) {
                playersHtml = ``
                let userIn = false
                let game = games[i]
                let ownerName = ''
                for (let j = 0; j < game.players.length; ++j) {
                    playersHtml += (j > 0 ? ` ` : ``) + game.players[j].nick_name
                    if (game.players[j].id == Lobby.User.id) {
                        userIn = true
                    }
                }

                statusHtml =
                    games[i].status_name == `active`
                        ? `<a href="` + game.URL + `">GO</a>`
                        : (userIn
                            ? `<button type="button" onclick="Lobby.Leave(` + game.id + `)">Leave</button>`
                            : `<button type="button" onclick="Lobby.Join(` + game.id + `)">Join</button>`
                        )
                html += `<tr id = "game-` + game.id + `">
                    <th scope="row">` + game.id + `</th>
                    <td>` + game.owner_name + `</td>
                    <td>` + playersHtml + `</td>
                    <td>` + game.player_count + `</td>
                    <td>` + statusHtml + `</td>
                </tr>`
            }
            $("#games").html(html)
            setTimeout(Lobby.Update, 10000)
        }).fail(function(data) {
            alert("UPDATE LOBBY FAIL")
        })
    }

    this.UpdateGames = function(url, idName) {
        $.ajax({
            type: "GET",
            url: url,
            data: {}
        }).done(function(data) {
            html = ``
            games = data.games
            for (let i = 0; i < games.length; ++i) {
                game = games[i]
                playersHtml = ``
                for (let j = 0; j < game.players.length; ++j) {
                    playersHtml += (j != 0 ? ` ` : ``) + game.players[j].nick_name
                }
                html += `<div class="col-md-2">
                    <div class="card text-center">
                        <div class="card-block">
                            <h4 class="card-title">Game #` + game.id + `</h4>
                            <p class="card-text">Creator: ` + game.owner_name + `
                            <p class="card-text">Status: ` + game.status_name + `</p>
                            <p class="card-text">Places: ` + game.player_count + `</p>
                            <p class="card-text">Players: ` + playersHtml + `</p>
                            ` + (game.status == 4
                                ? `<a class="card-link" href="/games/view/` + game.id + `">Replay</a>`
                                : ``
                            ) + `
                        </div>
                    </div>
                </div>`
            }
            $("#" + idName).html(html)
            setTimeout(Lobby.Update, 10000)
        }).fail(function(data) {
            alert("UPDATE GAMES FAIL")
        })
    }

    this.Update = function() {
        Lobby.Tabs[Lobby.State].Update()
    }

    this.SetActive = function(elem, idName) {
        $("a[class='nav-link active']").removeClass("active")
        elem.classList.add("active")
        for (let state in Lobby.Tabs) {
            $("#" + state).addClass("invisible")
        }
        $("#" + idName).removeClass("invisible")
        Lobby.State = idName
        Lobby.Update()
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
            url: "/api/lobby/games/active",
            data: {}
        }).done(function(data) {
            console.log(data)
            Lobby.Statuses = []
            if (data.games != null) {
                for (let i = 0; i < data.games.length; ++i) {
                    Lobby.Statuses[data.games[i].game_id] = data.games[i]
                }
            }
            setTimeout(Lobby.Update, 5000)
        }).fail(function(data) {
            alert("INIT FAIL")
        })
    }

    this.Init()

    this.Tabs = {
        "lobby-main": {
            State: "lobby-main",
            Url: "api/lobby/games/active",
            Update: function() {
                return Lobby.UpdateLobby(this.Url, this.State)
            },
        },
        "lobby-my-games": {
            State: "lobby-my-games",
            Url: "/api/lobby/games/my",
            Update: function() {
                return Lobby.UpdateGames(this.Url, this.State)
            },
        },
        "lobby-all-games": {
            State: "lobby-all-games",
            Url: "/api/lobby/games/all",
            Update: function() {
                return Lobby.UpdateGames(this.Url, this.State)
            },
        },
        "lobby-finished-games": {
            State: "lobby-finished-games",
            Url: "/api/lobby/games/finished",
            Update: function() {
                return Lobby.UpdateGames(this.Url, this.State)
            },
        },
    }

    return this
}

if (window.location.pathname == "/games") {
    window.Lobby = new lobbyHandler()
    window.Lobby.LoadUser()
}
