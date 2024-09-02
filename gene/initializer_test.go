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
			So(initializer.Check(8), ShouldBeNil)
			So(countUnique(initializer.Init(8).Raw), ShouldBeGreaterThan, 0)
		})

		Convey("permutation", func() {
			Convey("when ok", func() {
				initializer := PermuationInitializer{}
				So(initializer.Check(8), ShouldBeNil)
				So(countUnique(initializer.Init(8).Raw), ShouldEqual, 8)
			})

			Convey("when error", func() {
				initializer := PermuationInitializer{}
				So(initializer.Check(256), ShouldNotBeNil)
			})
		})
	})
}
