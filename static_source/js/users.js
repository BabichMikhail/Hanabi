function userHandler() {
    this._Current = null

    this.UpdateCurrent = function () {
        $.ajax({
            type: "GET",
            url: "/api/users/current",
            data: {},
            async: false,
        }).done(function(data) {
            console.log(data)
            if (data.status == "OK") {
                this._Current = data.user
            }
        }).fail(function(data) {
            alert("FAIL UPDATE CURRENT USER")
        })
        return this._Current
    }

    this.GetCurrent = function() {
        if (this._Current == null) {
            return this.UpdateCurrent()
        }
        return this._Current
    }

    return this
}

User = new userHandler()
