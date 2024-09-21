package gene

import (
	"testing"

	"github.com/sbiemont/galgogene/random"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChromosome(t *testing.T) {
	Convey("chromosome", t, func() {
		Convey("new chromosome", func() {
			chrm := NewChromosome(8, 0)
			So(chrm.Raw, ShouldResemble, []B{0, 0, 0, 0, 0, 0, 0, 0})
		})

		Convey("new chromosome random", func() {
			random.Seed(42)
			chrm := NewChromosomeRandom(8, 1)
			So(chrm.Raw, ShouldResemble, []B{0, 1, 0, 0, 1, 1, 0, 1})
		})

		Convey("len", func() {
			chrm := NewChromosome(8, 0)
			So(chrm.Len(), ShouldEqual, 8)
		})

		Convey("clone", func() {
			chrm := Chromosome{
				Raw:      []B{10, 20, 30, 40},
				maxValue: 42,
			}
			res := chrm.Clone()
			So(res, ShouldResemble, Chromosome{
				Raw:      []B{10, 20, 30, 40},
				maxValue: 42,
			})
		})

		Convey("new", func() {
			chrm := Chromosome{
				Raw:      []B{1, 2, 3, 4},
				maxValue: 42,
			}
			res := chrm.New()
			So(res, ShouldResemble, Chromosome{
				Raw:      []B{0, 0, 0, 0},
				maxValue: 42,
			})

			Convey("rand", func() {
				chrm := Chromosome{
					maxValue: 42,
				}
				random.Seed(42)
				res := chrm.Rand()
				So(res, ShouldEqual, B(32))
			})

			Convey("string", func() {
				chrm := Chromosome{
					Raw: []B{65, 66, 67, 68},
				}
				res := chrm.String()
				So(res, ShouldEqual, "ABCD")
			})
		})
	})
}
