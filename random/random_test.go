package random

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRandom(t *testing.T) {
	Convey("random", t, func() {
		Convey("ints", func() {
			Seed(42)
			result := OrderedInts(10, 20, 4)
			sort.Slice(result, func(i, j int) bool {
				return result[i] < result[j]
			})
			So(result, ShouldResemble, []int{13, 15, 16, 16})
		})
	})
}
