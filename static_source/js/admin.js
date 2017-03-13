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
            let newHtml = ``
            stats = data.stats
            if (!stats) {
                stats = []
            }
            for (let i = 0; i < stats.length; ++i) {
                stat = stats[i]
                // @todo
            }
            $("#stats").html(newHtml)
            setTimeout(Admin.Update, 10000)
        }).fail(function(data) {
            setTimeout(Admin.Update, 3000)
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

    this.Init = function () {
        setTimeout(Admin.Update, 10000)
    }

    this.Init()
    this.State = "admin-main"

    this.Tabs = {
        "admin-main": {
            Url: "api/admin/stats/read",
            Action: function() {
                return Admin.UpdateStats(this.Url)
            },
        },
        "admin-new-stats": {
            Url: "api/admin/stats/new",
            Action: function() {
                return Admin.CreateStatPage(this.Url)
            },
        },
    }

    return this
}
