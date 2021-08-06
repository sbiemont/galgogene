package operator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelection(t *testing.T) {
	Convey("new", t, func() {
		Convey("when new multi selection", func() {
			Convey("when empty", func() {
				selections, err := NewMultiSelection([]ProbaSelection{})
				So(err, ShouldBeError, "at least one selection is required")
				So(selections, ShouldBeEmpty)
			})

			Convey("when missing proba >= 1", func() {
				selections, err := NewMultiSelection([]ProbaSelection{
					NewProbaSelection(0.1, SelectionRoulette{}),
					NewProbaSelection(0.2, SelectionRoulette{}),
				})
				So(err, ShouldBeError, "selection with proba=1 shall only be the last one")
				So(selections, ShouldBeEmpty)
			})

			Convey("when proba >= 1 but not last", func() {
				selections, err := NewMultiSelection([]ProbaSelection{
					NewProbaSelection(0.1, SelectionRoulette{}),
					NewProbaSelection(1.0, SelectionRoulette{}),
					NewProbaSelection(0.2, SelectionRoulette{}),
				})
				So(err, ShouldBeError, "selection with proba=1 shall only be the last one")
				So(selections, ShouldBeEmpty)
			})

			Convey("when ok", func() {
				selections, err := NewMultiSelection([]ProbaSelection{
					NewProbaSelection(0.1, SelectionRoulette{}),
					NewProbaSelection(0.2, SelectionRoulette{}),
					NewProbaSelection(1.0, SelectionRoulette{}),
				})
				So(err, ShouldBeNil)
				So(selections, ShouldHaveLength, 3)
			})
		})
	})
}
