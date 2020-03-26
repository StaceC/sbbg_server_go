package game

import (
	"math/rand"
)

// NumberGenerator - Interface to allow mocking sequences of numbers
type NumberGenerator interface {
	GetInt() int
}

// RNG - Not seeded (so not random really atm)
type RNG struct {
	Max int
}

// NewRNG - New Random Number Generator
func NewRNG(max int) *RNG {
	rng := &RNG{max}

	return rng
}

// GetInt - Will return an int when called
// Rand includes 0, so pick a number between 0 and MaxNum - 1, then add 1
func (rng *RNG) GetInt() int {
	return rand.Intn(rng.Max-1) + 1
}

// SSNG - Set Sequence Number Generator
type SSNG struct {
	Numbers []int
	Index   int
}

// NewSSNG - New Set Sequence Generator
func NewSSNG(numbers []int) *SSNG {
	gen := &SSNG{
		Numbers: numbers,
		Index:   0,
	}

	return gen
}

// GetInt - Returns the next number in the passed in array
// Todo: Make safe
func (gen *SSNG) GetInt() int {
	i := gen.Index
	gen.Index++
	return gen.Numbers[i]
}
