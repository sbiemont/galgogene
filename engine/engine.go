package engine

import (
	"errors"
	"time"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// Number of parallel go-routines for computations
const nbGoRoutines int = 4

// Engine is the core element for running the algorithm
type Engine struct {
	Initializer     gene.Initializer
	Selection       operator.Selection
	CrossOver       operator.CrossOver
	Mutation        operator.Mutation
	Survivor        operator.Survivor
	Termination     operator.Termination
	OnNewGeneration func(gene.Population)
}

// Solution after running the engine
// * best population (with elite individual)
// * best population (with max total fitness)
// * termination operator triggered
// * error (if any)
type Solution struct {
	PopWithBestIndividual   gene.Population      // Population with best computed individual
	PopWithBestTotalFitness gene.Population      // Population with best total fitness computed
	Termination             operator.Termination // Termination that triggered the end of computation
}

// Run the engine
// * popSize: the number of individuals in a population
// * bitsSize: number of bits in a gene (in one individual)
// * fitness: the function that computes the gene's score
func (eng Engine) Run(popSize, bitsSize int, fitness gene.FitnessFct) (Solution, error) {
	if err := eng.check(); err != nil {
		return Solution{}, err
	}

	// Init new pop
	population := gene.NewPopulation(popSize, fitness, eng.Initializer)
	withBestIndividual := population
	withBestTotalFit := population
	errInit := population.Init(bitsSize)
	if errInit != nil {
		return Solution{}, errInit
	}
	eng.onNewGeneration(population)

	// Run until an ending condition is found
	var termination operator.Termination
	for ; termination == nil; termination = eng.Termination.End(population) {
		// New generation
		var err error
		population, err = eng.nextGeneration(population)
		if err != nil {
			return Solution{}, err
		}

		// Custom action
		eng.onNewGeneration(population)
		if population.Stats.TotalFitness > withBestTotalFit.Stats.TotalFitness {
			withBestTotalFit = population
		}
		if population.Stats.Elite.Fitness > withBestIndividual.Stats.Elite.Fitness {
			withBestIndividual = population
		}
	}

	return Solution{
		PopWithBestIndividual:   withBestIndividual,
		PopWithBestTotalFitness: withBestTotalFit,
		Termination:             termination,
	}, nil
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

	err := runParallelBatch(popSize, nbGoRoutines, func(from, to, _ int) error {
		for i := from; i < to; i++ {
			// Select 2 individuals
			ind1, err1 := eng.Selection.Select(parents)
			if err1 != nil {
				return err1
			}
			ind2, err2 := eng.Selection.Select(parents)
			if err2 != nil {
				return err2
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
		return nil
	})
	if err != nil {
		return gene.Population{}, err
	}

	// Compute all fitnesses
	newPop.ComputeFitness()

	// Survivors
	err = eng.Survivor.Survive(parents, &newPop)
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
	// Check presence
	switch {
	case eng.Initializer == nil:
		return errors.New("initializer must be set")
	case eng.Selection == nil:
		return errors.New("selection must be set")
	case eng.CrossOver == nil:
		return errors.New("crossover must be set")
	case eng.Survivor == nil:
		return errors.New("survivor must be set")
	case eng.Termination == nil:
		return errors.New("termination must be set")
	default:
		return nil
	}
}
