package gene

import (
	"galgogene.git/random"
)

var DefaultMaxValue uint8 = 1

// Bits represents a list of ordered bytes
// * With maxValue = 1, the data list will be 0, 1
// * with maxValue = 255, the data list will be 0, 1, .., 254, 255
type Bits struct {
	Raw      []uint8 // The raw data list
	MaxValue uint8   // The max value to be applied on each byte
}

// NewBits returns a full 0 initialized set of bits
func NewBits(size int, maxValue uint8) Bits {
	return Bits{
		Raw:      make([]uint8, size),
		MaxValue: maxValue,
	}
}

// NewBitsFrom creates a new empty set of bits from another set
func NewBitsFrom(bits Bits) Bits {
	return NewBits(bits.Len(), bits.MaxValue)
}

// NewBitsRandom returns a randomly initialized set of bits
func NewBitsRandom(size int, maxValue uint8) Bits {
	result := Bits{
		Raw:      make([]uint8, size),
		MaxValue: maxValue,
	}
	for i := 0; i < size; i++ {
		result.Raw[i] = result.modulo(random.Byte())
	}
	return result
}

// Len returns the data length
func (bits Bits) Len() int {
	return len(bits.Raw)
}

// Clone returns a full copy of Bits
func (bits Bits) Clone() Bits {
	clone := make([]uint8, bits.Len())
	copy(clone, bits.Raw)
	return Bits{
		Raw:      clone,
		MaxValue: bits.MaxValue,
	}
}

// Rand generates a random byte using the given max value
func (bits Bits) Rand() uint8 {
	return bits.modulo(random.Byte())
}

// Invert each bits of the recorded byte
func (bits Bits) Invert(i int) uint8 {
	return bits.modulo(^bits.Raw[i])
}

func (bits Bits) modulo(value uint8) uint8 {
	if bits.MaxValue == 255 {
		return value
	}

	return value % (bits.MaxValue + 1)
}
