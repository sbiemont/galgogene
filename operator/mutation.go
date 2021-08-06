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
type UniformMutation struct {
}

// Mate will mutate the first set of bits and leave the second one unchanged
func (um UniformMutation) Mutate(bits gene.Bits) gene.Bits {
	return mutate(bits, 0.5, func(b gene.Bits, _ int) uint8 {
		return b.Rand()
	})
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
