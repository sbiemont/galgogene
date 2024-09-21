package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"

	. "github.com/smartystreets/goconvey/convey"
)

// AppliedSurvivor only check if a survivor is used or not
type AppliedSurvivor struct {
	IsApplied bool
}

func (mut *AppliedSurvivor) Survive(_, _ gene.Population) gene.Population {
	mut.IsApplied = true
	return gene.Population{}
}

func TestSurvivor(t *testing.T) {
	Convey("survivor", t, func() {
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
		pop2 := func() gene.Population {
			return gene.Population{
				Individuals: []gene.Individual{
					{Fitness: 0.2, Rank: 2},
					{Fitness: 0.4, Rank: 4},
					{Fitness: 0.8, Rank: 8},
				},
			}
		}

		Convey("when elite", func() {
			p1 := pop1()
			p2 := pop2()
			res := EliteSurvivor{}.Survive(p1, p2)
			So(p1, ShouldResemble, pop1()) // pop1 unchanged
			So(p2, ShouldResemble, pop2()) // pop2 unchanged
			So(res.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.9, Rank: 9},
				{Fitness: 0.8, Rank: 8},
				{Fitness: 0.6, Rank: 6},
				{Fitness: 0.5, Rank: 5},
			})
		})

		Convey("when rank", func() {
			p1 := pop1()
			p2 := pop2()
			res := RankSurvivor{}.Survive(p1, p2)
			So(p1, ShouldResemble, pop1()) // pop1 unchanged
			So(p2, ShouldResemble, pop2()) // pop2 unchanged
			So(res.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.1, Rank: 1},
				{Fitness: 0.2, Rank: 2},
				{Fitness: 0.4, Rank: 4},
				{Fitness: 0.5, Rank: 5},
			})
		})

		Convey("when random", func() {
			p1 := pop1()
			p2 := pop2()
			random.Seed(42)
			res := RandomSurvivor{}.Survive(p1, p2)
			So(p1, ShouldResemble, pop1()) // pop1 unchanged
			So(p2, ShouldResemble, pop2()) // pop2 unchanged
			So(res.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.1, Rank: 1},
				{Fitness: 0.5, Rank: 5},
				{Fitness: 0.8, Rank: 8},
				{Fitness: 0.4, Rank: 4},
			})
		})

		Convey("when multi", func() {
			Convey("when first applied", func() {
				survivor1 := AppliedSurvivor{}
				survivor2 := AppliedSurvivor{}
				random.Seed(42)
				_ = MultiSurvivor{}.
					Use(0.5, &survivor1).
					Otherwise(&survivor2).
					Survive(gene.Population{}, gene.Population{})
				So(survivor1.IsApplied, ShouldBeTrue)
				So(survivor2.IsApplied, ShouldBeFalse)
			})

			Convey("when second applied", func() {
				survivor1 := AppliedSurvivor{}
				survivor2 := AppliedSurvivor{}
				random.Seed(42)
				_ = MultiSurvivor{}.
					Use(0.1, &survivor1).
					Otherwise(&survivor2).
					Survive(gene.Population{}, gene.Population{})
				So(survivor1.IsApplied, ShouldBeFalse)
				So(survivor2.IsApplied, ShouldBeTrue)
			})
		})
	})
}
