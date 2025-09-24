package utils

import (
	"math/rand"
)

// Dice simulates a dice roll of numDice dice with sizeDice sides each
// This mimics the original DikuMUD dice() function
func Dice(numDice, sizeDice int) int {
	if sizeDice < 1 {
		sizeDice = 1
	}
	
	sum := 0
	for i := 0; i < numDice; i++ {
		sum += rand.Intn(sizeDice) + 1
	}
	
	return sum
}

// Number creates a random number in interval [from;to]
// This mimics the original DikuMUD number() function
func Number(from, to int) int {
	if from > to {
		from, to = to, from
	}
	
	return rand.Intn(to-from+1) + from
}
