package models

import (
	"strconv"
	"strings"
)

func IntSliceToString(slice []int) string {
	values := []string{}
	for _, i := range slice {
		values = append(values, strconv.Itoa(i))
	}
	return strings.Join(values, ", ")
}
