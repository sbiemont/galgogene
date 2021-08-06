package operator

import (
	"testing"

	"genalgo.git/gene"
	. "github.com/smartystreets/goconvey/convey"
)

func newBits(bits []uint8) gene.Bits {
	return gene.Bits{
		Raw:      bits,
		MaxValue: gene.DefaultMaxValue,
	}
}

func TestMutations(t *testing.T) {
	Convey("mutate", t, func() {
		bits1 := newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1})
		toZero := func(gene.Bits, int) uint8 { return 0 }

		Convey("when none mutated", func() {
			res1 := mutate(bits1, 0.0, toZero)
			So(res1, ShouldResemble, newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1}))
		})

		Convey("when some mutated", func() {
			res1 := mutate(bits1, 0.5, toZero)
			So(res1, ShouldNotResemble, newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1}))
			So(res1, ShouldNotResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0}))
		})

		Convey("when all mutated", func() {
			res1 := mutate(bits1, 1.0, toZero)
			So(res1, ShouldResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0}))
		})
	})
}
