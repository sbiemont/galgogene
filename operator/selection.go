package operator

import (
	"errors"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"
)

// selection: https://en.wikipedia.org/wiki/Selection_(genetic_algorithm)

// Selection defines the selection method of one individual in a population
type Selection interface {
	Select(pop gene.Population) (gene.Individual, error)
}

// ------------------------------

// RouletteSelection defines a fitness proportionate selection
// https://en.wikipedia.org/wiki/Fitness_proportionate_selection
type RouletteSelection struct{}

// Select 1 individual using roulette method
// Calculate S = the sum of all fitnesses
// Generate a random number between 0 and S
// Starting from the top of the population, keep adding the fitnesses to the partial sum P, till P<S
// The individual for which P exceeds S is the chosen individual.
func (RouletteSelection) Select(pop gene.Population) (gene.Individual, error) {
	randFitness := random.Percent() * pop.Stats.TotalFitness

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

// TournamentSelection select the best individual between k individuals
type TournamentSelection struct {
	Fighters int // Number of fighters
}

// Select 1 individual between k figthers
// Choose k individuals from the population and retrieves the best one
func (st TournamentSelection) Select(pop gene.Population) (gene.Individual, error) {
	if st.Fighters == 0 {
		return gene.Individual{}, errors.New("selection tournament: fighters shall be > 0")
	}

	// Select k indexes from the population
	indexes := random.OrderedInts(0, len(pop.Individuals), st.Fighters)

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

// EliteSelection selects the best individual from the population
type EliteSelection struct{}

func (EliteSelection) Select(pop gene.Population) (gene.Individual, error) {
	return pop.Elite(), nil
}

// ------------------------------

// probaSelection is a probabilistic selection
type probaSelection struct {
	rate float64
	sel  Selection
}

// MultiSelection defines an ordered list of selections each one with a given probability in [0 ; 1]
// The first chosen selection ends processing. If no selection matches, an error is raised
type MultiSelection []probaSelection

// Use the given proba selection
func (ms MultiSelection) Use(rate float64, selection Selection) MultiSelection {
	ms = append(ms, probaSelection{
		rate: rate,
		sel:  selection,
	})
	return ms
}

// Otherwise defines the selection to be used if no selection have been picked
func (ms MultiSelection) Otherwise(selection Selection) multiSelection {
	return multiSelection{
		selections: ms,
		deflt:      selection,
	}
}

// multiSelection ends the selection with a default behavior
type multiSelection struct {
	selections []probaSelection
	deflt      Selection
}

// Select an individual
// First, randomly choose a selection
// Then, use the chosen selection on the current population
func (ms multiSelection) Select(pop gene.Population) (gene.Individual, error) {
	if ms.deflt == nil {
		return gene.Individual{}, errors.New("no default selector defined")
	}

	// Find for first selector to be used
	for _, proba := range ms.selections {
		if random.Peek(proba.rate) {
			return proba.sel.Select(pop)
		}
	}

	// Use default selector
	return ms.deflt.Select(pop)
}
