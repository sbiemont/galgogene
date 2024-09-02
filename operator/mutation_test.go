package operator

import (
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

func countUnique[T comparable](values []T) int {
	uniq := make(map[T]struct{})
	for _, val := range values {
		uniq[val] = struct{}{}
	}
	return len(uniq)
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

	Convey("swap permutation", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := SwapPermutation{}

		result := mutation.Mutate(bits)
		So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
		So(result.Raw, ShouldNotResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
		So(countUnique(result.Raw), ShouldEqual, 8)
	})

	Convey("inversion permutation", t, func() {
		bits := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := func(in gene.Bits, out *gene.Bits, _, _ int) {
			// Force positions
			pos1 := 2
			pos2 := 6
			for i := pos1; i <= pos2; i++ {
				out.Raw[i] = in.Raw[pos2-i+pos1]
			}
		}

		result := permutation(bits, mutation)
		So(bits.Raw, ShouldResemble, []uint8{1, 2, 3, 4, 5, 6, 7, 8})
		So(result.Raw, ShouldResemble, []uint8{1, 2, 7, 6, 5, 4, 3, 8})
	})
}
