package operator

import (
	"math/rand"

	"genalgo.git/gene"
	"genalgo.git/random"
)

// Mutation examples:
// https://www.tutorialspoint.com/genetic_algorithms/genetic_algorithms_mutation.htm

// Mutation defines a specific mutation on one set of bits and returns the mutated result
// Notes:
// * a mutation overrides some bits with new random values
// * a permutation randomly reorders some bits (without changing the values)
type Mutation interface {
	Mutate(bits gene.Bits) gene.Bits
}

// ------------------------------

// UniqueMutation selects one unique bit and flips its value (using the max value)
type UniqueMutation struct{}

// Mutate each bit with a probability of 50%
func (UniqueMutation) Mutate(bits gene.Bits) gene.Bits {
	i := rand.Intn(bits.Len())
	result := bits.Clone()
	result.Raw[i] = result.Rand()
	return result
}

// ------------------------------

// UniformMutation defines a random mutation of bits
type UniformMutation struct{}

// Mutate each bit with a probability of 50%
func (UniformMutation) Mutate(bits gene.Bits) gene.Bits {
	return mutate(bits, 0.5, func(b gene.Bits, _ int) uint8 {
		return b.Rand()
	})
}

// ------------------------------

// SwapPermutation defines a random swap of 2 bits
type SwapPermutation struct{}

// Mutate select 2 positions and swap the values
func (SwapPermutation) Mutate(bits gene.Bits) gene.Bits {
	return permutation(bits, func(in gene.Bits, out *gene.Bits, pos1, pos2 int) {
		out.Raw[pos1] = in.Raw[pos2]
		out.Raw[pos2] = in.Raw[pos1]
	})
}

// ------------------------------

// InversionPermutation picks 2 points and inverts the subtour
type InversionPermutation struct{}

// Mutate select 2 positions and inverts the subtour
// eg.:
//   * input:  AB.CDEF.GH
//   * output: AB.FEDC.GH
func (InversionPermutation) Mutate(bits gene.Bits) gene.Bits {
	return permutation(bits, func(in gene.Bits, out *gene.Bits, pos1, pos2 int) {
		for i := pos1; i <= pos2; i++ {
			out.Raw[i] = in.Raw[pos2-i+pos1]
		}
	})
}

// ------------------------------

// SramblePermutation picks 2 points and shuffles the subtour
type SramblePermutation struct{}

// Mutate select 2 positions and shuffles the subtour
// eg.:
//   * input:  AB.CDEF.GH
//   * output: AB.ECFD.GH
func (SramblePermutation) Mutate(bits gene.Bits) gene.Bits {
	return permutation(bits, func(in gene.Bits, out *gene.Bits, pos1, pos2 int) {
		indexes := rand.Perm(pos2 - pos1)
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
type MultiMutation []probaMutation

func (mm MultiMutation) Mutate(bits gene.Bits) gene.Bits {
	res := bits
	for _, m := range mm {
		if random.Peek(m.rate) {
			res = m.mut.Mutate(res)
		}
	}
	return res
}

// Use the given proba mutation
func (mm MultiMutation) Use(rate float64, mut Mutation) MultiMutation {
	return append(mm, probaMutation{
		rate: rate,
		mut:  mut,
	})
}

// ------------------------------

// mutate inverts some bits using a mutation rate
func mutate(bits gene.Bits, rate float64, fct func(gene.Bits, int) uint8) gene.Bits {
	result := bits.Clone()
	for i := 0; i < result.Len(); i++ {
		if random.Peek(rate) {
			result.Raw[i] = fct(result, i)
		}
	}
	return result
}

// mutationPositions helps creating 2 random positions in [0 ; bits.Len[
func mutationPositions(bits gene.Bits) (int, int) {
	pos := random.Ints(0, bits.Len(), 2)
	return pos[0], pos[1]
}

func permutation(bits gene.Bits, apply func(in gene.Bits, out *gene.Bits, pos1, pos2 int)) gene.Bits {
	pos := random.Ints(0, bits.Len(), 2)
	result := bits.Clone()
	apply(bits, &result, pos[0], pos[1])
	return result
}
