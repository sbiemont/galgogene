package operator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelector(t *testing.T) {
	Convey("new", t, func() {
		Convey("when new multi selector", func() {
			Convey("when empty", func() {
				selectors, err := NewMultiSelector([]ProbaSelector{})
				So(err, ShouldBeError, "at least one selector is required")
				So(selectors, ShouldBeEmpty)
			})

			Convey("when missing proba >= 1", func() {
				selectors, err := NewMultiSelector([]ProbaSelector{
					NewProbaSelector(0.1, SelectorRoulette{}),
					NewProbaSelector(0.2, SelectorRoulette{}),
				})
				So(err, ShouldBeError, "selector with proba=1 shall only be the last one")
				So(selectors, ShouldBeEmpty)
			})

			Convey("when proba >= 1 but not last", func() {
				selectors, err := NewMultiSelector([]ProbaSelector{
					NewProbaSelector(0.1, SelectorRoulette{}),
					NewProbaSelector(1.0, SelectorRoulette{}),
					NewProbaSelector(0.2, SelectorRoulette{}),
				})
				So(err, ShouldBeError, "selector with proba=1 shall only be the last one")
				So(selectors, ShouldBeEmpty)
			})

			Convey("when ok", func() {
				selectors, err := NewMultiSelector([]ProbaSelector{
					NewProbaSelector(0.1, SelectorRoulette{}),
					NewProbaSelector(0.2, SelectorRoulette{}),
					NewProbaSelector(1.0, SelectorRoulette{}),
				})
				So(err, ShouldBeNil)
				So(selectors, ShouldHaveLength, 3)
			})
		})
	})
}
