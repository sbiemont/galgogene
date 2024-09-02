package random

import (
	"math/rand"
	"sort"
)

// random: group here all calls to package "math/rand"

// OrderedInts builds an ordered list of k random integers in [min ; max[
// ex: [1, 2, 2, 9]
func OrderedInts(min, max, k int) []int {
	dm := max - min
	result := make([]int, k)
	for i := 0; i < k; i++ {
		result[i] = rand.Intn(dm) + min
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

// Bit returns a random byte in [0 ; 255]
func Byte() uint8 {
	return uint8(rand.Intn(256))
}

// Peek checks if the random generated rate in [0 ; 1[ matches the given rate
func Peek(rate float64) bool {
	return Percent() < rate
}

// Percent returns a random percentage in [0 ; 1[
func Percent() float64 {
	return rand.Float64()
}

// Perm returns a permutation of n ints
func Perm(n int) []int {
	return rand.Perm(n)
}

// Intn returns a random int in [0; n[
func Intn(n int) int {
	return rand.Intn(n)
}
