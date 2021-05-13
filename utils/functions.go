package utils

import (
	"math/rand"
	"time"
)

func GetRandomSliceObject(input []string) (string, int) {
	if len(input) == 1 {
		return input[0], 0
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(input))
	randomObject := input[randomIndex]
	return randomObject, randomIndex
}
