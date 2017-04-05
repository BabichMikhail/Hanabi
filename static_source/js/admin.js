function adminHandler() {
    this.UpdateStats = function(url) {
        $.ajax({
            type: "GET",
            url: url,
            data: {}
        }).done(function(data) {
            stats = data.data
            if (!stats) {
                stats = []
            }

            let html = `<div class="col-md-12">
                <table class="table table-striped">
                    <thead>
                        <th>#</th>
                        <th>AITypes</th>
                        <th>Count</th>
                        <th>Places</th>
                        <th>Points</th>
                        <th>ExecTime</th>
                        <th>ReadyAt</th>
                        <th>CreatedAt</th>
                        <th></th>
                    </thead>
                    <tbody>
            `
            for (let i = 0; i < stats.length; ++i) {
                let stat = stats[i]
                html += `<tr id="stat-` + stat.id + `">
                    <th>` + stat.id + `</th>
                    <td>` + stat.ai_names.join(' ') + `</td>
                    <td>` + stat.count + `</td>
                    <td>` + stat.ai_types.length + `</td>
                    <td>` + stat.points + `</td>
                    <td>` + stat.execution_time + `</td>
                    <td>` + stat.ready_at + `</td>
                    <td>` + stat.created_at + `</td>
                    <td><button class="btn btn-link" onClick="Admin.DeleteStat(` + stat.id + `)">Delete</button></td>
                </tr>`
            }
            html += `</tbody>
                </table>
            </div>`
            $("#stats").html(html)
            Admin.timeout = setTimeout(Admin.Update, 10000)
        }).fail(function(data) {
            Admin.timeout = setTimeout(Admin.Update, 3000)
        })
    }

    this.UpdateCreateStatPageCount = function() {
        let count = $("#ai-count").val()
        Admin.Tabs["admin-new-stats"] = {
            Url: "/api/admin/stats/new",
            Count: count,
            Action: function() {
                return Admin.CreateStatPage(count)
            },
        }
        Admin.Update()
    }

    this.CreateStatPage = function(count) {
        $.ajax({
            type: "GET",
            url: "/api/ai/names",
            data: {}
        }).done(function(data) {
            let html = ``
            aiNames = data.data
            let aiPlayersHtml = ``
            for (let j in [ 'h', 'e', 'l', 'l', 'o' ].slice(5 - count)) {
                aiPlayersHtml += `<select id="ai-type-` + j + `">`
                for (let i in aiNames) {
                    aiPlayersHtml += `<option value="` + i + `" ` + (i == 3 ? "selected" : "") + `>` + aiNames[i] + `</option>`
                }
                aiPlayersHtml += `</select>`
            }

            html += `<div class="col-md-4">
                <select id="ai-count" onChange="Admin.UpdateCreateStatPageCount()">`
            for (let i = 1; i <= 5; ++i) {
                html += `<option value="` + i + `" ` + (i == count ? "selected" : "") + `>` + i + `</option>`
            }
            html += `</select>
            </div>
            <div class="col-md-4">` +
                aiPlayersHtml + `
            </div>
            <div class="col-md-4">
                <label for="ai-games-count">Thouthands of games</label>
                <input type="number" id="ai-games-count"></input>
            </div>
            <div class="col-md-12">
                <label for="save-distribution-in-excel">Save distribution points in Excel</label>
                <input type="checkbox" id="save-distribution-in-excel" value="true"></input>
            <div class="col-md-12">
                <button type="submit" onClick="Admin.CreateStat()">Create</button>
            </div>
            `
            $('#stats').html(html)
        }).fail(function(data) {
            Admin.timeout = setTimeout(Admin.Update, 3000)
        })
    }

    this.CreateStat = function() {
        let gamesCount = $("input[id='ai-games-count']").val()
        let aiCount = $("select[id='ai-count']").val()
        let saveInExcel = $("input[id='save-distribution-in-excel']").val()
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
                save_distribution_in_excel: saveInExcel,
            },
        }).done(function(data) {
            if (data.status != "success") {
                console.log("Create status: " + data.status)
            }
        })
    }

    this.DeleteStat = function(id) {
        $.ajax({
            type: "POST",
            url: "/api/admin/stats/delete",
            data: {
                id: id,
            },
        }).done(function(data) {
            if (data.status != "success") {
                console.log("Delete status: " + data.status)
            }
            $(`tr[id="stat-` + id + `"]`).remove()
        })
    }

    this.Update = function() {
        clearTimeout(Admin.timeout)
        Admin.Tabs[Admin.State].Action()
    }

    this.SetActive = function(elem, idName) {
        $("a[class='nav-link active']").removeClass("active")
        elem.classList.add("active")
        Admin.State = idName
        Admin.Update()
    }

    this.timeout = setTimeout(this.Update, 10000)
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
            Count: 5,
            Action: function() {
                return Admin.CreateStatPage(this.Count)
            },
        },
    }

    return this
}
