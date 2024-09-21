package operator

import (
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"
)

// Mutation examples:
// https://www.tutorialspoint.com/genetic_algorithms/genetic_algorithms_mutation.htm

// Mutation defines a specific mutation on one set of bases and returns the mutated result
// Notes:
// * a mutation overrides some bases with new random values
// * a permutation randomly reorders some bases (without changing the values)
type Mutation interface {
	Mutate(chrm gene.Chromosome) gene.Chromosome
}

// ------------------------------

// UniqueMutation selects one unique bit and flips its value (using the max value)
type UniqueMutation struct{}

// Mutate a unique bit in the gene
func (UniqueMutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	i := random.IntN(chrm.Len())
	result := chrm.Clone()
	result.Raw[i] = result.Rand()
	return result
}

// ------------------------------

// UniformMutation defines a random mutation of bases
type UniformMutation struct{}

// Mutate each bit with a probability of 50%
func (UniformMutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	return mutate(chrm, 0.5, func(b gene.Chromosome, _ int) gene.B {
		return b.Rand()
	})
}

// ------------------------------

// SwapPermutation defines a random swap of 2 bases
type SwapPermutation struct{}

// Mutate select 2 positions and swap the values
func (SwapPermutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	return permutation(chrm, func(in gene.Chromosome, out *gene.Chromosome, pos1, pos2 int) {
		out.Raw[pos1] = in.Raw[pos2]
		out.Raw[pos2] = in.Raw[pos1]
	})
}

// ------------------------------

// InversionPermutation picks 2 points and inverts the subtour
type InversionPermutation struct{}

// Mutate select 2 positions and inverts the subtour
// eg.:
//   - input:  AB.CDEF.GH
//   - output: AB.FEDC.GH
func (InversionPermutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	return permutation(chrm, func(in gene.Chromosome, out *gene.Chromosome, pos1, pos2 int) {
		for i := pos1; i <= pos2; i++ {
			out.Raw[i] = in.Raw[pos2-i+pos1]
		}
	})
}

// ------------------------------

// ScramblePermutation picks 2 points and shuffles the subtour
type ScramblePermutation struct{}

// Mutate select 2 positions and shuffles the subtour
// eg.:
//   - input:  AB.CDEF.GH
//   - output: AB.ECFD.GH
func (ScramblePermutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	return permutation(chrm, func(in gene.Chromosome, out *gene.Chromosome, pos1, pos2 int) {
		indexes := random.Perm(pos2 - pos1)
		for i, index := range indexes {
			out.Raw[pos1+i] = in.Raw[pos1+index]
		}
	})
}

// ------------------------------

// probaMutation is a probabilistic mutation
type probaMutation struct {
	rate float64  // Mutation rate
	mut  Mutation // Mutation operator
}

// MultiMutation defines a serie of mutations with a specific probability of beeing chosen.
// All or no mutations may be applied
type MultiMutation struct {
	ApplyAll  bool // Set it to true, otherwise, processing stops at the first mutation to be applied
	mutations []probaMutation
}

func (mm MultiMutation) Mutate(chrm gene.Chromosome) gene.Chromosome {
	res := chrm
	for _, m := range mm.mutations {
		if random.Peek(m.rate) {
			res = m.mut.Mutate(res)
			if !mm.ApplyAll {
				return res
			}
		}
	}
	return res
}

// Use the given proba mutation
func (mm MultiMutation) Use(rate float64, mut Mutation) MultiMutation {
	return MultiMutation{
		ApplyAll: mm.ApplyAll,
		mutations: append(mm.mutations, probaMutation{
			rate: rate,
			mut:  mut,
		}),
	}
}

// ------------------------------

// mutate inverts some bases using a mutation rate
func mutate(chrm gene.Chromosome, rate float64, fct func(gene.Chromosome, int) gene.B) gene.Chromosome {
	result := chrm.Clone()
	for i := range result.Len() {
		if random.Peek(rate) {
			result.Raw[i] = fct(result, i)
		}
	}
	return result
}

func permutation(chrm gene.Chromosome, apply func(in gene.Chromosome, out *gene.Chromosome, pos1, pos2 int)) gene.Chromosome {
	pos := random.OrderedInts(0, chrm.Len(), 2)
	if pos[0] == pos[1] { // unchanged pos, leave bases unchanged
		return chrm
	}
	result := chrm.Clone()
	apply(chrm, &result, pos[0], pos[1])
	return result
}
