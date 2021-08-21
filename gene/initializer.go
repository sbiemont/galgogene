package gene

import (
	"errors"
	"math"
	"math/rand"
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

func (RandomInitializer) Check(_ int) error {
	return nil
}

func (izr RandomInitializer) Init(bitsSize int) Bits {
	return NewBitsRandom(bitsSize, izr.MaxValue)
}

// ------------------------------

// PermuationInitializer builds a list of shuffled permuations
type PermuationInitializer struct{}

func (PermuationInitializer) Check(bitsSize int) error {
	if bitsSize > math.MaxUint8 {
		return errors.New("bitsSize too big")
	}
	return nil
}

func (PermuationInitializer) Init(bitsSize int) Bits {
	result := NewBits(bitsSize, uint8(bitsSize))
	for i, value := range rand.Perm(bitsSize) {
		result.Raw[i] = uint8(value)
	}
	return result
}
