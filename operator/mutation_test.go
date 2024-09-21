package operator

import (
	"testing"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"

	. "github.com/smartystreets/goconvey/convey"
)

func newChromosome(bases []gene.B) gene.Chromosome {
	res := gene.NewChromosome(0, 42)
	res.Raw = bases
	return res
}

// AppliedMutation only check if a mutation is used or not
type AppliedMutation struct {
	IsApplied bool
}

func (mut *AppliedMutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	mut.IsApplied = true
	return chrm
}

func TestMutations(t *testing.T) {
	Convey("mutate", t, func() {
		chrm := newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 1})
		toZero := func(gene.Chromosome, int) gene.B { return 0 }

		Convey("when none mutated", func() {
			res := mutate(chrm, 0.0, toZero)
			So(res, ShouldResemble, newChromosome([]gene.B{1, 1, 1, 1, 1, 1, 1, 1}))
		})

		Convey("when some mutated", func() {
			random.Seed(42)
			res := mutate(chrm, 0.5, toZero)
			So(res, ShouldResemble, newChromosome([]gene.B{0, 0, 1, 0, 0, 0, 1, 0}))
		})

		Convey("when all mutated", func() {
			res := mutate(chrm, 1.0, toZero)
			So(res, ShouldResemble, newChromosome([]gene.B{0, 0, 0, 0, 0, 0, 0, 0}))
		})
	})

	Convey("swap permutation", t, func() {
		chrm := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
		mutation := SwapPermutation{}

		random.Seed(42)
		result := mutation.Mutate(chrm)
		So(chrm.Raw, ShouldResemble, []gene.B{1, 2, 3, 4, 5, 6, 7, 8})
		So(result.Raw, ShouldResemble, []gene.B{1, 2, 8, 4, 5, 6, 7, 3})
	})

	Convey("inversion permutation", t, func() {
		Convey("when permutation", func() {
			chrm := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
			mutation := func(in gene.Chromosome, out *gene.Chromosome, _, _ int) {
				// Force positions
				pos1 := 2
				pos2 := 6
				for i := pos1; i <= pos2; i++ {
					out.Raw[i] = in.Raw[pos2-i+pos1]
				}
			}

			result := permutation(chrm, mutation)
			So(chrm.Raw, ShouldResemble, []gene.B{1, 2, 3, 4, 5, 6, 7, 8})
			So(result.Raw, ShouldResemble, []gene.B{1, 2, 7, 6, 5, 4, 3, 8})
		})

		Convey("when mutation", func() {
			random.Seed(5)
			chrm := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
			res := InversionPermutation{}.Mutate(chrm)
			So(res.Raw, ShouldResemble, []gene.B{1, 2, 3, 6, 5, 4, 7, 8})
		})
	})

	Convey("scramble permutation", t, func() {
		random.Seed(9)
		chrm := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
		res := ScramblePermutation{}.Mutate(chrm)
		So(res.Raw, ShouldResemble, []gene.B{1, 2, 5, 4, 6, 3, 7, 8})
	})

	Convey("unique mutation", t, func() {
		random.Seed(42)
		chrm := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
		res := UniqueMutation{}.Mutate(chrm)
		So(res.Raw, ShouldResemble, []gene.B{1, 2, 16, 4, 5, 6, 7, 8})
	})

	Convey("uniform mutation", t, func() {
		random.Seed(42)
		chrm := newChromosome([]gene.B{1, 2, 3, 4, 5, 6, 7, 8})
		res := UniformMutation{}.Mutate(chrm)
		So(res.Raw, ShouldResemble, []gene.B{16, 2, 12, 31, 25, 6, 7, 8})
	})

	Convey("multi mutations", t, func() {
		mut1 := &AppliedMutation{}
		mut2 := &AppliedMutation{}

		Convey("when none", func() {
			random.Seed(42)
			mut1.IsApplied = false
			mut2.IsApplied = false
			_ = MultiMutation{}.Use(0.01, mut1).Use(0.01, mut2).Mutate(gene.Chromosome{})
			So(mut1.IsApplied, ShouldBeFalse)
			So(mut2.IsApplied, ShouldBeFalse)
		})

		Convey("when first", func() {
			random.Seed(42)
			mut1.IsApplied = false
			mut2.IsApplied = false
			_ = MultiMutation{}.Use(1, mut1).Use(1, mut2).Mutate(gene.Chromosome{})
			So(mut1.IsApplied, ShouldBeTrue)
			So(mut2.IsApplied, ShouldBeFalse)
		})

		Convey("when last", func() {
			random.Seed(42)
			mut1.IsApplied = false
			mut2.IsApplied = false
			_ = MultiMutation{}.Use(0.1, mut1).Use(1, mut2).Mutate(gene.Chromosome{})
			So(mut1.IsApplied, ShouldBeFalse)
			So(mut2.IsApplied, ShouldBeTrue)
		})

		Convey("when all", func() {
			random.Seed(42)
			mut1.IsApplied = false
			mut2.IsApplied = false
			_ = MultiMutation{ApplyAll: true}.Use(1, mut1).Use(1, mut2).Mutate(gene.Chromosome{})
			So(mut1.IsApplied, ShouldBeTrue)
			So(mut2.IsApplied, ShouldBeTrue)
		})
	})
}
