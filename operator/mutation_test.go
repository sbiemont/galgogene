package operator

import (
	"math/rand"
	"testing"

	"github.com/sbiemont/galgogene/gene"

	. "github.com/smartystreets/goconvey/convey"
)

func newBits(bits []uint8) gene.Bits {
	return gene.Bits{
		Raw:      bits,
		MaxValue: 8,
	}
}

func TestMutations(t *testing.T) {
	Convey("mutation position", t, func() {
		bits1 := newBits(make([]uint8, 8)) // empty 8 bits gene

		Convey("when middle", func() {
			rand.New(rand.NewSource(424242))
			pos1, pos2 := mutationPositions(bits1)
			So(pos1, ShouldEqual, 2)
			So(pos2, ShouldEqual, 6)
		})

		Convey("when first", func() {
			rand.New(rand.NewSource(42))
			pos1, pos2 := mutationPositions(bits1)
			So(pos1, ShouldEqual, 1)
			So(pos2, ShouldEqual, 3)
		})
	})

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

	Convey("bit flip mutation", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := UniqueMutation{}

		Convey("when middle", func() {
			rand.New(rand.NewSource(424242)) // pos: 6
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 3, 8})
		})

		Convey("when first", func() {
			rand.New(rand.NewSource(42)) // pos: 1
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 3, 3, 4, 5, 6, 7, 8})
		})
	})

	Convey("swap permutation", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := SwapPermutation{}

		Convey("when middle", func() {
			rand.New(rand.NewSource(424242))
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 2, 7, 4, 5, 6, 3, 8})
		})

		Convey("when first", func() {
			rand.New(rand.NewSource(42))
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 4, 3, 2, 5, 6, 7, 8})
		})
	})

	Convey("inversion permutation", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := InversionPermutation{}

		Convey("when middle", func() {
			rand.New(rand.NewSource(424242))
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 2, 7, 6, 5, 4, 3, 8})
		})

		Convey("when first", func() {
			rand.New(rand.NewSource(42))
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 4, 3, 2, 5, 6, 7, 8})
		})
	})

	Convey("scramble permutation", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := SramblePermutation{}

		Convey("when middle", func() {
			rand.New(rand.NewSource(424242))
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 2, 6, 3, 4, 5, 7, 8})
		})

		Convey("when first", func() {
			rand.New(rand.NewSource(42))
			result := mutation.Mutate(bits)
			So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []uint8{1, 3, 2, 4, 5, 6, 7, 8})
		})
	})
}
