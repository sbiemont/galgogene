package gene

import (
	"errors"
	"fmt"
)

type Serializer struct {
	pack int // Group by n bits
}

func NewSerializer(n int) (Serializer, error) {
	if n > 8 {
		return Serializer{}, errors.New("cannot group, n > 8")
	}

	return Serializer{
		pack: n,
	}, nil
}

// GroupBitsBy fetches bits by group and builds uint8 values
// Unfield bits are let at 0
//  eg.:
//  * GroupBitsBy(2) groups every 2 bits into 1 uint8 (0b000000xx)
//  * GroupBitsBy(4) groups every 4 bits into 1 uint8 (0b0000xxxx)
func (szr Serializer) ToBytes(bits Bits) ([]uint8, error) {
	if bits.Len()%szr.pack != 0 {
		return nil, fmt.Errorf("cannot group, total nb of bits (%d) should be modulo %d", bits.Len(), szr.pack)
	}

	// Group
	result := make([]uint8, bits.Len()/szr.pack)
	r := 0
	for i := 0; i < bits.Len(); i += szr.pack {
		var c uint8
		for k := 0; k < szr.pack; k++ {
			c |= bits.Raw[i+k] << (szr.pack - k - 1)
		}
		result[r] = c
		r++
	}
	return result, nil
}

// NewBitsFromBytes builds a list of bits from an array of bytes
//  * ToBits(n = 1) split bytes every 1 bit into 8 uint8 (0b0000000x)
//  * ToBits(n = 2) split bytes every 2 bits into 4 uint8 (0b000000xx)
//  * ToBits(n = 3) split bytes every 3 bits into 4 uint8 (0b00000xxx)
func (szr Serializer) ToBits(bytes []byte) Bits {
	result := Bits{
		Raw:      make([]uint8, szr.pack*len(bytes)),
		MaxValue: DefaultMaxValue,
	}
	var i int
	for _, byt := range bytes {
		for j := 0; j < szr.pack; j++ {
			k := szr.pack - j - 1
			result.Raw[i+j] = byt & (0x01 << k) >> k
			_ = k
		}
		i += szr.pack
	}
	return result
}
