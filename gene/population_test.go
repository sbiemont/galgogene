package gene

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func newIndividual(bits []uint8) Individual {
	return Individual{
		Code: Bits{
			Raw:      bits,
			MaxValue: DefaultMaxValue,
		},
	}
}

func TestPopulation(t *testing.T) {
	Convey("population", t, func() {
		ind1 := newIndividual([]uint8{1, 1, 1, 1, 1, 1, 1, 1})
		ind2 := newIndividual([]uint8{0, 0, 0, 0, 0, 0, 0, 0})
		ind3 := newIndividual([]uint8{1, 0, 0, 0, 0, 0, 0, 0})
		ind4 := newIndividual([]uint8{1, 1, 0, 0, 0, 0, 0, 0})

		fitness := func(b Bits) float64 {
			var sum float64
			for _, bit := range b.Raw {
				sum += float64(bit)
			}
			return sum
		}

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
				fitness:     fitness,
			}
			pop.ComputeFitness()

			fitnesses := []float64{
				pop.Individuals[0].Fitness,
				pop.Individuals[1].Fitness,
				pop.Individuals[2].Fitness,
				pop.Individuals[3].Fitness,
			}
			So(fitnesses, ShouldResemble, []float64{8, 0, 1, 2})
		})

		Convey("when total fitness", func() {
			pop := Population{
				Individuals: []Individual{ind1, ind2, ind3, ind4},
				fitness:     fitness,
			}

			So(pop.Stats.TotalFitness, ShouldEqual, 0)
			pop.ComputeFitness()
			So(pop.Stats.TotalFitness, ShouldEqual, 8+0+1+2)
		})

		Convey("when sort", func() {
			pop := Population{
				Individuals: []Individual{ind1, ind2, ind3, ind4},
				fitness:     fitness,
			}

			pop.Sort()
			So(pop.Individuals, ShouldResemble, []Individual{ind1, ind2, ind3, ind4}) // unchanged

			pop.ComputeFitness()
			pop.Sort()
			So(pop.Individuals, ShouldResemble, []Individual{
				{Code: ind1.Code, Fitness: 8},
				{Code: ind4.Code, Fitness: 2},
				{Code: ind3.Code, Fitness: 1},
				{Code: ind2.Code, Fitness: 0},
			})
		})
	})
}
