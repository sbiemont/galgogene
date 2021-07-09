package gene

import (
	"sort"
	"time"
)

// Individual represents the coded chain of bits with a given fitness
type Individual struct {
	Code    Bits
	Fitness float64
}

// NewIndividual initializes a new individual instance
func NewIndividual(code Bits) Individual {
	return Individual{
		Code: code,
	}
}

// FitnessFct defines the fitness function for a given individual
type FitnessFct func(Bits) float64

// Population represents an ordered list of individual with a common fitness function
type Population struct {
	Individuals   []Individual
	fitness       FitnessFct
	TotalFitness  float64
	TotalDuration time.Duration
	GenerationNb  int
}

// NewPopulation init an empty population of n individuals with a fitness function
func NewPopulation(size int, fitness FitnessFct) Population {
	return Population{
		Individuals: make([]Individual, size),
		fitness:     fitness,
	}
}

// NewPopulationFrom init an empty population of n individuals with the fitness function of the specified population
func NewPopulationFrom(size int, pop Population) Population {
	return NewPopulation(size, pop.fitness)
}

// Init the population with random bits of the given size
func (pop *Population) Init(bitsSize int) {
	for i := range pop.Individuals {
		pop.Individuals[i].Code = NewBitsRandom(bitsSize)
	}
	pop.ComputeFitness()
}

// ComputeFitness computes an set all fitnesses for each individual
func (pop *Population) ComputeFitness() {
	pop.TotalFitness = 0
	for i, individual := range pop.Individuals {
		fitness := pop.fitness(individual.Code)
		pop.Individuals[i].Fitness = fitness
		pop.TotalFitness += fitness
	}
}

// ComputeTotalFitness restart computation of total fitness
func (pop *Population) ComputeTotalFitness() {
	pop.TotalFitness = 0
	for _, individual := range pop.Individuals {
		pop.TotalFitness += individual.Fitness
	}
}

// Sort population by highest fitness first
func (pop Population) Sort() {
	sort.Slice(pop.Individuals, func(i, j int) bool {
		return pop.Individuals[i].Fitness > pop.Individuals[j].Fitness
	})
}

// Sort population by highest fitness first
func (pop Population) Elite() Individual {
	var elite Individual = pop.Individuals[0]
	for _, invididual := range pop.Individuals {
		if invididual.Fitness > elite.Fitness {
			elite = invididual
		}
	}
	return elite
}

// First extracts k first Individuals of the current population
func (pop Population) First(k int) Population {
	return Population{
		Individuals: pop.Individuals[0:k],
		fitness:     pop.fitness,
	}
}

// Last extracts k last Individuals of the current population
func (pop Population) Last(k int) Population {
	return Population{
		Individuals: pop.Individuals[len(pop.Individuals)-k:],
		fitness:     pop.fitness,
	}
}
