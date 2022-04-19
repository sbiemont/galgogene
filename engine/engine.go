package engine

import (
	"errors"
	"time"

	"galgogene.git/gene"
	"galgogene.git/operator"
)

type Engine struct {
	Initializer     gene.Initializer
	Selection       operator.Selection
	CrossOver       operator.CrossOver
	Mutation        operator.Mutation
	Survivor        operator.Survivor
	Termination     operator.Termination
	OnNewGeneration func(gene.Population)
}

func (eng Engine) Run(
	popSize,
	bitsSize int,
	fitness gene.FitnessFct,
) (gene.Population, gene.Population, operator.Termination, error) {
	if err := eng.check(); err != nil {
		return gene.Population{}, gene.Population{}, nil, err
	}

	// Init new pop
	population := gene.NewPopulation(popSize, fitness, eng.Initializer)
	best := population
	errInit := population.Init(bitsSize)
	if errInit != nil {
		return gene.Population{}, gene.Population{}, nil, errInit
	}
	eng.onNewGeneration(population)

	// Run until an ending condition is found
	var termination operator.Termination
	for ; termination == nil; termination = eng.Termination.End(population) {
		time.Sleep(50 * time.Nanosecond)
		// New generation
		var err error
		population, err = eng.nextGeneration(population)
		if err != nil {
			return gene.Population{}, gene.Population{}, nil, err
		}

		// Custom action
		eng.onNewGeneration(population)
		if population.Stats.TotalFitness > best.Stats.TotalFitness {
			best = population
		}
	}

	return population, best, termination, nil
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
		mut1, mut2 := ind1.Code, ind2.Code

		// Crossover
		if eng.CrossOver != nil {
			mut1, mut2 = eng.CrossOver.Mate(mut1, mut2)
		}

		// Mutation
		if eng.Mutation != nil {
			mut1 = eng.Mutation.Mutate(mut1)
			mut2 = eng.Mutation.Mutate(mut2)
		}

		// Add new individuals to the new generation
		newPop.Individuals[2*i] = gene.NewIndividual(mut1)
		newPop.Individuals[2*i+1] = gene.NewIndividual(mut2)
	}

	// Compute all fitnesses
	newPop.ComputeFitness()

	// Survivors
	err := eng.Survivor.Survive(parents, &newPop)
	if err != nil {
		return gene.Population{}, err
	}

	// Total fitness and additional data
	newPop.ComputeTotalFitness()
	newPop.AddRank()
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
