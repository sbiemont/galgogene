package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCrossOvers(t *testing.T) {
	Convey("crossover", t, func() {
		bits1 := newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1})
		bits2 := newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0})

		Convey("when 1 split", func() {
			res1, res2 := crossOver(bits1, bits2, []int{4})
			So(res1, ShouldResemble, newBits([]uint8{1, 1, 1, 1, 0, 0, 0, 0}))
			So(res2, ShouldResemble, newBits([]uint8{0, 0, 0, 0, 1, 1, 1, 1}))
			So(bits1, ShouldResemble, newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1})) // unchanged
			So(bits2, ShouldResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0})) // unchanged
		})

		Convey("when 2 splits", func() {
			res1, res2 := crossOver(bits1, bits2, []int{3, 6})
			So(res1, ShouldResemble, newBits([]uint8{1, 1, 1, 0, 0, 0, 1, 1}))
			So(res2, ShouldResemble, newBits([]uint8{0, 0, 0, 1, 1, 1, 0, 0}))
		})

		Convey("when 3 splits", func() {
			res1, res2 := crossOver(bits1, bits2, []int{2, 4, 6})
			So(res1, ShouldResemble, newBits([]uint8{1, 1, 0, 0, 1, 1, 0, 0}))
			So(res2, ShouldResemble, newBits([]uint8{0, 0, 1, 1, 0, 0, 1, 1}))
		})

		Convey("when 3 splits with full bytes", func() {
			bits1 := gene.Bits{
				Raw:      []uint8{0, 5, 10, 15, 20, 25, 30, 35},
				MaxValue: 255,
			}
			bits2 := gene.Bits{
				Raw:      []uint8{220, 225, 230, 235, 240, 245, 250, 255},
				MaxValue: 255,
			}
			res1, res2 := crossOver(bits1, bits2, []int{2, 4, 6})
			So(res1, ShouldResemble, gene.Bits{
				Raw:      []uint8{0, 5, 230, 235, 20, 25, 250, 255},
				MaxValue: 255,
			})
			So(res2, ShouldResemble, gene.Bits{
				Raw:      []uint8{220, 225, 10, 15, 240, 245, 30, 35},
				MaxValue: 255,
			})
		})

		Convey("when same indexes", func() {
			res1, res2 := crossOver(bits1, bits2, []int{7, 7, 7})
			So(res1, ShouldResemble, newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 0}))
			So(res2, ShouldResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 1}))
		})
	})

	Convey("uniform crossover", t, func() {
		bits1 := newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1})
		bits2 := newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0})

		Convey("when rate 0.5", func() {
			res1, res2 := uniformCrossOver(bits1, bits2, 0.5)
			So(res1, ShouldNotResemble, newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1}))
			So(res2, ShouldNotResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0}))
		})
	})

	Convey("davis' order crossover", t, func() {
		bits1 := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9})
		bits2 := newBits([]uint8{9, 8, 7, 6, 4, 5, 4, 3, 2, 1})

		res0 := davisOrderCrossOver(bits1, bits2, 0, 0)
		res1 := davisOrderCrossOver(bits1, bits2, 0, 1)
		res2 := davisOrderCrossOver(bits1, bits2, 0, 4)
		res3 := davisOrderCrossOver(bits1, bits2, 2, 5)
		res4 := davisOrderCrossOver(bits1, bits2, 7, 8)
		res5 := davisOrderCrossOver(bits1, bits2, 8, 8)
		So(res0, ShouldResemble, newBits([]uint8{1, 9, 8, 7, 6, 4, 5, 3, 2}))
		So(res1, ShouldResemble, newBits([]uint8{1, 2, 9, 8, 7, 6, 4, 5, 3}))
		So(res2, ShouldResemble, newBits([]uint8{1, 2, 3, 4, 5, 9, 8, 7, 6}))
		So(res3, ShouldResemble, newBits([]uint8{9, 8, 3, 4, 5, 6, 7, 2, 1}))
		So(res4, ShouldResemble, newBits([]uint8{7, 6, 4, 5, 3, 2, 1, 8, 9}))
		So(res5, ShouldResemble, newBits([]uint8{8, 7, 6, 4, 5, 3, 2, 1, 9}))
	})

	Convey("uniform order crossover", t, func() {
		bits1 := newBits([]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9})
		bits2 := newBits([]uint8{3, 4, 7, 2, 8, 9, 1, 6, 5})

		mask0 := []int{2, 4, 5, 6, 8}
		mask1 := []int{0, 1, 3, 7}
		So(uniformOrderCrossOver(bits1, bits2, mask0, mask1).Raw, ShouldResemble, []uint8{1, 2, 3, 4, 7, 9, 6, 8, 5})
		So(uniformOrderCrossOver(bits2, bits1, mask0, mask1).Raw, ShouldResemble, []uint8{3, 4, 1, 2, 5, 7, 8, 6, 9})
	})
}
