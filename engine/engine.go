package engine

import (
	"errors"
	"time"

	"genalgo.git/gene"
	"genalgo.git/operator"
)

type Engine struct {
	Selector        operator.Selector
	Mutator         operator.Mutator
	Survivor        operator.Survivor
	Ender           operator.Ender
	OnNewGeneration func(gene.Population)
}

func (eng Engine) Run(popSize, bitsSize int, fitness gene.FitnessFct) (gene.Population, operator.Ender, error) {
	if err := eng.check(); err != nil {
		return gene.Population{}, nil, err
	}

	// Init new pop
	population := gene.NewPopulation(popSize, fitness)
	population.Init(bitsSize)
	eng.onNewGeneration(population)

	// Run until an ending condition is found
	var ender operator.Ender
	for ; ender == nil; ender = eng.Ender.End(population) {
		// New generation
		var err error
		population, err = eng.nextGeneration(population)
		if err != nil {
			return gene.Population{}, nil, err
		}

		// Custom action
		eng.onNewGeneration(population)
	}

	return population, ender, nil
}

// onNewGeneration calls the user method (only if defined)
func (eng Engine) onNewGeneration(population gene.Population) {
	if eng.OnNewGeneration != nil {
		eng.OnNewGeneration(population)
	}
}

// nextGeneration builds a new generation of individuals
func (eng Engine) nextGeneration(parents gene.Population) (gene.Population, error) {
	start := time.Now()

	// Init
	popSize := len(parents.Individuals)
	newPop := gene.NewPopulationFrom(2*popSize, parents)

	for i := 0; i < popSize; i++ {
		// Select 2 individuals
		ind1, err1 := eng.Selector.Select(parents)
		if err1 != nil {
			return gene.Population{}, err1
		}
		ind2, err2 := eng.Selector.Select(parents)
		if err2 != nil {
			return gene.Population{}, err2
		}

		// Mutate and add to the new population
		mut1, mut2 := eng.Mutator.Mate(ind1.Code, ind2.Code)
		newPop.Individuals[2*i] = gene.NewIndividual(mut1)
		newPop.Individuals[2*i+1] = gene.NewIndividual(mut2)
	}

	// Compute all fitnesses
	newPop.ComputeFitness()

	// Survivors + total fitness and additional data
	eng.Survivor.Survive(parents, &newPop)
	newPop.ComputeTotalFitness()
	newPop.Stats.GenerationNb = parents.Stats.GenerationNb + 1
	newPop.Stats.TotalDuration = parents.Stats.TotalDuration + time.Since(start)
	return newPop, nil
}

func (eng Engine) check() error {
	if eng.Selector == nil {
		return errors.New("selector must be set")
	}

	if eng.Mutator == nil {
		return errors.New("mutator must be set")
	}

	if eng.Survivor == nil {
		return errors.New("survivor must be set")
	}

	if eng.Ender == nil {
		return errors.New("ender must be set")
	}

	return nil
}
