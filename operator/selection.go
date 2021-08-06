package operator

import (
	"errors"
	"math/rand"

	"genalgo.git/gene"
	"genalgo.git/random"
)

// selection: "https://en.wikipedia.org/wiki/Selection_(genetic_algorithm)"

// Selection defines the selection method of one individual in a population
type Selection interface {
	Select(pop gene.Population) (gene.Individual, error)
}

// ------------------------------

// SelectionRoulette defines a fitness proportionate selection
// https://en.wikipedia.org/wiki/Fitness_proportionate_selection
type SelectionRoulette struct{}

// Select 1 individual using roulette method
// Calculate S = the sum of all fitnesses
// Generate a random number between 0 and S
// Starting from the top of the population, keep adding the fitnesses to the partial sum P, till P<S
// The individual for which P exceeds S is the chosen individual.
func (SelectionRoulette) Select(pop gene.Population) (gene.Individual, error) {
	randFitness := rand.Float64() * pop.Stats.TotalFitness

	var currFitness float64
	for _, individual := range pop.Individuals {
		currFitness += individual.Fitness
		if currFitness >= randFitness {
			return individual, nil
		}
	}

	// Failed to peek an individual
	return gene.Individual{}, errors.New("selection roulette failed")
}

// ------------------------------

// SelectionTournament select the best individual between k individuals
type SelectionTournament struct {
	Fighters int // Number of fighters
}

// Select 1 individual between k figthers
// Choose k individuals from the population and retrieves the best one
func (st SelectionTournament) Select(pop gene.Population) (gene.Individual, error) {
	if st.Fighters == 0 {
		return gene.Individual{}, errors.New("selection tournament: fighters shall be > 0")
	}

	// Select k indexes from the population
	indexes := random.Ints(0, len(pop.Individuals), st.Fighters)

	// Select the best of chosen ones
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

type SelectionElite struct{}

func (SelectionElite) Select(pop gene.Population) (gene.Individual, error) {
	return pop.Elite(), nil
}

// ------------------------------

type ProbaSelection struct {
	rate float64
	sel  Selection
}

func NewProbaSelection(rate float64, sel Selection) ProbaSelection {
	return ProbaSelection{
		rate: rate,
		sel:  sel,
	}
}

// MultiSelection defines an ordered list of selections each one with a given probability in [0 ; 1]
// The first chosen selection ends processing. If no selection matches, an error is raised
type MultiSelection []ProbaSelection

// NewMultiSelection checks for consistency probailities
func NewMultiSelection(selections []ProbaSelection) (MultiSelection, error) {
	n := len(selections)
	if n == 0 {
		return nil, errors.New("at least one selection is required")
	}
	var sum float64

	// Check proba = 1 shall only be the last one
	for i, selection := range selections {
		sum += selection.rate
		if (selection.rate >= 1.0 && i != n-1) || // proba 1 but not last
			(selection.rate < 1.0 && i == n-1) { // proba not 1 but last
			return nil, errors.New("selection with proba=1 shall only be the last one")
		}
	}

	// OK
	return selections, nil
}

// Select an individual
// First, randomly choose a selection
// Then, use the chosen selection on the current population
func (ms MultiSelection) Select(pop gene.Population) (gene.Individual, error) {
	for _, probaSelection := range ms {
		if random.Peek(probaSelection.rate) {
			return probaSelection.sel.Select(pop)
		}
	}

	// Error, a selection with the max probability shall be defined to avoid this kind of error
	return gene.Individual{}, errors.New("selection multi proba, cannot peek any individual")
}
