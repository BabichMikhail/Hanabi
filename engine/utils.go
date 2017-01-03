package engine

import "math/rand"

func RandomIntPermutation(values []int) []int {
	for i := len(values) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		value := values[i]
		values[i] = values[j]
		values[j] = value
	}
	return values
}
