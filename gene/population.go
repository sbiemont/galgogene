package gene

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sbiemont/galgogene/random"
)

// Individual represents the coded chain of bases with a given fitness
type Individual struct {
	ID      uuid.UUID  // Unique identifier for the individual
	Code    Chromosome // Genetic data representation
	Fitness float64    // Current fitness of the individual
	Rank    int        // Generation number of the individual (starts at 0)
}

// NewIndividual initializes a new individual instance
func NewIndividual(code Chromosome, fitness float64) Individual {
	return Individual{
		ID:      uuid.New(),
		Code:    code,
		Fitness: fitness,
	}
}

// PopulationStats gathers general data for a population
type PopulationStats struct {
	TotalFitness  float64
	TotalDuration time.Duration
	GenerationNb  int
	Elite         Individual
}

// Fitness defines the fitness function for a given individual
type Fitness func(Chromosome) float64

// Population represents an ordered list of individual with a common fitness function
type Population struct {
	Individuals []Individual
	Stats       PopulationStats
}

// NewPopulation init an empty population of n individuals with a fitness function
func NewPopulation(size int) Population {
	return Population{
		Individuals: make([]Individual, size),
	}
}

// Init the population with random chromosome of the given size
func (pop *Population) Init(chrmSize int, initializer Initializer, fitness Fitness) error {
	// Full init
	for i := range pop.Individuals {
		chrm, err := initializer.Init(chrmSize)
		if err != nil {
			return err
		}

		// Update current individual
		pop.Individuals[i].Code = chrm
		pop.Individuals[i].Fitness = fitness(chrm)
	}

	pop.ComputeTotalFitness()
	return nil
}

// ComputeTotalFitness restart computation of total fitness
// Compute
//   - Total fitness
//   - Elite
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

// ComputeRank move all individual to the upper rank
func (pop *Population) ComputeRank() {
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

// Shuffle the population
func (pop Population) Shuffle() {
	random.Shuffle(pop.Individuals)
}

// Sort population by highest fitness first
func (pop Population) Elite() Individual {
	return pop.Stats.Elite
}

// First extracts k first Individuals of the current population
func (pop Population) First(k int) Population {
	return Population{
		Individuals: pop.Individuals[0:k],
	}
}

// Last extracts k last Individuals of the current population
func (pop Population) Last(k int) Population {
	return Population{
		Individuals: pop.Individuals[len(pop.Individuals)-k:],
	}
}

// Len returns the popultation number of individuals
func (pop Population) Len() int {
	return len(pop.Individuals)
}

// Unique returns a slice of unique individuals
func (pop Population) Unique() []Individual {
	// Extract unique individuals key signature
	unique := make(map[string]Individual)
	for _, individual := range pop.Individuals {
		key := individual.Code.String()
		unique[key] = individual
	}

	// Flatten individuals
	result := make([]Individual, 0, len(unique))
	for _, individual := range unique {
		result = append(result, individual)
	}
	return result
}
