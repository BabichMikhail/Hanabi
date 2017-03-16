function adminHandler() {
    this.CreateStat = function(url) {
        // @todo
    }

    this.UpdateStats = function(url, idName) {
        $.ajax({
            type: "GET",
            url: url,
            data: {}
        }).done(function(data) {
            stats = data.stats
            if (!stats) {
                stats = []
            }

            let html = `<div class="col-md-12">
                <table class=table table-striped">
                    <thead>
                        <th>#</th>
                        <th>AITypes</th>
                        <th>Places</th>
                        <th>Points</th>
                        <th>Execution Time</th>
                        <th>Ready at</th>
                        <th>Created at</th>
                    </thead>
                    <tbody>
                    </tbody>
                </table>
            `
            for (let i = 0; i < stats.length; ++i) {
                stat = stats[i]
                // @todo
            }
            $("#stats").html(html)
            setTimeout(Admin.Update, 10000)
        }).fail(function(data) {
            setTimeout(Admin.Update, 3000)
        })
    }

    this.CreateStatPage = function() {
        $.ajax({
            type: "GET",
            url: "/api/ai/names",
            data: {}
        }).done(function(data) {
            console.log(data)
            let html = ``
            aiNames = data.data
            console.log(aiNames)
            console.log(aiNames[0])
            let aiPlayersHtml = ``
            console.log(aiNames.length)
            for (let j in [ 'h', 'e', 'l', 'l', 'o' ]) {
                aiPlayersHtml += `<select id="ai-type-` + j + `">`
                for (let i in aiNames) {
                    aiPlayersHtml += `<option value="` + i + `" ` + (i == 3 ? "selected" : "") + `>` + aiNames[i] + `</option>`
                }
                aiPlayersHtml += `</select>`
            }

            html += `<div class="col-md-4">
                <select id="ai-count">
                    <option value="2">2</option>
                    <option value="3">3</option>
                    <option value="4">4</option>
                    <option value="5" selected>5</option>
                </select>
            </div>
            <div class="col-md-4">` +
                aiPlayersHtml + `
            </div>
            <div class="col-md-4">
                <label for="ai-games-count">Thouthands of games</label>
                <input type="number" id="ai-games-count"></input>
            </div>
            <div class="col-md-12">
                <button type="submit" onClick="Admin.CreateStat()">Create</button>
            </div>
            `
            console.log(aiPlayersHtml)
            $('#stats').html(html)
            setTimeout(Admin.Update, 60000)
        }).fail(function(data) {
            setTimeout(Admin.Update, 3000)
        })
    }

    this.CreateStat = function() {
        let gamesCount = $("input[id='ai-games-count']").val()
        let aiCount = $("select[id='ai-count']").val()
        console.log(gamesCount)
        console.log(aiCount)
        let aiTypes = []
        while (aiCount > 0) {
            --aiCount
            aiTypes[aiCount] = +$("select[id='ai-type-" + aiCount + "']").val()
        }
        console.log(aiTypes)
        $.ajax({
            type: "POST",
            url: "/api/admin/stats/create",
            data: {
                count: gamesCount,
                ai_types: JSON.stringify(aiTypes),
            },
        }).done(function(data) {
            console.log(data.status)
        })
    }

    this.Update = function() {
        Admin.Tabs[Admin.State].Action()
    }

    this.SetActive = function(elem, idName) {
        $("a[class='nav-link active']").removeClass("active")
        elem.classList.add("active")
        Admin.State = idName
        Admin.Update()
    }

    setTimeout(this.Update, 10000)
    this.State = "admin-main"

    this.Tabs = {
        "admin-main": {
            Url: "/api/admin/stats/read",
            Action: function() {
                return Admin.UpdateStats(this.Url)
            },
        },
        "admin-new-stats": {
            Url: "/api/admin/stats/new",
            Action: function() {
                return Admin.CreateStatPage(this.Url)
            },
        },
    }

    return this
}
