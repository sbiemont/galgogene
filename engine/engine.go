package engine

import (
	"errors"
	"time"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// Engine is the core element for running the algorithm
type Engine struct {
	Initializer     gene.Initializer
	Selection       operator.Selection
	CrossOver       operator.CrossOver
	Mutation        operator.Mutation
	Survivor        operator.Survivor
	Termination     operator.Termination
	Fitness         gene.Fitness
	OnNewGeneration func(gene.Population)
}

func (eng Engine) check() error {
	// Check presence
	switch {
	case eng.Fitness == nil:
		return errors.New("fitness must be set")
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

// Run the engine
// * popSize:        the number of individuals in a population
// * offspringSize:  the number of individuals in the offspring population
// * chromosomeSize: number of bases in a chromosome (in one individual)
func (eng Engine) Run(popSize, offspringSize, chromosomeSize int) (Solution, error) {
	start := time.Now()
	if err := eng.check(); err != nil {
		return Solution{}, err
	}

	// Init population size
	makeEven := func(v int) int {
		return v + v%2
	}
	if popSize <= 0 {
		popSize = 2
	}
	if offspringSize <= 0 {
		offspringSize = popSize
	}
	popSize = makeEven(popSize)
	offspringSize = makeEven(offspringSize)

	// Init first pop
	population := gene.NewPopulation(popSize)
	errInit := population.Init(chromosomeSize, eng.Initializer, eng.Fitness)
	if errInit != nil {
		return Solution{}, errInit
	}
	withBestIndividual := population
	withBestTotalFit := population
	eng.onNewGeneration(population)

	// Init channels
	chPopulation := make(chan gene.Population, 1)
	chSelection := make(chan gene.Chromosome, 20)
	chCrossover := make(chan gene.Chromosome, 20)
	chMutation := make(chan gene.Chromosome, 20)
	chIndividuals := make(chan gene.Individual, 20)
	chErr := make(chan error)
	go eng.selection(offspringSize, chPopulation, chSelection, chErr)
	go eng.crossover(chSelection, chCrossover)
	go eng.mutation(chCrossover, chMutation)
	go eng.fitness(chMutation, chIndividuals)
	defer close(chErr)
	defer close(chIndividuals)
	defer close(chMutation)
	defer close(chCrossover)
	defer close(chSelection)
	defer close(chPopulation)

	// Run until an ending condition is found
	var termination operator.Termination
	for ; termination == nil; termination = eng.Termination.End(population) {
		chPopulation <- population // start selection process
		var err error
		population, err = eng.nextGeneration(start, population, offspringSize, chIndividuals, chErr)
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
func (eng Engine) nextGeneration(
	start time.Time,
	parents gene.Population,
	offspringSize int,
	chIndividuals <-chan gene.Individual,
	chErr <-chan error,
) (gene.Population, error) {
	// Init
	offsprings := gene.NewPopulation(offspringSize)
	for i := range offspringSize {
		select {
		case err := <-chErr:
			return gene.Population{}, err
		case ind := <-chIndividuals:
			offsprings.Individuals[i] = ind
		}
	}
	offsprings.ComputeTotalFitness()

	// Survivors
	// new population has changed, compute global data like total fitness
	newPop := eng.Survivor.Survive(parents, offsprings)
	newPop.ComputeTotalFitness()
	newPop.ComputeRank()
	newPop.Stats.GenerationNb = parents.Stats.GenerationNb + 1
	newPop.Stats.TotalDuration = time.Since(start)
	return newPop, nil
}

// selection process: generate 1 selection per individual in the original population
func (eng Engine) selection(offspringSize int, in <-chan gene.Population, out chan<- gene.Chromosome, chErr chan<- error) {
	for population := range in {
		for range offspringSize {
			ind, err := eng.Selection.Select(population)
			if err != nil {
				chErr <- err
				return
			}
			out <- ind.Code
		}
	}
}

// crossover process: use 2 chromosomes and produce 2 new ones
func (eng Engine) crossover(in <-chan gene.Chromosome, out chan<- gene.Chromosome) {
	for chrm1 := range in {
		chrm2 := <-in
		if eng.CrossOver != nil {
			chrm1, chrm2 = eng.CrossOver.Mate(chrm1, chrm2)
		}
		out <- chrm1
		out <- chrm2
	}
}

// mutation process: mutate all chromosomes using the defined mutation function
func (eng Engine) mutation(in <-chan gene.Chromosome, out chan<- gene.Chromosome) {
	for mut := range in {
		if eng.Mutation != nil {
			mut = eng.Mutation.Mutate(mut)
		}
		out <- mut
	}
}

// fitness process: compute each individual fitness
func (eng Engine) fitness(in <-chan gene.Chromosome, out chan<- gene.Individual) {
	for chrm := range in {
		fitness := eng.Fitness(chrm)
		out <- gene.NewIndividual(chrm, fitness)
	}
}
