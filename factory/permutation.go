package factory

import (
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// Permutation factory for genes with permutation strategy
type Permutation struct {
	Initializer permutationInitializer
	Selection   commonSelection
	CrossOver   permutationCrossOver
	Mutation    permutationMutation
	Survivor    commonSurvivor
	Termination commonTermination
}

// Initializer

type permutationInitializer struct{}

func (f permutationInitializer) Permutation() gene.PermutationInitializer {
	return gene.PermutationInitializer{}
}

// CrossOver

type permutationCrossOver struct{}

func (f permutationCrossOver) DavisOrder() operator.DavisOrderCrossOver {
	return operator.DavisOrderCrossOver{}
}

func (f permutationCrossOver) UniformOrder() operator.UniformOrderCrossOver {
	return operator.UniformOrderCrossOver{}
}

func (f permutationCrossOver) Multi() operator.MultiCrossOver {
	return operator.MultiCrossOver{}
}

// Mutation

type permutationMutation struct{}

func (f permutationMutation) Swap() operator.SwapPermutation {
	return operator.SwapPermutation{}
}

func (f permutationMutation) Inversion() operator.InversionPermutation {
	return operator.InversionPermutation{}
}

func (f permutationMutation) Scramble() operator.ScramblePermutation {
	return operator.ScramblePermutation{}
}

func (f permutationMutation) Multi() operator.MultiMutation {
	return operator.MultiMutation{}
}
