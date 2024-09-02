package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
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
				ind, err := TournamentSelection{Fighters: 1}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldBeGreaterThanOrEqualTo, 0.1)
			})

			Convey("when k=2", func() {
				ind, err := TournamentSelection{Fighters: 2}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldBeGreaterThanOrEqualTo, 0.1)
			})

			Convey("when k=3", func() {
				ind, err := TournamentSelection{Fighters: 3}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldBeGreaterThanOrEqualTo, 0.1)
			})

			Convey("when k=4", func() {
				ind, err := TournamentSelection{Fighters: 4}.Select(pop)
				So(err, ShouldBeNil)
				So(ind.Fitness, ShouldBeGreaterThanOrEqualTo, 0.1)
			})
		})

		// Convey("when new multi selection", func() {
		// 	Convey("when empty", func() {
		// 		selections, err := NewMultiSelection([]ProbaSelection{})
		// 		So(err, ShouldBeError, "at least one selection is required")
		// 		So(selections, ShouldBeEmpty)
		// 	})

		// 	Convey("when missing proba >= 1", func() {
		// 		selections, err := NewMultiSelection([]ProbaSelection{
		// 			NewProbaSelection(0.1, SelectionRoulette{}),
		// 			NewProbaSelection(0.2, SelectionRoulette{}),
		// 		})
		// 		So(err, ShouldBeError, "selection with proba=1 shall only be the last one")
		// 		So(selections, ShouldBeEmpty)
		// 	})

		// 	Convey("when proba >= 1 but not last", func() {
		// 		selections, err := NewMultiSelection([]ProbaSelection{
		// 			NewProbaSelection(0.1, SelectionRoulette{}),
		// 			NewProbaSelection(1.0, SelectionRoulette{}),
		// 			NewProbaSelection(0.2, SelectionRoulette{}),
		// 		})
		// 		So(err, ShouldBeError, "selection with proba=1 shall only be the last one")
		// 		So(selections, ShouldBeEmpty)
		// 	})

		// 	Convey("when ok", func() {
		// 		selections, err := NewMultiSelection([]ProbaSelection{
		// 			NewProbaSelection(0.1, SelectionRoulette{}),
		// 			NewProbaSelection(0.2, SelectionRoulette{}),
		// 			NewProbaSelection(1.0, SelectionRoulette{}),
		// 		})
		// 		So(err, ShouldBeNil)
		// 		So(selections, ShouldHaveLength, 3)
		// 	})
		// })
	})
}
