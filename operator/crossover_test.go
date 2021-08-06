package operator

import (
	"testing"

	"genalgo.git/gene"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCrossOvers(t *testing.T) {
	Convey("cross over", t, func() {
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

	Convey("uniform cross over", t, func() {
		bits1 := newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1})
		bits2 := newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0})

		Convey("when rate 0.5", func() {
			res1, res2 := uniformCrossOver(bits1, bits2, 0.5)
			So(res1, ShouldNotResemble, newBits([]uint8{1, 1, 1, 1, 1, 1, 1, 1}))
			So(res2, ShouldNotResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0}))
		})
	})
}
