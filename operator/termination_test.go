package operator

import (
	"testing"
	"time"

	"github.com/sbiemont/galgogene/gene"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTerminations(t *testing.T) {
	Convey("terminations", t, func() {
		Convey("when generation Termination", func() {
			Convey("when ko", func() {
				termination := GenerationTermination{K: 10}
				pop := gene.Population{Stats: gene.PopulationStats{GenerationNb: 9}}
				So(termination.End(pop, pop, pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				termination := GenerationTermination{K: 10}
				pop := gene.Population{Stats: gene.PopulationStats{GenerationNb: 10}}
				So(termination.End(pop, pop, pop), ShouldEqual, &termination)
			})
		})

		Convey("when improvement termination", func() {
			Convey("when no k defined", func() {
				termination := ImprovementTermination{}
				pop := gene.Population{Stats: gene.PopulationStats{TotalFitness: 42}}
				So(termination.End(pop, pop, pop), ShouldBeNil)

				pop.Stats.TotalFitness = 43
				So(termination.End(pop, pop, pop), ShouldBeNil)

				// TotalFitness is still 43 => no more improvment
				So(termination.End(pop, pop, pop), ShouldEqual, &termination)
			})

			Convey("when k defined", func() {
				termination := ImprovementTermination{
					K: 3,
				}
				pop := gene.Population{Stats: gene.PopulationStats{TotalFitness: 42}}
				So(termination.End(pop, pop, pop), ShouldBeNil)               // 0 -> 42
				So(termination.End(pop, pop, pop), ShouldBeNil)               // k: 1
				So(termination.End(pop, pop, pop), ShouldBeNil)               // k: 2
				So(termination.End(pop, pop, pop), ShouldEqual, &termination) // k: 3
			})
		})

		Convey("when above fitness termination", func() {
			Convey("when ko", func() {
				pop := gene.Population{Stats: gene.PopulationStats{Elite: gene.Individual{Fitness: 0.6}}}
				termination := FitnessTermination{Fitness: 0.8}
				So(termination.End(pop, pop, pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				pop := gene.Population{Stats: gene.PopulationStats{Elite: gene.Individual{Fitness: 0.7}}}
				termination := FitnessTermination{Fitness: 0.7}
				So(termination.End(pop, pop, pop), ShouldEqual, &termination)
			})
		})

		Convey("when duration termination", func() {
			Convey("when ko", func() {
				termination := DurationTermination{Duration: time.Minute}
				pop := gene.Population{Stats: gene.PopulationStats{TotalDuration: time.Second}}
				So(termination.End(pop, pop, pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				termination := DurationTermination{Duration: time.Minute}
				pop := gene.Population{Stats: gene.PopulationStats{TotalDuration: time.Minute}}
				So(termination.End(pop, pop, pop), ShouldEqual, &termination)
			})
		})

		Convey("when multi termination", func() {
			termination1 := GenerationTermination{K: 10}
			termination2 := DurationTermination{Duration: time.Minute}
			termination := MultiTermination{&termination1, &termination2}

			Convey("when ko", func() {
				pop := gene.Population{
					Stats: gene.PopulationStats{
						GenerationNb:  9,
						TotalDuration: time.Second,
					},
				}
				So(termination.End(pop, pop, pop), ShouldBeNil)
			})

			Convey("when ok, termination #1", func() {
				pop := gene.Population{
					Stats: gene.PopulationStats{
						GenerationNb:  10,
						TotalDuration: time.Second,
					},
				}
				So(termination.End(pop, pop, pop), ShouldEqual, &termination1)
			})

			Convey("when ok, termination #2", func() {
				pop := gene.Population{
					Stats: gene.PopulationStats{
						GenerationNb:  9,
						TotalDuration: time.Minute,
					},
				}
				So(termination.End(pop, pop, pop), ShouldEqual, &termination2)
			})
		})
	})
}
