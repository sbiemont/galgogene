package operator

import (
	"time"

	"genalgo.git/gene"
)

// Termination defines an ending condition for the engine
type Termination interface {
	// End returns true when the processing should end ; false otherwise
	End(pop gene.Population) Termination
}

func condition(cond bool, termination Termination) Termination {
	if cond {
		return termination
	}

	return nil
}

// ------------------------------

// TerminationGeneration should end processing when the ith generation is reached
type TerminationGeneration struct {
	K int // The max generation to be reached
}

func (end *TerminationGeneration) End(pop gene.Population) Termination {
	return condition(pop.Stats.GenerationNb >= end.K, end)
}

// ------------------------------

// TerminationImprovement should end processing when the total fitness
// has not increased since the previous generation
type TerminationImprovement struct {
	previousTotalFitness float64
}

func (end *TerminationImprovement) End(pop gene.Population) Termination {
	stop := end.previousTotalFitness == pop.Stats.TotalFitness
	end.previousTotalFitness = pop.Stats.TotalFitness
	return condition(stop, end)
}

// ------------------------------

// TerminationAboveFitness should end processing when the elite reaches the defined fitness
type TerminationAboveFitness struct {
	Fitness float64 // Min fitness
}

func (end *TerminationAboveFitness) End(pop gene.Population) Termination {
	for _, individual := range pop.Individuals {
		if individual.Fitness >= end.Fitness {
			return end
		}
	}

	return nil
}

// ------------------------------

// TerminationDuration should end processing when the total duration of each generation reaches a maximum
type TerminationDuration struct {
	Duration time.Duration // Max duration
}

func (end *TerminationDuration) End(pop gene.Population) Termination {
	return condition(pop.Stats.TotalDuration >= end.Duration, end)
}

// ------------------------------

// MultiTermination should en processing when one of the defined terminations
type MultiTermination []Termination

func (end MultiTermination) End(pop gene.Population) Termination {
	for _, termination := range end {
		if termination.End(pop) != nil {
			return termination // ok, termination found
		}
	}

	return nil // no termination found
}
