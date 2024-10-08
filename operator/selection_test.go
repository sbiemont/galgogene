package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSelection(t *testing.T) {
	Convey("selection", t, func() {
		pop1 := func() gene.Population {
			return gene.Population{
				Individuals: []gene.Individual{
					{Fitness: 0.1, Rank: 1},
					{Fitness: 0.5, Rank: 5},
					{Fitness: 0.6, Rank: 6},
					{Fitness: 0.9, Rank: 9},
				},
			}
		}

		Convey("when roulette", func() {
			Convey("when total fitness = 0", func() {
				p1 := pop1()
				p1.Stats.TotalFitness = 0
				ind, err := RouletteSelection{}.Select(p1)
				So(err, ShouldBeNil)
				So(ind, ShouldResemble, gene.Individual{Fitness: 0.1, Rank: 1})
			})

			Convey("when total fitness = 0.5", func() {
				p1 := pop1()
				p1.Stats.TotalFitness = 0.5
				ind, err := RouletteSelection{}.Select(p1)
				So(err, ShouldBeNil)
				if ind.Fitness == 0.1 {
					So(ind, ShouldResemble, gene.Individual{Fitness: 0.1, Rank: 1})
				} else {
					So(ind, ShouldResemble, gene.Individual{Fitness: 0.5, Rank: 5})
				}
			})
		})

		Convey("when tournament", func() {
			pop := pop1()
			Convey("when k=0", func() {
				ind, err := TournamentSelection{Fighters: 0}.Select(pop)
				So(err, ShouldNotBeNil)
				So(ind, ShouldResemble, gene.Individual{})
			})

			Convey("when k=1", func() {
				random.Seed(1)
				ind, err := TournamentSelection{Fighters: 1}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldEqual, 0.9)
			})

			Convey("when k=2", func() {
				random.Seed(3)
				ind, err := TournamentSelection{Fighters: 2}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldEqual, 0.5)
			})

			Convey("when k=3", func() {
				random.Seed(3)
				ind, err := TournamentSelection{Fighters: 3}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldEqual, 0.6)
			})

			Convey("when k=4", func() {
				random.Seed(3)
				ind, err := TournamentSelection{Fighters: 4}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldEqual, 0.6)
			})
		})

		Convey("when multi selection", func() {
			Convey("when ok", func() {
				selections := MultiSelection{}.
					Use(0.1, RouletteSelection{}).
					Use(0.2, EliteSelection{}).
					Otherwise(TournamentSelection{Fighters: 42})

				So(selections, ShouldResemble, multiSelection{
					selections: []probaSelection{
						{
							rate: 0.1,
							sel:  RouletteSelection{},
						},
						{
							rate: 0.2,
							sel:  EliteSelection{},
						},
					},
					deflt: TournamentSelection{Fighters: 42},
				})
			})
		})
	})
}
