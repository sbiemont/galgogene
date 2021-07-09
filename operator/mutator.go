package operator

import (
	"genalgo.git/gene"
	"genalgo.git/random"
)

// Mutator defines the method to be used for mutating a selection of 2 set of bits
type Mutator interface {
	// Mate 2 codes to generate 2 new codes (with the same size)
	Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits)
}

// ------------------------------

// https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)

// OnePointCrossOver performs cross-over with 1 randomly choosen point
type OnePointCrossOver struct{}

func (OnePointCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return crossOver(bits1, bits2, random.Ints(0, len(bits1), 1))
}

// TwoPointsCrossOver performs cross-over with 2 randomly choosen points
type TwoPointsCrossOver struct{}

func (TwoPointsCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return crossOver(bits1, bits2, random.Ints(0, len(bits1), 2))
}

// UniformCrossOver performs a bit by bit cross-over from both parents with an equal probability of beeing chosen
type UniformCrossOver struct{}

func (UniformCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return uniformCrossOver(bits1, bits2, 0.5)
}

// ------------------------------

// Mutate defines a random mutation of bits (with a rate in [0 ; 1])
// * 0: no mutation will happen
// * 1: all bits will be inverted
type Mutate struct {
	Rate float64 // The mutation rate on the current bits
}

// Mate will mutate the first set of bits and leave the second one unchanged
func (mut Mutate) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return mutate(bits1, mut.Rate), bits2
}

// ------------------------------

// ProbaMutator is a probabilistic mutator
type ProbaMutator struct {
	rate float64
	mut  Mutator
}

// NewProbaMutator build a new full instance of ProbaMutator
func NewProbaMutator(rate float64, mut Mutator) ProbaMutator {
	return ProbaMutator{
		rate: rate,
		mut:  mut,
	}
}

// MultiMutators defines a serie of mutators with a specific probability of beeing chosen.
// All or no mutators may be applied
type MultiMutators []ProbaMutator

func (mm MultiMutators) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	res1, res2 := bits1, bits2
	for _, m := range mm {
		if random.Peek(m.rate) {
			res1, res2 = m.mut.Mate(res1, res2)
		}
	}
	return res1, res2
}

// ------------------------------

// mutate inverts some bits using a mutation rate
func mutate(bits gene.Bits, rate float64) gene.Bits {
	n := len(bits)
	result := gene.NewBits(n)
	for i := 0; i < n; i++ {
		if random.Peek(rate) {
			result[i] = 1 - bits[i] // invert
		} else {
			result[i] = bits[i] // copy
		}
	}
	return result
}

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

	sz := len(bits1)
	res1 := gene.NewBits(sz)
	res2 := gene.NewBits(sz)

	i1 := 0
	for _, i2 := range append(indexes, sz) {
		// Copy sub-slices
		if i1 != i2 {
			_ = copy(res1[i1:i2], (*gA)[i1:i2])
			_ = copy(res2[i1:i2], (*gB)[i1:i2])
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
	sz := len(bits1)
	res1 := gene.NewBits(sz)
	res2 := gene.NewBits(sz)

	for i := 0; i < sz; i++ {
		if random.Peek(rate) {
			// Copy without change
			res1[i] = bits1[i]
			res2[i] = bits2[i]
		} else {
			// Swap values
			res1[i] = bits2[i]
			res2[i] = bits1[i]
		}
	}

	return res1, res2
}
