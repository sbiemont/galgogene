package gene

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func newIndividual(chrm []B, fitness float64) Individual {
	return Individual{
		Fitness: fitness,
		Code: Chromosome{
			Raw: chrm,
		},
	}
}

func TestPopulation(t *testing.T) {
	Convey("population", t, func() {
		ind1 := newIndividual([]B{1, 1, 1, 1, 1, 1, 1, 1}, 0.1)
		ind2 := newIndividual([]B{0, 0, 0, 0, 0, 0, 0, 0}, 0.2)
		ind3 := newIndividual([]B{1, 0, 0, 0, 0, 0, 0, 0}, 0.3)
		ind4 := newIndividual([]B{1, 1, 0, 0, 0, 0, 0, 0}, 0.4)

		Convey("when first", func() {
			pop := Population{
				Individuals: []Individual{ind1, ind2, ind3, ind4},
			}

			So(pop.First(1).Individuals, ShouldResemble, []Individual{ind1})
			So(pop.First(4).Individuals, ShouldResemble, []Individual{ind1, ind2, ind3, ind4})
		})

		Convey("when last", func() {
			pop := Population{
				Individuals: []Individual{ind1, ind2, ind3, ind4},
			}

			So(pop.Last(1).Individuals, ShouldResemble, []Individual{ind4})
			So(pop.Last(4).Individuals, ShouldResemble, []Individual{ind1, ind2, ind3, ind4})
		})

		Convey("when compute fitness", func() {
			pop := Population{
				Individuals: []Individual{ind1, ind2, ind3, ind4},
			}
			So(pop.Stats.TotalFitness, ShouldEqual, 0)
			pop.ComputeTotalFitness()
			So(pop.Stats, ShouldResemble, PopulationStats{
				TotalFitness: 0.1 + 0.2 + 0.3 + 0.4,
				Elite:        ind4,
			})
		})

		Convey("when sort by rank", func() {
			pop := Population{
				Individuals: []Individual{
					{Rank: 4},
					{Rank: 3},
					{Rank: 2},
					{Rank: 1},
				},
			}

			pop.SortByRank()
			So(pop.Individuals, ShouldResemble, []Individual{
				{Rank: 1},
				{Rank: 2},
				{Rank: 3},
				{Rank: 4},
			})
		})

		Convey("when move rank", func() {
			pop := Population{
				Individuals: []Individual{
					{Rank: 4},
					{Rank: 3},
					{Rank: 2},
					{Rank: 1},
				},
			}

			pop.ComputeRank()
			So(pop.Individuals, ShouldResemble, []Individual{
				{Rank: 5},
				{Rank: 4},
				{Rank: 3},
				{Rank: 2},
			})

		})
	})
}
