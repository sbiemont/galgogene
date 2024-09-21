package gene

import (
	"github.com/sbiemont/galgogene/random"
)

// Genetic base (can be binary or uint8)
type B uint8

// Byte cast internal data
func (b B) Byte() byte {
	return byte(b)
}

// Chromosome represents a list of ordered bytes
// * With maxValue = 1, the data list will be 0, 1
// * with maxValue = 255, the data list will be 0, 1, .., 254, 255
type Chromosome struct {
	Raw      []B // The raw data list
	maxValue B   // The max value to be applied on each byte (only useful for random operators)
}

// NewChromosome returns a full 0 initialized set of bases
func NewChromosome(size int, maxValue B) Chromosome {
	return Chromosome{
		Raw:      make([]B, size),
		maxValue: maxValue,
	}
}

// NewChromosomeRandom returns a randomly initialized set of bases
func NewChromosomeRandom(size int, maxValue B) Chromosome {
	result := NewChromosome(size, maxValue)
	for i := range size {
		result.Raw[i] = result.Rand()
	}
	return result
}

// Len returns the data length
func (chrm Chromosome) Len() int {
	return len(chrm.Raw)
}

// Clone returns a full copy of the current chromosome
func (chrm Chromosome) Clone() Chromosome {
	clone := chrm.New()
	copy(clone.Raw, chrm.Raw)
	return clone
}

// New returns a new empty chromosome based on the current properties
func (chrm Chromosome) New() Chromosome {
	return NewChromosome(chrm.Len(), chrm.maxValue)
}

// Rand generates a random base using the given max value
func (chrm Chromosome) Rand() B {
	// Uppercast, compute, downcast
	value := random.Uint64()
	return B(value % (uint64(chrm.maxValue) + 1))
}

// String exports the chromsome as a string
func (chrm Chromosome) String() string {
	res := make([]byte, chrm.Len())
	for i, it := range chrm.Raw {
		res[i] = it.Byte()
	}
	return string(res)
}
