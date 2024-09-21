package gene

import (
	"fmt"

	"github.com/sbiemont/galgogene/random"
)

// Initializer is in charge of the individuals code initialization
type Initializer interface {
	// Init the individual code using the input parameters
	Init(chrmSize int) (Chromosome, error)
}

// ------------------------------

// RandomInitializer is a full random chromosome initializer
type RandomInitializer struct {
	MaxValue B
}

// NewRandomInitializer builds a new instance
func NewRandomInitializer(maxValue B) RandomInitializer {
	return RandomInitializer{
		MaxValue: maxValue,
	}
}

func (izr RandomInitializer) Init(chrmSize int) (Chromosome, error) {
	if izr.MaxValue == 0 {
		return Chromosome{}, fmt.Errorf("initializer max value cannot be 0")
	}

	return NewChromosomeRandom(chrmSize, izr.MaxValue), nil
}

// ------------------------------

// PermutationInitializer builds a list of shuffled permutations
type PermutationInitializer struct{}

func (PermutationInitializer) Init(chrmSize int) (Chromosome, error) {
	if chrmSize == 0 {
		return Chromosome{}, fmt.Errorf("chrmSize cannot be 0")
	}

	result := NewChromosome(chrmSize, B(chrmSize))
	for i, value := range random.Perm(chrmSize) {
		result.Raw[i] = B(value)
	}
	return result, nil
}
