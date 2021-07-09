package operator

import (
	"testing"
	"time"

	"genalgo.git/gene"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnders(t *testing.T) {
	Convey("enders", t, func() {
		Convey("when generation ender", func() {
			Convey("when ko", func() {
				ender := EnderGeneration{K: 10}
				pop := gene.Population{GenerationNb: 9}
				So(ender.End(pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				ender := EnderGeneration{K: 10}
				pop := gene.Population{GenerationNb: 10}
				So(ender.End(pop), ShouldEqual, &ender)
			})
		})

		Convey("when improvement ender", func() {
			ender := EnderImprovement{}
			pop := gene.Population{TotalFitness: 42}
			So(ender.End(pop), ShouldBeNil)

			pop.TotalFitness = 43
			So(ender.End(pop), ShouldBeNil)

			// TotalFitness is still 43 => no more improvment
			So(ender.End(pop), ShouldEqual, &ender)
		})

		Convey("when above fitness ender", func() {
			pop := gene.Population{
				Individuals: []gene.Individual{
					{Fitness: 0.5},
					{Fitness: 0.6},
					{Fitness: 0.7},
				},
			}

			Convey("when ko", func() {
				ender := EnderAboveFitness{Fitness: 0.8}
				So(ender.End(pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				ender := EnderAboveFitness{Fitness: 0.7}
				So(ender.End(pop), ShouldEqual, &ender)
			})
		})

		Convey("when below fitness ender", func() {
			pop := gene.Population{
				Individuals: []gene.Individual{
					{Fitness: 0.5},
					{Fitness: 0.6},
					{Fitness: 0.7},
				},
			}

			Convey("when ko", func() {
				ender := EnderBelowFitness{Fitness: 0.4}
				So(ender.End(pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				ender := EnderBelowFitness{Fitness: 0.5}
				So(ender.End(pop), ShouldEqual, &ender)
			})
		})

		Convey("when duration ender", func() {
			Convey("when ko", func() {
				ender := EnderDuration{Duration: time.Minute}
				pop := gene.Population{TotalDuration: time.Second}
				So(ender.End(pop), ShouldBeNil)
			})

			Convey("when ok", func() {
				ender := EnderDuration{Duration: time.Minute}
				pop := gene.Population{TotalDuration: time.Minute}
				So(ender.End(pop), ShouldEqual, &ender)
			})
		})

		Convey("when multi ender", func() {
			ender1 := EnderGeneration{K: 10}
			ender2 := EnderDuration{Duration: time.Minute}
			ender := MultiEnder{&ender1, &ender2}

			Convey("when ko", func() {
				pop := gene.Population{
					GenerationNb:  9,
					TotalDuration: time.Second,
				}
				So(ender.End(pop), ShouldBeNil)
			})

			Convey("when ok, ender #1", func() {
				pop := gene.Population{
					GenerationNb:  10,
					TotalDuration: time.Second,
				}
				So(ender.End(pop), ShouldEqual, &ender1)
			})

			Convey("when ok, ender #2", func() {
				pop := gene.Population{
					GenerationNb:  9,
					TotalDuration: time.Minute,
				}
				So(ender.End(pop), ShouldEqual, &ender2)
			})
		})
	})
}
