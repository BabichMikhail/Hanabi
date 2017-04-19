function lobbyHandler() {
    this.Create = function (isAIGame) {
        let createUrl = isAIGame ? "/api/lobby/createAI" : "/api/lobby/create"
        let count = $("input[name=playersCount]").val()
        count = count < 2 ? 2 : count
        count = count > 5 ? 5 : count
        $.ajax({
            type: "POST",
            url: createUrl,
            data: {
                playersCount: count,
            },
        }).done(function(result) {
            console.log(result)
            if (result.status == "success") {
                if (Lobby.State == "lobby-finished-games") {
                    return
                }
                gameData = result.data
                let gameId = gameData.id
                let html = `<tr id="game-` +  gameId + `">
                    <th scope="row">` + gameId + `</th>
                    <td>` + gameData.owner + `</td>
                    <td>` + gameData.owner + `</td>
                    <td>` + count  + `</td>
                    <td>` + gameData.status + `</td>
                    <td>
                        <a class="btn-link" href="#" onclick="Lobby.Leave(` + gameId + `)">Leave</button>
                    </td>
                </tr>`
                $("#games").prepend(html)
            }
        }).fail(function(result) {
            console.log("create fail")
        })
    }

    this.Leave = function(id) {
        $.ajax({
            type: "POST",
            url: "/api/lobby/leave/" + id,
            data: {}
        }).done(function(result) {
            console.log(result)
            if (result.status == "success") {
                if (result.data == "delete") {
                    $("#game-" + id).remove()
                } else {
                    Lobby.Update()
                }
            }
        }).fail(function(result) {
            console.log("leave fail")
        })
    }

    this.LoadUser = function() {
        $.ajax({
            type: "GET",
            url: "/api/users/current",
            data: {}
        }).done(function(result) {
            console.log(result)
            if (result.status == "success") {
                Lobby.User = result.data
            }
        }).fail(function(result) {
            console.log("load user fail")
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
        }).done(function(result) {
            console.log(result)
            let data = result.data
            if (result.status == "success" && data.game_status == "active") {
                setTimeout(Lobby.OpenGame, 1000, data.URL)
            }
            Lobby.Update()
        }).fail(function(result) {
            console.log("join fail")
        })
    }

    this.State = "lobby"

    this.UpdateGames = function(url, idName) {
        $.ajax({
            type: "GET",
            url: url,
            data: {}
        }).done(function(result) {
            newHtml = ``
            games = result.data
            if (!games) {
                games = []
            }
            for (let i = 0; i < games.length; ++i) {
                game = games[i]
                playersHtml = ``
                for (let j = 0; j < game.players.length; ++j) {
                    playersHtml += (j != 0 ? ` ` : ``) + game.players[j].nick_name
                }
                let actionHtml = ``
                if (game.status_name == "finished") {
                    actionHtml = `<a class="btn-link" href="/games/view/` + game.id + `">Replay</a>`
                } else if (game.status_name == "wait") {
                    actionHtml = game.user_in
                        ? `<a class="btn-link" href="#" onclick="Lobby.Leave(` + game.id + `)">Leave</a>`
                        : `<a class="btn-link" href="#" onclick="Lobby.Join(` + game.id + `)">Join</a>`
                } else if (game.status_name == "active" && game.user_in) {
                    actionHtml = `<a class="btn-link" href="/games/room/` + game.id + `">Go</a>`
                }

                newHtml += `<tr id="game-` + game.id + `">
                    <th scope="row">` + game.id + `</th>
                    <td>` + game.owner_name + `</td>
                    <td>` + playersHtml + `</td>
                    ` + (idName == 'lobby-finished-games' ? `<td>` + game.points + `</td>` : ``) +`
                    <td>` + game.player_count + `</td>
                    <td>` + game.status_name + `</td>
                    <td>` + actionHtml + `</td>
                </tr>`
                tableHeadHtml =
                    `<th>#</th>` +
                    `<th>Creator</th>` +
                    `<th>Users</th>` +
                    (idName == 'lobby-finished-games' ? `<th>Points</th>` : ``) +
                    `<th>Places</th>` +
                    `<th>Status</th>` +
                    `<th></th>`
            }
            $("#games").html(newHtml)
            $("#table-head").html(tableHeadHtml)
            Lobby.timeout = setTimeout(Lobby.Update, 10000)
        }).fail(function(result) {
            console.log("update games fail")
        })
    }

    this.Update = function() {
        clearTimeout(Lobby.timeout)
        Lobby.Tabs[Lobby.State].Update()
    }

    this.SetActive = function(elem, idName) {
        $("a[class='nav-link active']").removeClass("active")
        elem.classList.add("active")
        Lobby.State = idName
        Lobby.Update()
    }

    this.Init = function () {
        $.ajax({
            type: "GET",
            url: "/api/lobby/games/active",
            data: {}
        }).done(function(result) {
            console.log(result)
        }).fail(function(result) {
            console.log("init fail")
        })
    }

    this.Init()
    this.State = "lobby-main"

    this.Tabs = {
        "lobby-main": {
            State: "lobby-main",
            Url: "api/lobby/games/active",
            Update: function() {
                return Lobby.UpdateGames(this.Url, this.State)
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

    this.timeout = setTimeout(this.Update, 5000)
    return this
}
