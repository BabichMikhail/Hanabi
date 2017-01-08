package game

import "math/rand"

func RandomIntPermutation(values []int) []int {
	for i := len(values) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
	return values
}
