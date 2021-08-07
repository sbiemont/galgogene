package operator

import (
	"math/rand"
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

	Convey("one swap", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := SwapMutation{}

		Convey("when middle", func() {
			rand.Seed(424242)
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 2, 7, 4, 5, 6, 3, 8})
		})

		Convey("when first", func() {
			rand.Seed(42)
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 4, 3, 2, 5, 6, 7, 8})
		})
	})

	Convey("two swap", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := TwoSwapMutation{}

		Convey("when middle", func() {
			rand.Seed(424242)
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 2, 7, 6, 5, 4, 3, 8})
		})

		Convey("when first", func() {
			rand.Seed(42)
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 4, 3, 2, 5, 6, 7, 8})
		})
	})
}
