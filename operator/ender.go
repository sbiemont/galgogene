package operator

import (
	"time"

	"genalgo.git/gene"
)

// Ender defines an ending condition for the engine
type Ender interface {
	// End returns true when the processing should end ; false otherwise
	End(pop gene.Population) Ender
}

func condition(cond bool, ender Ender) Ender {
	if cond {
		return ender
	}

	return nil
}

// ------------------------------

// EnderGeneration should end processing when the ith generation is reached
type EnderGeneration struct {
	K int // The max generation to be reached
}

func (end *EnderGeneration) End(pop gene.Population) Ender {
	return condition(pop.GenerationNb >= end.K, end)
}

// ------------------------------

// EnderImprovement should end processing when the total fitness
// has not increased since the previous generation
type EnderImprovement struct {
	previousTotalFitness float64
}

func (end *EnderImprovement) End(pop gene.Population) Ender {
	stop := end.previousTotalFitness == pop.TotalFitness
	end.previousTotalFitness = pop.TotalFitness
	return condition(stop, end)
}

// ------------------------------

// EnderAboveFitness should end processing when the elite reaches the defined fitness
type EnderAboveFitness struct {
	Fitness float64 // Min fitness
}

func (end *EnderAboveFitness) End(pop gene.Population) Ender {
	for _, individual := range pop.Individuals {
		if individual.Fitness >= end.Fitness {
			return end
		}
	}

	return nil
}

// ------------------------------

// EnderBelowFitness should end processing when the elite reaches the defined fitness
type EnderBelowFitness struct {
	Fitness float64 // Max fitness
}

func (end *EnderBelowFitness) End(pop gene.Population) Ender {
	for _, individual := range pop.Individuals {
		if individual.Fitness <= end.Fitness {
			return end
		}
	}

	return nil
}

// ------------------------------

// EnderDuration should end processing when the total duration of each generation reaches a maximum
type EnderDuration struct {
	Duration time.Duration // Max duration
}

func (end *EnderDuration) End(pop gene.Population) Ender {
	return condition(pop.TotalDuration >= end.Duration, end)
}

// ------------------------------

// MultiEnder should en processing when one of the defined enders
type MultiEnder []Ender

func (end MultiEnder) End(pop gene.Population) Ender {
	for _, ender := range end {
		if ender.End(pop) != nil {
			return ender // ok, ender found
		}
	}

	return nil // no ender found
}
