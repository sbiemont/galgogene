package factory

import (
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// Random factory for genes with random values
type Random struct {
	Initializer randomInitializer
	Selection   commonSelection
	Survivor    commonSurvivor
	Mutation    randomMutation
	CrossOver   randomCrossOver
	Termination commonTermination
}

// Initializer

type randomInitializer struct{}

func (f randomInitializer) Random(maxValue uint8) gene.RandomInitializer {
	return gene.RandomInitializer{
		MaxValue: maxValue,
	}
}

// CrossOver

type randomCrossOver struct{}

func (f randomCrossOver) OnePoint() operator.OnePointCrossOver {
	return operator.OnePointCrossOver{}
}

func (f randomCrossOver) TwoPoints() operator.TwoPointsCrossOver {
	return operator.TwoPointsCrossOver{}
}

func (f randomCrossOver) Uniform() operator.UniformCrossOver {
	return operator.UniformCrossOver{}
}

func (f randomCrossOver) Multi() operator.MultiCrossOver {
	return operator.MultiCrossOver{}
}

// Mutation

type randomMutation struct{}

func (f randomMutation) Unique() operator.UniqueMutation {
	return operator.UniqueMutation{}
}

func (f randomMutation) Uniform() operator.UniformMutation {
	return operator.UniformMutation{}
}

func (f randomMutation) Multi() operator.MultiMutation {
	return operator.MultiMutation{}
}
