package operator

import (
	"errors"
	"math/rand"

	"genalgo.git/gene"
	"genalgo.git/random"
)

// selection: "https://en.wikipedia.org/wiki/Selection_(genetic_algorithm)"

// Selector defines the selection method of one individual in a population
type Selector interface {
	Select(pop gene.Population) (gene.Individual, error)
}

// ------------------------------

// SelectorRoulette defines a fitness proportionate selection
// https://en.wikipedia.org/wiki/Fitness_proportionate_selection
type SelectorRoulette struct{}

// Select 1 individual using roulette method
// Calculate S = the sum of all fitnesses
// Generate a random number between 0 and S
// Starting from the top of the population, keep adding the fitnesses to the partial sum P, till P<S
// The individual for which P exceeds S is the chosen individual.
func (SelectorRoulette) Select(pop gene.Population) (gene.Individual, error) {
	randFitness := rand.Float64() * pop.Stats.TotalFitness

	var currFitness float64
	for _, individual := range pop.Individuals {
		currFitness += individual.Fitness
		if currFitness >= randFitness {
			return individual, nil
		}
	}

	// Failed to peek an individual
	return gene.Individual{}, errors.New("selector roulette failed")
}

// ------------------------------

// SelectorTournament select the best individual between k individuals
type SelectorTournament struct {
	Fighters int // Number of fighters
}

// Select 1 individual between k figthers
// Choose k individuals from the population and retrieves the best one
func (st SelectorTournament) Select(pop gene.Population) (gene.Individual, error) {
	if st.Fighters == 0 {
		return gene.Individual{}, errors.New("selector tournament: fighters shall be > 0")
	}

	// Select k indexes from the population
	indexes := random.Ints(0, len(pop.Individuals), st.Fighters)

	// Select the best of choosen ones
	best := &pop.Individuals[indexes[0]]
	for _, index := range indexes[1:] {
		current := &pop.Individuals[index]
		if current.Fitness > best.Fitness {
			best = current
		}
	}
	return *best, nil
}

// ------------------------------

type ProbaSelector struct {
	rate float64
	sel  Selector
}

func NewProbaSelector(rate float64, sel Selector) ProbaSelector {
	return ProbaSelector{
		rate: rate,
		sel:  sel,
	}
}

// MultiSelector defines an ordered list of selectors each one with a given probability in [0 ; 1]
// The first choosen selector ends processing. If no selector matches, an error is raised
type MultiSelector []ProbaSelector

// NewMultiSelector checks for consistency probailities
func NewMultiSelector(selectors []ProbaSelector) (MultiSelector, error) {
	n := len(selectors)
	if n == 0 {
		return nil, errors.New("at least one selector is required")
	}
	var sum float64

	// Check proba = 1 shall only be the last one
	for i, selector := range selectors {
		sum += selector.rate
		if (selector.rate >= 1.0 && i != n-1) || // proba 1 but not last
			(selector.rate < 1.0 && i == n-1) { // proba not 1 but last
			return nil, errors.New("selector with proba=1 shall only be the last one")
		}
	}

	// OK
	return selectors, nil
}

// Select an individual
// First, randomly choose a selector
// Then, use the choosen selector on the current population
func (ms MultiSelector) Select(pop gene.Population) (gene.Individual, error) {
	for _, probaSelector := range ms {
		if random.Peek(probaSelector.rate) {
			return probaSelector.sel.Select(pop)
		}
	}

	// Error, a selector with the max probability shall be defined to avoid this kind of error
	return gene.Individual{}, errors.New("selector multi proba, cannot peek any individual")
}
