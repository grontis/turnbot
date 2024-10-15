package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func RollDice(sides int) int {
	seed := uint64(time.Now().UnixNano())
	rand.Seed(seed)
	return rand.Intn(sides) + 1
}
