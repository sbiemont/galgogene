package gene

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitializer(t *testing.T) {
	Convey("initializer", t, func() {
		Convey("random", func() {
			rand.New(rand.NewSource(424242))
			initializer := RandomInitializer{
				MaxValue: 8,
			}
			So(initializer.Check(8), ShouldBeNil)
			So(initializer.Init(8), ShouldResemble, Bits{
				Raw:      []uint8{2, 3, 1, 5, 8, 4, 2, 7},
				MaxValue: 8,
			})
		})

		Convey("permutation", func() {
			Convey("when ok", func() {
				rand.New(rand.NewSource(424242))
				initializer := PermuationInitializer{}
				So(initializer.Check(8), ShouldBeNil)
				So(initializer.Init(8), ShouldResemble, Bits{
					Raw:      []uint8{5, 7, 3, 6, 0, 1, 2, 4},
					MaxValue: 8,
				})
			})

			Convey("when error", func() {
				initializer := PermuationInitializer{}
				So(initializer.Check(256), ShouldNotBeNil)
			})
		})
	})
}
