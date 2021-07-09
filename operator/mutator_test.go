package operator

import (
	"testing"

	"genalgo.git/gene"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMutations(t *testing.T) {
	Convey("cross over", t, func() {
		bits1 := gene.Bits{1, 1, 1, 1, 1, 1, 1, 1}
		bits2 := gene.Bits{0, 0, 0, 0, 0, 0, 0, 0}

		Convey("when 1 split", func() {
			res1, res2 := crossOver(bits1, bits2, []int{4})
			So(res1, ShouldResemble, gene.Bits{1, 1, 1, 1, 0, 0, 0, 0})
			So(res2, ShouldResemble, gene.Bits{0, 0, 0, 0, 1, 1, 1, 1})
		})

		Convey("when 2 splits", func() {
			res1, res2 := crossOver(bits1, bits2, []int{3, 6})
			So(res1, ShouldResemble, gene.Bits{1, 1, 1, 0, 0, 0, 1, 1})
			So(res2, ShouldResemble, gene.Bits{0, 0, 0, 1, 1, 1, 0, 0})
		})

		Convey("when 3 splits", func() {
			res1, res2 := crossOver(bits1, bits2, []int{2, 4, 6})
			So(res1, ShouldResemble, gene.Bits{1, 1, 0, 0, 1, 1, 0, 0})
			So(res2, ShouldResemble, gene.Bits{0, 0, 1, 1, 0, 0, 1, 1})
		})

		Convey("when same indexes", func() {
			res1, res2 := crossOver(bits1, bits2, []int{7, 7, 7})
			So(res1, ShouldResemble, gene.Bits{1, 1, 1, 1, 1, 1, 1, 0})
			So(res2, ShouldResemble, gene.Bits{0, 0, 0, 0, 0, 0, 0, 1})
		})
	})

	Convey("uniform cross over", t, func() {
		bits1 := gene.Bits{1, 1, 1, 1, 1, 1, 1, 1}
		bits2 := gene.Bits{0, 0, 0, 0, 0, 0, 0, 0}

		Convey("when rate 0.5", func() {
			res1, res2 := uniformCrossOver(bits1, bits2, 0.5)
			So(res1, ShouldNotResemble, gene.Bits{1, 1, 1, 1, 1, 1, 1, 1})
			So(res2, ShouldNotResemble, gene.Bits{0, 0, 0, 0, 0, 0, 0, 0})
		})
	})

	Convey("mutate", t, func() {
		bits1 := gene.Bits{1, 1, 1, 1, 1, 1, 1, 1}

		Convey("when none mutated", func() {
			res1 := mutate(bits1, 0.0)
			So(res1, ShouldResemble, gene.Bits{1, 1, 1, 1, 1, 1, 1, 1})
		})

		Convey("when some mutated", func() {
			res1 := mutate(bits1, 0.5)
			So(res1, ShouldNotResemble, gene.Bits{1, 1, 1, 1, 1, 1, 1, 1})
			So(res1, ShouldNotResemble, gene.Bits{0, 0, 0, 0, 0, 0, 0, 0})
		})

		Convey("when all mutated", func() {
			res1 := mutate(bits1, 1.0)
			So(res1, ShouldResemble, gene.Bits{0, 0, 0, 0, 0, 0, 0, 0})
		})
	})
}
