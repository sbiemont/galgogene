package operator

import (
	"genalgo.git/gene"
	"genalgo.git/random"
)

// Mutation defines a specific mutation on one set of bits and returns the mutated result
type Mutation interface {
	Mutate(bits gene.Bits) gene.Bits
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

// SwapMutation defines a random swap of 2 bits
type SwapMutation struct{}

// Mutate select 2 positions and swap the values
func (SwapMutation) Mutate(bits gene.Bits) gene.Bits {
	pos := random.Ints(0, bits.Len(), 2)
	pos1, pos2 := pos[0], pos[1]

	result := bits.Clone()
	result.Raw[pos1] = bits.Raw[pos2]
	result.Raw[pos2] = bits.Raw[pos1]
	return result
}

// ------------------------------

// TwoSwapMutation picks 2 points and inverts the subtour
type TwoSwapMutation struct{}

// Mutate select 2 positions and inverts the subtour
// eg.:
//   * input:  AB.CDEF.GH
//   * output: AB.FEDC.GH
func (TwoSwapMutation) Mutate(bits gene.Bits) gene.Bits {
	pos := random.Ints(0, bits.Len(), 2)
	pos1, pos2 := pos[0], pos[1]

	result := bits.Clone()
	for i := pos1; i <= pos2; i++ {
		result.Raw[i] = bits.Raw[pos2-i+pos1]
	}
	return result
}

// ------------------------------

// ProbaMutation is a probabilistic mutation
type ProbaMutation struct {
	rate float64  // Mutation rate
	mut  Mutation // Mutation operator
}

// NewProbaMutation build a new full instance of ProbaMutation
func NewProbaMutation(rate float64, mut Mutation) ProbaMutation {
	return ProbaMutation{
		rate: rate,
		mut:  mut,
	}
}

// MultiMutation defines a serie of mutations with a specific probability of beeing chosen.
// All or no mutations may be applied
type MultiMutation []ProbaMutation

func (mm MultiMutation) Mutate(bits gene.Bits) gene.Bits {
	res := bits
	for _, m := range mm {
		if random.Peek(m.rate) {
			res = m.mut.Mutate(res)
		}
	}
	return res
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
