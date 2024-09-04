package random

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRandom(t *testing.T) {
	Convey("random", t, func() {
		Convey("ints", func() {
			result := OrderedInts(10, 20, 4)
			sort.Slice(result, func(i, j int) bool {
				return result[i] < result[j]
			})
			So(result[0], ShouldBeBetweenOrEqual, 10, 20)
			So(result[1], ShouldBeBetweenOrEqual, result[0], 20)
			So(result[2], ShouldBeBetweenOrEqual, result[1], 20)
			So(result[3], ShouldBeBetweenOrEqual, result[2], 20)
		})

		Convey("byte", func() {
			So(Byte(), ShouldBeBetweenOrEqual, 0, 255)
		})
	})
}
