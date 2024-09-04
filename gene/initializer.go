package gene

import (
	"errors"
	"math"

	"github.com/sbiemont/galgogene/random"
)

// Initializer is in charge of the individuals code initialization
type Initializer interface {
	// Check if parameters are valid
	Check(bitsSize int) error

	// Init the individual code using the input parameters
	Init(bitsSize int) Bits
}

// ------------------------------

// RandomInitializer is a full random bits initializer
type RandomInitializer struct {
	MaxValue uint8
}

// NewRandomInitializer builds a new instance
func NewRandomInitializer(maxValue uint8) RandomInitializer {
	return RandomInitializer{
		MaxValue: maxValue,
	}
}

func (izr RandomInitializer) Check(_ int) error {
	if izr.MaxValue == 0 {
		return errors.New("initializer max value cannot be 0")
	}
	return nil
}

func (izr RandomInitializer) Init(bitsSize int) Bits {
	return NewBitsRandom(bitsSize, izr.MaxValue)
}

// ------------------------------

// PermutationInitializer builds a list of shuffled permutations
type PermutationInitializer struct{}

func (PermutationInitializer) Check(bitsSize int) error {
	if bitsSize > math.MaxUint8 {
		return errors.New("bitsSize too big")
	}
	return nil
}

func (PermutationInitializer) Init(bitsSize int) Bits {
	result := NewBits(bitsSize, uint8(bitsSize))
	for i, value := range random.Perm(bitsSize) {
		result.Raw[i] = uint8(value)
	}
	return result
}
