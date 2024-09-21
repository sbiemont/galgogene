package gene

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func countUnique[T comparable](values []T) int {
	uniq := make(map[T]struct{})
	for _, val := range values {
		uniq[val] = struct{}{}
	}
	return len(uniq)
}

func TestInitializer(t *testing.T) {
	Convey("initializer", t, func() {
		Convey("random", func() {
			initializer := RandomInitializer{
				MaxValue: 8,
			}
			chrm, err := initializer.Init(8)
			So(err, ShouldBeNil)
			So(countUnique(chrm.Raw), ShouldBeGreaterThan, 0)
		})

		Convey("permutation", func() {
			Convey("when ok", func() {
				initializer := PermutationInitializer{}
				chrm, err := initializer.Init(8)
				So(err, ShouldBeNil)
				So(countUnique(chrm.Raw), ShouldEqual, 8)
			})

			Convey("when error", func() {
				initializer := PermutationInitializer{}
				_, err := initializer.Init(0)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
