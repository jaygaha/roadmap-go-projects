package utils

import (
	"math/rand"
)

// random number generator between min and max
func GenerateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min // returns a random number between min and max
}

// CloseGuessNumber returns the absolute value of an int
func CloseGuessNumber(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
