package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"

	. "github.com/smartystreets/goconvey/convey"
)

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

		Convey("when keep best fitness", func() {
			p1 := pop1()
			p2 := pop2()
			srv := EliteSurvivor{}
			err := srv.Survive(p1, &p2)
			So(err, ShouldBeNil)
			So(p1.Individuals, ShouldResemble, []gene.Individual{ // pop1 unchanged
				{Fitness: 0.1, Rank: 1},
				{Fitness: 0.5, Rank: 5},
				{Fitness: 0.6, Rank: 6},
				{Fitness: 0.9, Rank: 9},
			})
			So(p2.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.9, Rank: 9},
				{Fitness: 0.8, Rank: 8},
				{Fitness: 0.6, Rank: 6},
				{Fitness: 0.5, Rank: 5},
			})
		})

		Convey("when keep minimum ranking", func() {
			p1 := pop1()
			p2 := pop2()
			srv := RankSurvivor{}
			err := srv.Survive(p1, &p2)
			So(err, ShouldBeNil)
			So(p1.Individuals, ShouldResemble, []gene.Individual{ // pop1 unchanged
				{Fitness: 0.1, Rank: 1},
				{Fitness: 0.5, Rank: 5},
				{Fitness: 0.6, Rank: 6},
				{Fitness: 0.9, Rank: 9},
			})
			So(p2.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.1, Rank: 1},
				{Fitness: 0.2, Rank: 2},
				{Fitness: 0.4, Rank: 4},
				{Fitness: 0.5, Rank: 5},
			})
		})
	})
}
