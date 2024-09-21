package operator

import (
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"
)

// CrossOver defines the method to be used for mutating a selection of 2 chromosomes
type CrossOver interface {
	// Mate 2 codes to generate 2 new codes (with the same size)
	Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome)
}

// ------------------------------

// https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm)

// OnePointCrossOver performs cross-over with 1 randomly chosen point
type OnePointCrossOver struct{}

func (OnePointCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	return crossOver(chrm1, chrm2, random.OrderedInts(0, chrm1.Len(), 1))
}

// ------------------------------

// TwoPointsCrossOver performs cross-over with 2 randomly chosen points
type TwoPointsCrossOver struct{}

func (TwoPointsCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	return crossOver(chrm1, chrm2, random.OrderedInts(0, chrm1.Len(), 2))
}

type ThreePointsCrossOver struct{}

func (ThreePointsCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	return crossOver(chrm1, chrm2, random.OrderedInts(0, chrm1.Len(), 3))
}

// ------------------------------

// UniformCrossOver performs a bit by bit cross-over from both parents with an equal probability of beeing chosen
type UniformCrossOver struct{}

func (UniformCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	return uniformCrossOver(chrm1, chrm2, 0.5)
}

// ------------------------------

// DavisOrderCrossOver performs a Davis' order crossover (permutation)
type DavisOrderCrossOver struct{}

func (DavisOrderCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	pos := random.OrderedInts(0, chrm1.Len(), 2)
	return davisOrderCrossOver(chrm1, chrm2, pos[0], pos[1]), davisOrderCrossOver(chrm2, chrm1, pos[0], pos[1])
}

// ------------------------------

// UniformOrderCrossOver performs a uniform order crossover (permutation)
type UniformOrderCrossOver struct{}

func (UniformOrderCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	var mask0 []int
	var mask1 []int
	for i := range chrm1.Len() {
		if random.Peek(0.5) {
			mask1 = append(mask1, i)
		} else {
			mask0 = append(mask0, i)
		}
	}
	return uniformOrderCrossOver(chrm1, chrm2, mask0, mask1), uniformOrderCrossOver(chrm2, chrm1, mask0, mask1)
}

// ------------------------------

// PartiallyMatchCrossOver (PMX) performs an order crossover (permutation)
type PartiallyMatchCrossOver struct{}

func (PartiallyMatchCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	pos := random.OrderedInts(0, chrm1.Len(), 2)
	return partiallyMatchCrossOver(chrm1, chrm2, pos[0], pos[1]), partiallyMatchCrossOver(chrm2, chrm1, pos[0], pos[1])
}

// ------------------------------

// probaCrossOver is a probabilistic crossover
type probaCrossOver struct {
	rate float64   // Crossover rate
	co   CrossOver // Crossover operator
}

// MultiCrossOver defines a serie of crossovers with a specific probability of beeing chosen.
// All or no crossovers may be applied
type MultiCrossOver struct {
	ApplyAll   bool // Set it to true, otherwise, processing stops at the first crossover to be applied
	crossovers []probaCrossOver
}

// Use the given proba crossover
func (mco MultiCrossOver) Use(rate float64, co CrossOver) MultiCrossOver {
	return MultiCrossOver{
		ApplyAll: mco.ApplyAll,
		crossovers: append(mco.crossovers, probaCrossOver{
			rate: rate,
			co:   co,
		}),
	}
}

func (mco MultiCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	res1, res2 := chrm1, chrm2
	for _, m := range mco.crossovers {
		if random.Peek(m.rate) {
			res1, res2 = m.co.Mate(res1, res2)
			if !mco.ApplyAll {
				return res1, res2
			}
		}
	}
	return res1, res2
}

// ------------------------------

// Helpers

// crossOver chromosome #1 with #2 using an ordered list of indexes
// Returns the 2 resulting set of bases
// index:  [0 1 2 3 4 5 6 7]
// bases 1: [0 0 0 0 0 0 0 0]
// bases 2: [1 1 1 1 1 1 1 1]
// Example
// indexes:   [    2   4   6  ]
// result 1:  [0 0 1 1 0 0 1 1]
// result 2:  [1 1 0 0 1 1 0 0]
func crossOver(chrm1, chrm2 gene.Chromosome, indexes []int) (gene.Chromosome, gene.Chromosome) {
	var gA, gB *gene.Chromosome = &chrm1, &chrm2
	res1 := chrm1.New()
	res2 := chrm2.New()

	i1 := 0
	for _, i2 := range append(indexes, chrm1.Len()) {
		// Copy sub-slices
		if i1 != i2 {
			_ = copy(res1.Raw[i1:i2], (*gA).Raw[i1:i2])
			_ = copy(res2.Raw[i1:i2], (*gB).Raw[i1:i2])
		}

		// Swap and move to next index
		gA, gB = gB, gA
		i1 = i2
	}
	return res1, res2
}

// uniformCrossOver swap bases with uniform distribution
// bases 1: [0 0 0 0 0 0 0 0]
// bases 2: [1 1 1 1 1 1 1 1]
// Example
// result 1: [0 1 1 0 0 0 1 0]
// result 2: [1 0 0 1 1 1 0 1]
func uniformCrossOver(chrm1, chrm2 gene.Chromosome, rate float64) (gene.Chromosome, gene.Chromosome) {
	res1 := chrm1.New()
	res2 := chrm2.New()

	for i := range chrm1.Len() {
		if random.Peek(rate) {
			// Copy without change
			res1.Raw[i] = chrm1.Raw[i]
			res2.Raw[i] = chrm2.Raw[i]
		} else {
			// Swap values
			res1.Raw[i] = chrm2.Raw[i]
			res2.Raw[i] = chrm1.Raw[i]
		}
	}
	return res1, res2
}

// finder
type finder struct {
	idx        int            // current index
	usedValues map[gene.B]int // all values and thiere counts
}

func newFinder() finder {
	return finder{
		usedValues: make(map[gene.B]int),
	}
}

// mark an index as visited
func (fnd finder) useValue(value gene.B) {
	count := fnd.usedValues[value]
	fnd.usedValues[value] = count + 1
}

// find the next unused value in the chromosome
func (fnd *finder) nextUnused(chrm gene.Chromosome) gene.B {
	for fnd.idx < chrm.Len() {
		value := chrm.Raw[fnd.idx]
		count := fnd.usedValues[value]
		fnd.idx++
		if count == 0 {
			return value
		}
		if count > 0 {
			fnd.usedValues[value] = count - 1
		}
	}
	return 0
}

// davisOrderCrossOver replaces values from chrm2 to chrm1
// pos:   0 1 2 3 4 5 6 7 8
// chrm1: A B C D E F G H I
// chrm2: I H G F E D C B A
// pos: [1 ; 3]
// res:   - - C D E F - - - (copy chrm1 from pos1 to pos2)
// res:   I H C D E F G B A (copy chrm2 one by one except data already present in res)
func davisOrderCrossOver(chrm1, chrm2 gene.Chromosome, pos1, pos2 int) gene.Chromosome {
	// Find unused value
	fnd := newFinder()
	res := chrm1.New()

	// Copy range part
	for i := pos1; i <= pos2; i++ {
		value := chrm1.Raw[i]
		res.Raw[i] = value
		fnd.useValue(value)
	}

	// Fill begining with unused values
	for i := range pos1 {
		res.Raw[i] = fnd.nextUnused(chrm2)
	}

	// Fill ending with unused values
	for i := pos2 + 1; i < chrm1.Len(); i++ {
		res.Raw[i] = fnd.nextUnused(chrm2)
	}
	return res
}

func uniformOrderCrossOver(chrm1, chrm2 gene.Chromosome, mask0 []int, mask1 []int) gene.Chromosome {
	res := chrm1.New()
	fnd := newFinder()

	// mask1: copy values and add to uniq
	for _, idx := range mask1 {
		value := chrm1.Raw[idx]
		res.Raw[idx] = value
		fnd.useValue(value)
	}

	// mask0: value has to be found in unused values
	for _, idx := range mask0 {
		res.Raw[idx] = fnd.nextUnused(chrm2)
	}
	return res
}

// partiallyMatchCrossOver applies a PMX on chrm1 using [pos1 ; pos2[ from chrm2
// example
// chrm1: 1 2 3 4 5 6 7 8 (dst)
// chrm2: 3 7 5 1 6 8 2 4 (src)
// pos  :       x x x
//
// So each choosen base from src to dst:
// - 1 => 4
// - 6 => 5
// - 8 => 6 => 5
//
// Begin of chrm1:           1 2 3 => 4 2 3
// Then copy the src => dst: 4 5 6 => 1 6 8
// End of chrm1:             7 8   => 7 5
// res1: 4 2 3 1 6 8 7 5
func partiallyMatchCrossOver(chrm1, chrm2 gene.Chromosome, pos1, pos2 int) gene.Chromosome {
	// Ref on src and dest using given positions
	src := chrm2.Raw[pos1:pos2]
	dst := chrm1.Raw[pos1:pos2]
	res1 := chrm1.New()

	// First part: apply pmx
	for i := range pos1 {
		res1.Raw[i] = pmxConvert(chrm1.Raw[i], src, dst)
	}

	// Middle part: copy the source (chrm2) into result (res1)
	copy(res1.Raw[pos1:pos2], src)

	// Last part: also apply pmx
	for i := pos2; i < chrm1.Len(); i++ {
		res1.Raw[i] = pmxConvert(chrm1.Raw[i], src, dst)
	}

	return res1
}

// pmxConvert checks if the given value is found in src.
// - if not found: returns the value
// - if found: get the value at the same position in dst
//   - relaunch pmx-convert with the new value until no convertion is found
func pmxConvert(value gene.B, src, dst []gene.B) gene.B {
	// Find value position in src
	idx := indexOf(src, value)
	if idx == -1 { // not found, no convertion
		return value
	}

	// Src index found, convert to dst value
	return pmxConvert(dst[idx], src, dst)
}

// indexOf resturns the first index where the value is found into the input slice ; otherwise, returns -1
func indexOf[T comparable](slc []T, elt T) int {
	for i, item := range slc {
		if item == elt {
			return i
		}
	}
	return -1
}
