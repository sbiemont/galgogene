package random

import "math/rand"

// Ints builds an ordered list of k random integers in [min ; max[
func Ints(min, max, k int) []int {
	result := make([]int, k)
	last := min
	for i := 0; i < k; i++ {
		rd := inRangeInt(last, max)
		result[i] = rd
		last = rd
	}
	return result
}

// Bit returns a random bit 0 or 1
func Bit() uint8 {
	return uint8(rand.Intn(2))
}

// Peek checks if the random generated rate in [0 ; 1[ matches the given one
func Peek(rate float64) bool {
	return Percent() < rate
}

// Percent returns a random percentage in [0 ; 1[
func Percent() float64 {
	return rand.Float64()
}

func inRangeInt(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}
