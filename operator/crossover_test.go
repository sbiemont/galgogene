package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"

	. "github.com/smartystreets/goconvey/convey"
)

// AppliedCrossOver only check if a crossover is used or not
type AppliedCrossOver struct {
	IsApplied bool
}

func (mut *AppliedCrossOver) Mate(chrm1, chrm2 gene.Chromosome) (gene.Chromosome, gene.Chromosome) {
	mut.IsApplied = true
	return chrm1, chrm2
}

func TestCrossOvers(t *testing.T) {
	Convey("crossover", t, func() {
		chrm1 := newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 1})
		chrm2 := newChromosome([]gene.B{0, 0, 0, 0, 0, 0, 0, 0})

		Convey("when 1 split", func() {
			res1, res2 := crossOver(chrm1, chrm2, []int{4})
			So(res1, ShouldResemble, newChromosome([]gene.B{1, 1, 1, 1, 0, 0, 0, 0}))
			So(res2, ShouldResemble, newChromosome([]gene.B{0, 0, 0, 0, 1, 1, 1, 1}))
			So(chrm1, ShouldResemble, newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 1})) // unchanged
			So(chrm2, ShouldResemble, newChromosome([]gene.B{0, 0, 0, 0, 0, 0, 0, 0})) // unchanged
		})

		Convey("when 2 splits", func() {
			res1, res2 := crossOver(chrm1, chrm2, []int{3, 6})
			So(res1, ShouldResemble, newChromosome([]gene.B{1, 1, 1, 0, 0, 0, 1, 1}))
			So(res2, ShouldResemble, newChromosome([]gene.B{0, 0, 0, 1, 1, 1, 0, 0}))
		})

		Convey("when 3 splits", func() {
			res1, res2 := crossOver(chrm1, chrm2, []int{2, 4, 6})
			So(res1, ShouldResemble, newChromosome([]gene.B{1, 1, 0, 0, 1, 1, 0, 0}))
			So(res2, ShouldResemble, newChromosome([]gene.B{0, 0, 1, 1, 0, 0, 1, 1}))
		})

		Convey("when 3 splits with full bytes", func() {
			chrm1 := gene.Chromosome{
				Raw: []gene.B{0, 5, 10, 15, 20, 25, 30, 35},
			}
			chrm2 := gene.Chromosome{
				Raw: []gene.B{220, 225, 230, 235, 240, 245, 250, 255},
			}
			res1, res2 := crossOver(chrm1, chrm2, []int{2, 4, 6})
			So(res1, ShouldResemble, gene.Chromosome{
				Raw: []gene.B{0, 5, 230, 235, 20, 25, 250, 255},
			})
			So(res2, ShouldResemble, gene.Chromosome{
				Raw: []gene.B{220, 225, 10, 15, 240, 245, 30, 35},
			})
		})

		Convey("when same indexes", func() {
			res1, res2 := crossOver(chrm1, chrm2, []int{7, 7, 7})
			So(res1, ShouldResemble, newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 0}))
			So(res2, ShouldResemble, newChromosome([]gene.B{0, 0, 0, 0, 0, 0, 0, 1}))
		})

		Convey("when one point crossover", func() {
			random.Seed(42)
			res1, res2 := OnePointCrossOver{}.Mate(chrm1, chrm2)
			So(res1.Raw, ShouldResemble, []gene.B{1, 1, 0, 0, 0, 0, 0, 0})
			So(res2.Raw, ShouldResemble, []gene.B{0, 0, 1, 1, 1, 1, 1, 1})
		})

		Convey("when two point crossover", func() {
			random.Seed(42)
			res1, res2 := TwoPointsCrossOver{}.Mate(chrm1, chrm2)
			So(res1.Raw, ShouldResemble, []gene.B{1, 1, 0, 0, 0, 0, 0, 1})
			So(res2.Raw, ShouldResemble, []gene.B{0, 0, 1, 1, 1, 1, 1, 0})
		})
	})

	Convey("uniform crossover", t, func() {
		chrm1 := newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 1})
		chrm2 := newChromosome([]gene.B{0, 0, 0, 0, 0, 0, 0, 0})

		Convey("when rate 0.5", func() {
			res1, res2 := uniformCrossOver(chrm1, chrm2, 0.5)
			So(res1, ShouldNotResemble, newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 1}))
			So(res2, ShouldNotResemble, newChromosome([]gene.B{0, 0, 0, 0, 0, 0, 0, 0}))
		})

		Convey("when uniform", func() {
			random.Seed(42)
			res1, res2 := UniformCrossOver{}.Mate(chrm1, chrm2)
			So(res1.Raw, ShouldResemble, []gene.B{1, 1, 0, 1, 1, 1, 0, 1})
			So(res2.Raw, ShouldResemble, []gene.B{0, 0, 1, 0, 0, 0, 1, 0})
		})
	})

	Convey("davis' order crossover", t, func() {
		chrm1 := newChromosome([]gene.B{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'})
		chrm2 := newChromosome([]gene.B{'I', 'H', 'G', 'F', 'E', 'D', 'C', 'B', 'A'})

		type test struct {
			name     string
			pos      [2]int
			expected []gene.B
		}

		tests := []test{
			{
				name:     "pos [0,0]",
				pos:      [2]int{0, 0},
				expected: []gene.B{'A', 'I', 'H', 'G', 'F', 'E', 'D', 'C', 'B'},
			},
			{
				name:     "pos [0,1]",
				pos:      [2]int{0, 1},
				expected: []gene.B{'A', 'B', 'I', 'H', 'G', 'F', 'E', 'D', 'C'},
			},
			{
				name:     "pos [0,4]",
				pos:      [2]int{0, 4},
				expected: []gene.B{'A', 'B', 'C', 'D', 'E', 'I', 'H', 'G', 'F'},
			},
			{
				name:     "pos [0,8]",
				pos:      [2]int{0, 8},
				expected: []gene.B{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'},
			},
			{
				name:     "pos [2,5]",
				pos:      [2]int{2, 5},
				expected: []gene.B{'I', 'H', 'C', 'D', 'E', 'F', 'G', 'B', 'A'},
			},
			{
				name:     "pos [7,8]",
				pos:      [2]int{7, 8},
				expected: []gene.B{'G', 'F', 'E', 'D', 'C', 'B', 'A', 'H', 'I'},
			},
			{
				name:     "pos [8,8]",
				pos:      [2]int{8, 8},
				expected: []gene.B{'H', 'G', 'F', 'E', 'D', 'C', 'B', 'A', 'I'},
			},
		}

		for _, t := range tests {
			res := davisOrderCrossOver(chrm1, chrm2, t.pos[0], t.pos[1])
			SoMsg(t.name, res.Raw, ShouldResemble, t.expected)
		}
	})

	Convey("davis' order crossover with duplicated values", t, func() {
		chrm1 := newChromosome([]gene.B{'A', 'A', 'B', 'B', 'C', 'C'})
		chrm2 := newChromosome([]gene.B{'C', 'C', 'B', 'B', 'A', 'A'})

		So(davisOrderCrossOver(chrm1, chrm2, 1, 3).Raw, ShouldResemble, []gene.B{'C', 'A', 'B', 'B', 'C', 'A'})
	})

	Convey("uniform order crossover", t, func() {
		chrm1 := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8, 9})
		chrm2 := newChromosome([]gene.B{3, 4, 7, 2, 8, 9, 1, 6, 5})

		mask0 := []int{2, 4, 5, 6, 8}
		mask1 := []int{0, 1, 3, 7}
		So(uniformOrderCrossOver(chrm1, chrm2, mask0, mask1).Raw, ShouldResemble, []gene.B{1, 2, 3, 4, 7, 9, 6, 8, 5})
		So(uniformOrderCrossOver(chrm2, chrm1, mask0, mask1).Raw, ShouldResemble, []gene.B{3, 4, 1, 2, 5, 7, 8, 6, 9})
	})

	Convey("partially matched crossover", t, func() {
		chrm1 := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
		chrm2 := newChromosome([]gene.B{3, 7, 5, 1, 6, 8, 2, 4})

		pos1, pos2 := 3, 6
		So(partiallyMatchCrossOver(chrm1, chrm2, pos1, pos2).Raw, ShouldResemble, []gene.B{4, 2, 3, 1, 6, 8, 7, 5})
		So(partiallyMatchCrossOver(chrm2, chrm1, pos1, pos2).Raw, ShouldResemble, []gene.B{3, 7, 8, 4, 5, 6, 2, 1})

		pos1, pos2 = 0, 1
		So(partiallyMatchCrossOver(chrm1, chrm2, pos1, pos2).Raw, ShouldResemble, []gene.B{3, 2, 1, 4, 5, 6, 7, 8})
		So(partiallyMatchCrossOver(chrm2, chrm1, pos1, pos2).Raw, ShouldResemble, []gene.B{1, 7, 5, 3, 6, 8, 2, 4})

		pos1, pos2 = 7, 8
		So(partiallyMatchCrossOver(chrm1, chrm2, pos1, pos2).Raw, ShouldResemble, []gene.B{1, 2, 3, 8, 5, 6, 7, 4})
		So(partiallyMatchCrossOver(chrm2, chrm1, pos1, pos2).Raw, ShouldResemble, []gene.B{3, 7, 5, 1, 6, 4, 2, 8})
	})

	Convey("multi crossovers", t, func() {
		co1 := &AppliedCrossOver{}
		co2 := &AppliedCrossOver{}

		Convey("when none", func() {
			random.Seed(42)
			co1.IsApplied = false
			co2.IsApplied = false
			_, _ = MultiCrossOver{}.Use(0.01, co1).Use(0.01, co2).Mate(gene.Chromosome{}, gene.Chromosome{})
			So(co1.IsApplied, ShouldBeFalse)
			So(co2.IsApplied, ShouldBeFalse)
		})

		Convey("when first", func() {
			random.Seed(42)
			co1.IsApplied = false
			co2.IsApplied = false
			_, _ = MultiCrossOver{}.Use(1, co1).Use(1, co2).Mate(gene.Chromosome{}, gene.Chromosome{})
			So(co1.IsApplied, ShouldBeTrue)
			So(co2.IsApplied, ShouldBeFalse)
		})

		Convey("when last", func() {
			random.Seed(42)
			co1.IsApplied = false
			co2.IsApplied = false
			_, _ = MultiCrossOver{}.Use(0.1, co1).Use(1, co2).Mate(gene.Chromosome{}, gene.Chromosome{})
			So(co1.IsApplied, ShouldBeFalse)
			So(co2.IsApplied, ShouldBeTrue)
		})

		Convey("when all", func() {
			random.Seed(42)
			co1.IsApplied = false
			co2.IsApplied = false
			_, _ = MultiCrossOver{ApplyAll: true}.Use(1, co1).Use(1, co2).Mate(gene.Chromosome{}, gene.Chromosome{})
			So(co1.IsApplied, ShouldBeTrue)
			So(co2.IsApplied, ShouldBeTrue)
		})
	})
}

func TestFinder(t *testing.T) {
	Convey("new finder", t, func() {
		f := newFinder()
		So(f, ShouldResemble, finder{
			idx:        0,
			usedValues: map[gene.B]int{}, // empty but not nil map
		})
	})

	Convey("finder + use + duplicated values", t, func() {
		chrm := gene.Chromosome{
			// Indexes:   0   1   2   3   4   5   6   7   8   9  10  11  12  13
			Raw: []gene.B{42, 43, 44, 45, 46, 47, 48, 42, 43, 44, 45, 46, 47, 48},
		}

		f := newFinder()
		// 42
		f.useValue(43)
		// 44
		f.useValue(45)
		f.useValue(46)
		f.useValue(47)
		// 48

		So(f.nextUnused(chrm), ShouldEqual, 42)
		So(f.nextUnused(chrm), ShouldEqual, 44)
		So(f.nextUnused(chrm), ShouldEqual, 48)
		So(f.nextUnused(chrm), ShouldEqual, 42) // duplicated values
		So(f.nextUnused(chrm), ShouldEqual, 43)
		So(f.nextUnused(chrm), ShouldEqual, 44)
		So(f.nextUnused(chrm), ShouldEqual, 45)
		So(f.nextUnused(chrm), ShouldEqual, 46)
		So(f.nextUnused(chrm), ShouldEqual, 47)
		So(f.nextUnused(chrm), ShouldEqual, 48)
		So(f.nextUnused(chrm), ShouldEqual, 0)
	})
}
