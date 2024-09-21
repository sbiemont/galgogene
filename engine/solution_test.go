package engine

import (
	"testing"

	"github.com/sbiemont/galgogene/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSolution(t *testing.T) {
	Convey("solution", t, func() {
		Convey("when check type", func() {
			sol := Solution{
				Termination: &operator.ImprovementTermination{},
			}
			So(sol.TerminationType(&operator.ImprovementTermination{}), ShouldBeTrue)
			So(sol.TerminationType(&operator.DurationTermination{}), ShouldBeFalse)
		})
	})
}
