package random

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRanom(t *testing.T) {
	Convey("random", t, func() {
		Convey("ints", func() {
			result := Ints(10, 20, 4)
			So(result[0], ShouldBeBetweenOrEqual, 10, 20)
			So(result[1], ShouldBeBetweenOrEqual, result[0], 20)
			So(result[2], ShouldBeBetweenOrEqual, result[1], 20)
			So(result[3], ShouldBeBetweenOrEqual, result[2], 20)
		})

		Convey("byte", func() {
			Convey("when min", func() {
				rand.Seed(273)
				So(Byte(), ShouldEqual, 0)
			})

			Convey("when random", func() {
				rand.Seed(1)
				So(Byte(), ShouldEqual, 33)
			})

			Convey("when max", func() {
				rand.Seed(74)
				So(Byte(), ShouldEqual, 255)
			})
		})
	})
}
