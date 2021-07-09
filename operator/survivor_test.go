package operator

import (
	"testing"

	"genalgo.git/gene"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSurvivor(t *testing.T) {
	Convey("survivor", t, func() {
		pop1 := gene.Population{
			Individuals: []gene.Individual{
				{Fitness: 0.1},
				{Fitness: 0.5},
				{Fitness: 0.6},
				{Fitness: 0.9},
			},
		}
		pop2 := gene.Population{
			Individuals: []gene.Individual{
				{Fitness: 0.2},
				{Fitness: 0.4},
				{Fitness: 0.8},
			},
		}

		Convey("when copy all parents", func() {
			srv := SurvivorAddAllParents{}
			srv.Survive(pop1, &pop2)
			So(pop1.Individuals, ShouldResemble, []gene.Individual{ // pop1 unchanged
				{Fitness: 0.1},
				{Fitness: 0.5},
				{Fitness: 0.6},
				{Fitness: 0.9},
			})
			So(pop2.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.2}, // from pop2
				{Fitness: 0.4},
				{Fitness: 0.8},
				{Fitness: 0.1}, // from pop1
				{Fitness: 0.5},
				{Fitness: 0.6},
				{Fitness: 0.9},
			})
		})

		Convey("when copy k parents", func() {
			srv := SurvivorAddParentsElite{K: 2}
			srv.Survive(pop1, &pop2)
			So(pop1.Individuals, ShouldResemble, []gene.Individual{ // pop1 has been ordered
				{Fitness: 0.9},
				{Fitness: 0.6},
				{Fitness: 0.5},
				{Fitness: 0.1},
			})
			So(pop2.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.2}, // from pop2
				{Fitness: 0.4},
				{Fitness: 0.8},
				{Fitness: 0.9}, // from pop1 (ordered)
				{Fitness: 0.6},
			})
		})

		Convey("when keep best fitness", func() {
			srv := SurvivorElite{K: 2}
			srv.Survive(pop1, &pop2)
			So(pop1.Individuals, ShouldResemble, []gene.Individual{ // pop1 unchanged
				{Fitness: 0.1},
				{Fitness: 0.5},
				{Fitness: 0.6},
				{Fitness: 0.9},
			})
			So(pop2.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.8}, // from pop2
				{Fitness: 0.4},
			})
		})

		Convey("when multi survivors", func() {
			srv := MultiSurvivor{
				SurvivorAddParentsElite{K: 2},
				SurvivorElite{K: 2},
			}
			srv.Survive(pop1, &pop2)
			So(pop1.Individuals, ShouldResemble, []gene.Individual{ // pop1 has been ordered
				{Fitness: 0.9},
				{Fitness: 0.6},
				{Fitness: 0.5},
				{Fitness: 0.1},
			})
			So(pop2.Individuals, ShouldResemble, []gene.Individual{
				{Fitness: 0.9}, // from pop1
				{Fitness: 0.8}, // from pop2
			})
		})
	})
}
