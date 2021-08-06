package operator

import (
	"genalgo.git/gene"
	"genalgo.git/random"
)

// CrossOver defines the method to be used for mutating a selection of 2 set of bits
type CrossOver interface {
	// Mate 2 codes to generate 2 new codes (with the same size)
	Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits)
}

// ------------------------------

// https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)

// OnePointCrossOver performs cross-over with 1 randomly chosen point
type OnePointCrossOver struct{}

func (OnePointCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return crossOver(bits1, bits2, random.Ints(0, bits1.Len(), 1))
}

// TwoPointsCrossOver performs cross-over with 2 randomly chosen points
type TwoPointsCrossOver struct{}

func (TwoPointsCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return crossOver(bits1, bits2, random.Ints(0, bits1.Len(), 2))
}

// UniformCrossOver performs a bit by bit cross-over from both parents with an equal probability of beeing chosen
type UniformCrossOver struct{}

func (UniformCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return uniformCrossOver(bits1, bits2, 0.5)
}

// DavisOrderCrossOver performs a Davis' order crossover (permutation)
type DavisOrderCrossOver struct{}

func (DavisOrderCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	pos := random.Ints(0, bits1.Len(), 2)
	return davisOrderCrossOver(bits1, bits2, pos[0], pos[1]), davisOrderCrossOver(bits2, bits1, pos[0], pos[1])
}

// ------------------------------

// ProbaCrossOver is a probabilistic crossover
type ProbaCrossOver struct {
	rate float64   // Crossover rate
	co   CrossOver // Crossover operator
}

// NewProbaCrossOver build a new full instance of ProbaCrossOver
func NewProbaCrossOver(rate float64, co CrossOver) ProbaCrossOver {
	return ProbaCrossOver{
		rate: rate,
		co:   co,
	}
}

// MultiCrossOver defines a serie of crossovers with a specific probability of beeing chosen.
// All or no crossovers may be applied
type MultiCrossOver []ProbaCrossOver

func (mco MultiCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	res1, res2 := bits1, bits2
	for _, m := range mco {
		if random.Peek(m.rate) {
			res1, res2 = m.co.Mate(res1, res2)
		}
	}
	return res1, res2
}

// ------------------------------

// crossOver bits #1 with #2 using an ordered list of indexes
// Returns the 2 resulting set of bits
// index:  [0 1 2 3 4 5 6 7]
// bits 1: [0 0 0 0 0 0 0 0]
// bits 2: [1 1 1 1 1 1 1 1]
// Example
// indexes:   [    2   4   6  ]
// result 1:  [0 0 1 1 0 0 1 1]
// result 2:  [1 1 0 0 1 1 0 0]
func crossOver(bits1, bits2 gene.Bits, indexes []int) (gene.Bits, gene.Bits) {
	var gA, gB *gene.Bits = &bits1, &bits2

	sz := bits1.Len()
	res1 := gene.NewBits(sz, bits1.MaxValue)
	res2 := gene.NewBits(sz, bits1.MaxValue)

	i1 := 0
	for _, i2 := range append(indexes, sz) {
		// Copy sub-slices
		if i1 != i2 {
			_ = copy(res1.Raw[i1:i2], (*gA).Raw[i1:i2])
			_ = copy(res2.Raw[i1:i2], (*gB).Raw[i1:i2])
		}

		// Swap
		var tmp *gene.Bits = gA
		gA = gB
		gB = tmp

		// Next index
		i1 = i2
	}

	return res1, res2
}

// uniformCrossOver swap bits with uniform distribution
// bits 1: [0 0 0 0 0 0 0 0]
// bits 2: [1 1 1 1 1 1 1 1]
// Example
// result 1: [0 1 1 0 0 0 1 0]
// result 2: [1 0 0 1 1 1 0 1]
func uniformCrossOver(bits1, bits2 gene.Bits, rate float64) (gene.Bits, gene.Bits) {
	sz := bits1.Len()
	res1 := gene.NewBits(sz, bits1.MaxValue)
	res2 := gene.NewBits(sz, bits1.MaxValue)

	for i := 0; i < sz; i++ {
		if random.Peek(rate) {
			// Copy without change
			res1.Raw[i] = bits1.Raw[i]
			res2.Raw[i] = bits2.Raw[i]
		} else {
			// Swap values
			res1.Raw[i] = bits2.Raw[i]
			res2.Raw[i] = bits1.Raw[i]
		}
	}

	return res1, res2
}

func davisOrderCrossOver(bits1, bits2 gene.Bits, pos1, pos2 int) gene.Bits {
	// Find unused value
	sz := bits1.Len()
	uniq := make(map[uint8]interface{})
	var idx int
	unusedValue := func() uint8 {
		for idx < sz {
			value := bits2.Raw[idx]
			_, ok := uniq[value]
			idx++
			if !ok {
				uniq[value] = nil
				return value
			}
		}
		return 0
	}

	res := gene.NewBits(sz, bits1.MaxValue)

	// Copy range part
	for i := pos1; i <= pos2; i++ {
		value := bits1.Raw[i]
		res.Raw[i] = value
		uniq[value] = nil
	}

	// Fill begining with unused values
	for i := 0; i < pos1; i++ {
		res.Raw[i] = unusedValue()
	}

	// Fill ending with unused values
	for i := pos2 + 1; i < sz; i++ {
		res.Raw[i] = unusedValue()
	}

	return res
}
