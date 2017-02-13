function viewHandler() {
    this.NewCard = function(color, knownColor, value, knownValue) {
        return {
            color: color,
            knownColor: knownColor,
            value: value,
            knownValue: knownValue,
        }
    }

    this.NewAction = function(type, playerPosition, value) {
        return {
            type: type,
            pos: playerPosition,
            value: value,
        }
    }

    return this
}
