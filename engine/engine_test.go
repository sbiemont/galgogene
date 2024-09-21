package engine

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEngine(t *testing.T) {
	Convey("check", t, func() {
		Convey("when missing fitness", func() {
			So(Engine{}.check(), ShouldBeError, "fitness must be set")
		})

		Convey("when missing initializer", func() {
			eng := Engine{
				Fitness: func(c gene.Chromosome) float64 { return 0 },
			}
			So(eng.check(), ShouldBeError, "initializer must be set")
		})

		Convey("when minimalist", func() {
			eng := Engine{
				Initializer: gene.RandomInitializer{MaxValue: 1},
				Selection:   operator.RouletteSelection{},
				CrossOver:   operator.OnePointCrossOver{},
				Survivor:    operator.RankSurvivor{},
				Termination: &operator.DurationTermination{},
				Fitness:     func(c gene.Chromosome) float64 { return 0 },
			}
			So(eng.check(), ShouldBeNil)
		})
	})
}
