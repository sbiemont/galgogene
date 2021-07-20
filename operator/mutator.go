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

// ------------------------------

// Mutate defines a random mutation of bits (with a rate in [0 ; 1])
// * 0: no mutation will happen
// * 1: all bits **may** be mutated
type Mutate struct {
	Rate float64 // The mutation rate on the current bits
}

// Mate will mutate the first set of bits and leave the second one unchanged
func (mut Mutate) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return mutate(bits1, mut.Rate, func(bits gene.Bits, _ int) uint8 {
		return bits.Rand()
	}), bits2
}

// ------------------------------

// Invert defines a random invertion of bits (with a rate in [0 ; 1])
// * 0: no invertion will happen
// * 1: all bits are inverted
type Invert struct {
	Rate float64 // The invertion rate on the current bits
}

// Mate will invert the first set of bits and leave the second one unchanged
func (mut Invert) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	return mutate(bits1, mut.Rate, func(bits gene.Bits, i int) uint8 {
		return bits.Invert(i)
	}), bits2
}

// ------------------------------

// ProbaMutator is a probabilistic mutator
type ProbaMutator struct {
	rate float64 // Mutation rate
	mut  Mutator // Mutation operator
}

// NewProbaMutator build a new full instance of ProbaMutator
func NewProbaMutator(rate float64, mut Mutator) ProbaMutator {
	return ProbaMutator{
		rate: rate,
		mut:  mut,
	}
}

// MultiMutator defines a serie of mutators with a specific probability of beeing chosen.
// All or no mutators may be applied
type MultiMutator []ProbaMutator

func (mm MultiMutator) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
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
func mutate(bits gene.Bits, rate float64, fct func(gene.Bits, int) uint8) gene.Bits {
	result := bits.Clone()
	for i := 0; i < result.Len(); i++ {
		if random.Peek(rate) {
			result.Raw[i] = fct(result, i)
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
