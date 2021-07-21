package gene

import (
	"errors"
	"fmt"

	"genalgo.git/random"
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

// NewBitsFromBytes builds a list of bits from an array of bytes
func NewBitsFromBytes(bytes []byte) Bits {
	result := Bits{
		Raw:      make([]uint8, 8*len(bytes)),
		MaxValue: DefaultMaxValue,
	}
	var i int
	for _, b := range bytes {
		result.Raw[i] = b & 0x80 >> 7
		result.Raw[i+1] = b & 0x40 >> 6
		result.Raw[i+2] = b & 0x20 >> 5
		result.Raw[i+3] = b & 0x10 >> 4
		result.Raw[i+4] = b & 0x08 >> 3
		result.Raw[i+5] = b & 0x04 >> 2
		result.Raw[i+6] = b & 0x02 >> 1
		result.Raw[i+7] = b & 0x01
		i += 8
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

// GroupBitsBy fetches bits by group and builds uint8 values
// Unfield bits are let at 0
//  eg.:
//  * GroupBitsBy(2) groups every 2 bits into 1 uint8 (0b000000xx)
//  * GroupBitsBy(4) groups every 4 bits into 1 uint8 (0b0000xxxx)
func (bits Bits) GroupBitsBy(n int) ([]uint8, error) {
	if bits.Len()%n != 0 {
		return nil, fmt.Errorf("cannot group, total nb of bits (%d) should be modulo %d", bits.Len(), n)
	}
	if n > 8 {
		return nil, errors.New("cannot group, n > 8")
	}

	// Group
	result := make([]uint8, bits.Len()/n)
	r := 0
	for i := 0; i < bits.Len(); i += n {
		var c uint8
		for k := 0; k < n; k++ {
			c |= bits.Raw[i+k] << (n - k - 1)
		}
		result[r] = c
		r++
	}
	return result, nil
}
