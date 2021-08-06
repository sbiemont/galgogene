package engine

import (
	"errors"
	"time"

	"genalgo.git/gene"
	"genalgo.git/operator"
	"genalgo.git/random"
)

type Engine struct {
	Selection       operator.Selection
	CrossOver       operator.CrossOver
	Mutation        operator.Mutation
	Survivor        operator.Survivor
	Termination     operator.Termination
	OnNewGeneration func(gene.Population)
}

func (eng Engine) Run(popSize, bitsSize int, fitness gene.FitnessFct) (gene.Population, operator.Termination, error) {
	return eng.RunWithMaxValue(popSize, bitsSize, 0x01, fitness)
}

func (eng Engine) RunWithMaxValue(
	popSize,
	bitsSize int,
	maxValue uint8,
	fitness gene.FitnessFct,
) (gene.Population, operator.Termination, error) {
	if err := eng.check(); err != nil {
		return gene.Population{}, nil, err
	}

	// Init new pop
	population := gene.NewPopulation(popSize, fitness)
	population.Init(bitsSize, maxValue)
	eng.onNewGeneration(population)

	// Run until an ending condition is found
	var termination operator.Termination
	for ; termination == nil; termination = eng.Termination.End(population) {
		time.Sleep(50 * time.Nanosecond)
		// New generation
		var err error
		population, err = eng.nextGeneration(population)
		if err != nil {
			return gene.Population{}, nil, err
		}

		// Custom action
		eng.onNewGeneration(population)
	}

	return population, termination, nil
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
		ind1, err1 := eng.Selection.Select(parents)
		if err1 != nil {
			return gene.Population{}, err1
		}
		ind2, err2 := eng.Selection.Select(parents)
		if err2 != nil {
			return gene.Population{}, err2
		}

		// Crossover
		mut1, mut2 := eng.CrossOver.Mate(ind1.Code, ind2.Code)

		// Mutation
		if eng.Mutation != nil {
			if random.Peek(0.5) {
				mut1 = eng.Mutation.Mutate(mut1)
			} else {
				mut2 = eng.Mutation.Mutate(mut2)
			}
		}

		// Add new individuals to the new generation
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
	if eng.Selection == nil {
		return errors.New("selection must be set")
	}

	if eng.CrossOver == nil {
		return errors.New("crossover must be set")
	}

	if eng.Survivor == nil {
		return errors.New("survivor must be set")
	}

	if eng.Termination == nil {
		return errors.New("termination must be set")
	}

	return nil
}
