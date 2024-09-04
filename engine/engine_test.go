package engine

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEngineCheck(t *testing.T) {
	Convey("check", t, func() {
		Convey("when empty", func() {
			So(Engine{}.check(), ShouldBeError, "initializer must be set")
		})

		Convey("when minimalist", func() {
			eng := Engine{
				Initializer: gene.RandomInitializer{MaxValue: 1},
				Selection:   operator.RouletteSelection{},
				CrossOver:   operator.OnePointCrossOver{},
				Survivor:    operator.ChildrenSurvivor{},
				Termination: &operator.DurationTermination{},
			}
			So(eng.check(), ShouldBeNil)
		})
	})
}
