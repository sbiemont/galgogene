package operator

import (
	"galgogene.git/gene"
	"galgogene.git/random"
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

// UniformOrderCrossOver performs a uniform order crossover (permutation)
type UniformOrderCrossOver struct{}

func (UniformOrderCrossOver) Mate(bits1, bits2 gene.Bits) (gene.Bits, gene.Bits) {
	var mask0 []int
	var mask1 []int
	for i := 0; i < bits1.Len(); i++ {
		if random.Peek(0.5) {
			mask1 = append(mask1, i)
		} else {
			mask0 = append(mask0, i)
		}
	}

	return uniformOrderCrossOver(bits1, bits2, mask0, mask1), uniformOrderCrossOver(bits2, bits1, mask0, mask1)
}

// ------------------------------

// probaCrossOver is a probabilistic crossover
type probaCrossOver struct {
	rate float64   // Crossover rate
	co   CrossOver // Crossover operator
}

// MultiCrossOver defines a serie of crossovers with a specific probability of beeing chosen.
// All or no crossovers may be applied
type MultiCrossOver []probaCrossOver

// Use the given proba crossover
func (mco MultiCrossOver) Use(rate float64, co CrossOver) MultiCrossOver {
	return append(mco, probaCrossOver{
		rate: rate,
		co:   co,
	})
}

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
	res1 := gene.NewBitsFrom(bits1)
	res2 := gene.NewBitsFrom(bits1)

	i1 := 0
	for _, i2 := range append(indexes, bits1.Len()) {
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
	res1 := gene.NewBitsFrom(bits1)
	res2 := gene.NewBitsFrom(bits1)

	for i := 0; i < bits1.Len(); i++ {
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

type finder struct {
	idx  int
	uniq map[uint8]interface{}
}

func newFinder() finder {
	return finder{
		uniq: make(map[uint8]interface{}),
	}
}

func (fnd *finder) nextUnused(bits gene.Bits) uint8 {
	for fnd.idx < bits.Len() {
		value := bits.Raw[fnd.idx]
		_, ok := fnd.uniq[value]
		fnd.idx++
		if !ok {
			fnd.used(value)
			return value
		}
	}
	return 0
}

func (fnd finder) used(value uint8) {
	fnd.uniq[value] = nil
}

func davisOrderCrossOver(bits1, bits2 gene.Bits, pos1, pos2 int) gene.Bits {
	// Find unused value
	fnd := newFinder()
	res := gene.NewBitsFrom(bits1)

	// Copy range part
	for i := pos1; i <= pos2; i++ {
		value := bits1.Raw[i]
		res.Raw[i] = value
		fnd.used(value)
	}

	// Fill begining with unused values
	for i := 0; i < pos1; i++ {
		res.Raw[i] = fnd.nextUnused(bits2)
	}

	// Fill ending with unused values
	for i := pos2 + 1; i < bits1.Len(); i++ {
		res.Raw[i] = fnd.nextUnused(bits2)
	}

	return res
}

func uniformOrderCrossOver(bits1, bits2 gene.Bits, mask0 []int, mask1 []int) gene.Bits {
	res := gene.NewBitsFrom(bits1)
	fnd := newFinder()

	// mask1: copy values and add to uniq
	for _, idx := range mask1 {
		value := bits1.Raw[idx]
		res.Raw[idx] = value
		fnd.used(value)
	}

	// mask0: value has to be found in unused values
	for _, idx := range mask0 {
		res.Raw[idx] = fnd.nextUnused(bits2)
	}

	return res
}
