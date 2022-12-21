package operator

import (
	"time"

	"github.com/sbiemont/galgogene/gene"
)

// Termination defines an ending condition for the engine
type Termination interface {
	// End returns true when the processing should end ; false otherwise
	End(pop gene.Population) Termination
}

// ------------------------------

// GenerationTermination should end processing when the ith generation is reached
type GenerationTermination struct {
	K int // The max generation to be reached
}

func (end *GenerationTermination) End(pop gene.Population) Termination {
	return condition(pop.Stats.GenerationNb >= end.K, end)
}

// ------------------------------

// ImprovementTermination should end processing when the total fitness
// has not increased since the previous generation
type ImprovementTermination struct {
	K                    int // The number of generations with the same improvement (default: 1)
	k                    int // The internal number of generations
	previousTotalFitness float64
}

func (end *ImprovementTermination) End(pop gene.Population) Termination {
	if end.previousTotalFitness == pop.Stats.TotalFitness {
		end.k++ // one more generation with same fitness
	} else {
		end.k = 0 // reset
	}

	k := getDefault(end.K, 1)
	end.previousTotalFitness = pop.Stats.TotalFitness
	return condition(end.k >= k, end)
}

// ------------------------------

// FitnessTermination should end processing when the elite reaches the defined fitness
type FitnessTermination struct {
	Fitness float64 // Min fitness
}

func (end *FitnessTermination) End(pop gene.Population) Termination {
	for _, individual := range pop.Individuals {
		if individual.Fitness >= end.Fitness {
			return end
		}
	}

	return nil
}

// ------------------------------

// DurationTermination should end processing when the total duration of each generation reaches a maximum
type DurationTermination struct {
	Duration time.Duration // Max duration
}

func (end *DurationTermination) End(pop gene.Population) Termination {
	return condition(pop.Stats.TotalDuration >= end.Duration, end)
}

// ------------------------------

// MultiTermination should en processing when one of the defined terminations
type MultiTermination []Termination

func (end MultiTermination) Use(t Termination) MultiTermination {
	return append(end, t)
}

func (end MultiTermination) End(pop gene.Population) Termination {
	for _, termination := range end {
		if termination.End(pop) != nil {
			return termination // ok, termination found
		}
	}

	return nil // no termination found
}

// ------------------------------

// Helper, returns the termination if cond is true
func condition(cond bool, termination Termination) Termination {
	if cond {
		return termination
	}

	return nil
}

// Helper, get the default value
func getDefault(value, deflt int) int {
	if value == 0 {
		return deflt
	}
	return value
}
