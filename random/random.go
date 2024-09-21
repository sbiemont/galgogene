package random

import (
	"math/rand/v2"
	"sort"
)

// random: group here all calls to package "math/rand"
var gen *rand.Rand

func init() {
	Seed(rand.Uint64())
}

func Seed(seed uint64) {
	gen = rand.New(rand.NewPCG(42, seed))
}

// OrderedInts builds an ordered list of k random integers in [min ; max[
// ex: (1, 2, 2, 9)
func OrderedInts(min, max, k int) []int {
	dm := max - min
	result := make([]int, k)
	for i := range k {
		result[i] = gen.IntN(dm) + min
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

// UInt64 returns a random uint64
func Uint64() uint64 {
	return gen.Uint64()
}

// Peek checks if the random generated rate in [0 ; 1[ matches the given rate
func Peek(rate float64) bool {
	return Percent() < rate
}

// Percent returns a random percentage in [0 ; 1[
func Percent() float64 {
	return gen.Float64()
}

// Perm returns a permutation of n ints
func Perm(n int) []int {
	return gen.Perm(n)
}

// IntN returns a random int in [0; n[
func IntN(n int) int {
	return gen.IntN(n)
}

// Shuffle randomizes the order of elements
func Shuffle[T any](items []T) {
	gen.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
}
