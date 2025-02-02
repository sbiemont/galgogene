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
	OnNewGeneration func(pop gene.Population, withBestIndividual gene.Population, withBestTotalFitness gene.Population)
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
	eng.onNewGeneration(population, population, population)

	// New channels
	chSelection := make(chan gene.Population, 1)
	chCrossover := make(chan gene.Chromosome, 20)
	chMutation := make(chan gene.Chromosome, 20)
	chFitness := make(chan gene.Chromosome, 20)
	chIndividuals := make(chan gene.Individual, 20)
	chOffsprings := make(chan gene.Population)
	chErr := make(chan error)
	chSolution := make(chan Solution)

	go eng.selection(offspringSize, chSelection, chCrossover, chErr)
	go eng.crossover(chCrossover, chMutation)
	go eng.mutation(chMutation, chFitness)
	go eng.fitness(chFitness, chIndividuals)
	go eng.offsprings(offspringSize, chIndividuals, chOffsprings)

	defer close(chErr)
	defer close(chSolution)

	// Run until an ending condition or an error is found
	go eng.run(start, population, chSelection, chOffsprings, chSolution)
	select {
	case sol := <-chSolution:
		return sol, nil
	case err := <-chErr:
		return Solution{}, err
	}
}

func (eng Engine) run(
	start time.Time,
	population gene.Population,
	chSelection chan<- gene.Population,
	chOffsprings <-chan gene.Population,
	chSolution chan<- Solution,
) {
	withBestIndividual := population
	withBestTotalFit := population

	for {
		// End ?
		termination := eng.Termination.End(population, withBestIndividual, withBestTotalFit)
		if termination != nil {
			chSolution <- Solution{
				PopWithBestIndividual:   withBestIndividual,
				PopWithBestTotalFitness: withBestTotalFit,
				Termination:             termination,
			}
			return
		}

		// Start selection process
		chSelection <- population

		// Wait for offspring to be ready
		offsprings := <-chOffsprings
		population = eng.survivors(start, population, offsprings)

		// Custom action
		if population.Stats.TotalFitness > withBestTotalFit.Stats.TotalFitness {
			withBestTotalFit = population
		}
		if population.Stats.Elite.Fitness > withBestIndividual.Stats.Elite.Fitness {
			withBestIndividual = population
		}
		eng.onNewGeneration(population, withBestIndividual, withBestTotalFit)
	}
}

// onNewGeneration calls the user method (only if defined)
func (eng Engine) onNewGeneration(population, withBestIndividual, withBestTotalFit gene.Population) {
	if eng.OnNewGeneration != nil {
		eng.OnNewGeneration(population, withBestIndividual, withBestTotalFit)
	}
}

// selection process: generate 1 selection per individual in the original population
func (eng Engine) selection(offspringSize int, in <-chan gene.Population, out chan<- gene.Chromosome, chErr chan<- error) {
	defer close(out)
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
	defer close(out)
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
	defer close(out)
	for mut := range in {
		if eng.Mutation != nil {
			mut = eng.Mutation.Mutate(mut)
		}
		out <- mut
	}
}

// fitness process: compute each individual fitness
func (eng Engine) fitness(in <-chan gene.Chromosome, out chan<- gene.Individual) {
	defer close(out)
	for chrm := range in {
		fitness := eng.Fitness(chrm)
		out <- gene.NewIndividual(chrm, fitness)
	}
}

// Group every n individuals into a new population
func (eng Engine) offsprings(offspringSize int, in <-chan gene.Individual, out chan<- gene.Population) {
	defer close(out)
	offsprings := gene.NewPopulation(offspringSize)
	var i int
	for ind := range in {
		offsprings.Individuals[i] = ind
		i++
		if i == offspringSize { // valid current offspring and begin next
			out <- offsprings
			offsprings = gene.NewPopulation(offspringSize)
			i = 0
		}
	}
}

// Survivors builds a new population of individuals
// The new population has changed, so compute global data like total fitness
func (eng Engine) survivors(start time.Time, parents gene.Population, offsprings gene.Population) gene.Population {
	newPop := eng.Survivor.Survive(parents, offsprings)
	newPop.ComputeTotalFitness()
	newPop.ComputeRank()
	newPop.Stats.GenerationNb = parents.Stats.GenerationNb + 1
	newPop.Stats.TotalDuration = time.Since(start)
	return newPop
}
