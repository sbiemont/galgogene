package gene

import (
	"errors"
	"sort"
	"time"
)

// Individual represents the coded chain of bits with a given fitness
type Individual struct {
	Code    Bits    // Genetic data representation
	Fitness float64 // Current fitness of the individual
	Rank    int     // Generation number of the individual (starts at 0)
}

// NewIndividual initializes a new individual instance
func NewIndividual(code Bits) Individual {
	return Individual{
		Code: code,
	}
}

// PopulationStats gathers general data for a population
type PopulationStats struct {
	TotalFitness  float64
	TotalDuration time.Duration
	GenerationNb  int
	Elite         Individual
}

// FitnessFct defines the fitness function for a given individual
type FitnessFct func(Bits) float64

// Population represents an ordered list of individual with a common fitness function
type Population struct {
	Individuals []Individual
	fitness     FitnessFct
	Stats       PopulationStats
	initializer Initializer
}

// NewPopulation init an empty population of n individuals with a fitness function
func NewPopulation(size int, fitness FitnessFct, initializer Initializer) Population {
	return Population{
		Individuals: make([]Individual, size),
		fitness:     fitness,
		initializer: initializer,
	}
}

// NewPopulationFrom init an empty population of n individuals with the fitness function of the specified population
func NewPopulationFrom(size int, pop Population) Population {
	return NewPopulation(size, pop.fitness, pop.initializer)
}

// Init the population with random bits of the given size
func (pop *Population) Init(bitsSize int) error {
	if pop.initializer == nil {
		return errors.New("initializer shall be set")
	}

	// Check only once
	err := pop.initializer.Check(bitsSize)
	if err != nil {
		return err
	}

	// Full init
	for i := range pop.Individuals {
		pop.Individuals[i].Code = pop.initializer.Init(bitsSize)
	}
	pop.ComputeFitness()
	return nil
}

// ComputeFitness computes an set all fitnesses for each individual
// Compute
//  * Individual fitness
//  * Total fitness
//  * Elite
func (pop *Population) ComputeFitness() {
	pop.Stats.Elite = pop.Individuals[0]
	for i, individual := range pop.Individuals {
		fitness := pop.fitness(individual.Code)
		pop.Individuals[i].Fitness = fitness
		pop.Stats.TotalFitness += fitness
		if individual.Fitness > pop.Stats.Elite.Fitness {
			pop.Stats.Elite = individual
		}
	}
}

// ComputeTotalFitness restart computation of total fitness
// Compute
//  * Total fitness
//  * Elite
func (pop *Population) ComputeTotalFitness() {
	pop.Stats.TotalFitness = 0
	pop.Stats.Elite = pop.Individuals[0]
	for _, individual := range pop.Individuals {
		pop.Stats.TotalFitness += individual.Fitness
		if individual.Fitness > pop.Stats.Elite.Fitness {
			pop.Stats.Elite = individual
		}
	}
}

// AddRank move all individual to the upper rank
func (pop *Population) AddRank() {
	for i := range pop.Individuals {
		pop.Individuals[i].Rank++
	}
}

// SortByFitness sorts the population by highest fitness first
func (pop Population) SortByFitness() {
	sort.Slice(pop.Individuals, func(i, j int) bool {
		return pop.Individuals[i].Fitness > pop.Individuals[j].Fitness
	})
}

// SortByRank sorts the population by newest individual first
func (pop Population) SortByRank() {
	sort.Slice(pop.Individuals, func(i, j int) bool {
		return pop.Individuals[i].Rank < pop.Individuals[j].Rank
	})
}

// Sort population by highest fitness first
func (pop Population) Elite() Individual {
	return pop.Stats.Elite
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

// Len returns the popultation number of individuals
func (pop Population) Len() int {
	return len(pop.Individuals)
}

// MapCount counts the number of different fitnesses
func (pop Population) MapCount() int {
	count := make(map[float64]interface{})
	for _, individual := range pop.Individuals {
		count[individual.Fitness] = nil
	}
	return len(count)
}
