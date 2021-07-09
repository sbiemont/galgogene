package gene

import (
	"genalgo.git/random"
)

// Bits represents a list of ordered bit (0 or 1)
type Bits []uint8

// NewBits returns a full 0 initialized set of bits
func NewBits(size int) Bits {
	return make(Bits, size)
}

// NewBitsRandom returns a randomly initialized set of bits
func NewBitsRandom(size int) Bits {
	result := make(Bits, size)
	for i := 0; i < size; i++ {
		result[i] = random.Bit()
	}
	return result
}

// NewBitsFromBytes builds a list of bits from an array of bytes
// Beware: each byte is reverted to build an ordered list of 8 bits
func NewBitsFromBytes(bytes []byte) Bits {
	result := make(Bits, 8*len(bytes))
	var i int
	for _, b := range bytes {
		result[i] = b & 0x01
		result[i+1] = b & 0x02 >> 1
		result[i+2] = b & 0x04 >> 2
		result[i+3] = b & 0x08 >> 3
		result[i+4] = b & 0x10 >> 4
		result[i+5] = b & 0x20 >> 5
		result[i+6] = b & 0x40 >> 6
		result[i+7] = b & 0x80 >> 7
		i += 8
	}
	return result
}

// ToBytes converts bits int bytes
// Beware: each byte is reverted to build an ordered list of 8 bits
// If bits are missing, zeros are added
func (bits Bits) ToBytes() []byte {
	nbBytes := int(len(bits) / 8)

	// Add zeros at the end to get modulo 8
	k := 8 - len(bits)%8
	fullBits := bits
	if k < 8 {
		nbBytes++ // one more complete byte
		fullBits = append(bits, make([]uint8, k)...)
	}

	r := 0
	result := make([]byte, nbBytes)
	for i := 0; i < len(fullBits); i += 8 {
		var c byte
		c |= fullBits[i]
		c |= fullBits[i+1] << 1
		c |= fullBits[i+2] << 2
		c |= fullBits[i+3] << 3
		c |= fullBits[i+4] << 4
		c |= fullBits[i+5] << 5
		c |= fullBits[i+6] << 6
		c |= fullBits[i+7] << 7

		result[r] = c
		r++
	}
	return result
}
