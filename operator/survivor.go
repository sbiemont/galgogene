package operator

import (
	"errors"

	"galgogene.git/gene"
	"galgogene.git/random"
)

// Survivor defines an action to be applied on the current generation
type Survivor interface {
	// Survive allow to choose some individual from the parents population and/or update the survivors
	Survive(parents gene.Population, survivors *gene.Population) error
}

// ------------------------------

// EliteSurvivor selects the elite from the parents + children population
type EliteSurvivor struct{}

func (svr EliteSurvivor) Survive(parents gene.Population, survivors *gene.Population) error {
	survivors.Individuals = append(survivors.Individuals, parents.Individuals...)
	survivors.SortByFitness()
	survivors.Individuals = (*survivors).First(parents.Len()).Individuals
	return nil
}

// ------------------------------

// RankSurvivor selects the newer individuals from the parents + children population
type RankSurvivor struct{}

func (svr RankSurvivor) Survive(parents gene.Population, survivors *gene.Population) error {
	survivors.Individuals = append(survivors.Individuals, parents.Individuals...)
	survivors.SortByRank()
	survivors.Individuals = (*survivors).First(parents.Len()).Individuals
	return nil
}

// ------------------------------

// ChildrenSurvivor only let the children population survive
type ChildrenSurvivor struct{}

func (svr ChildrenSurvivor) Survive(parents gene.Population, survivors *gene.Population) error {
	survivors.Individuals = (*survivors).First(parents.Len()).Individuals
	return nil
}

// ------------------------------

type probaSurvivor struct {
	rate     float64
	survivor Survivor
}

// MultiSurvivor defines a list of **ordered** surviving actions
type MultiSurvivor struct {
	survivors []probaSurvivor
	deflt     Survivor
}

// Use the given probabilistic survivor
func (svr MultiSurvivor) Use(rate float64, survivor Survivor) MultiSurvivor {
	svr.survivors = append(svr.survivors, probaSurvivor{
		rate:     rate,
		survivor: survivor,
	})
	return svr
}

// Otherwise defines the survivor to be used if no survivor have been picked
func (svr MultiSurvivor) Otherwise(survivor Survivor) MultiSurvivor {
	svr.deflt = survivor
	return svr
}

func (svr MultiSurvivor) Survive(parents gene.Population, survivors *gene.Population) error {
	// Check
	if svr.deflt == nil {
		return errors.New("no default survivor defined")
	}

	// Run first survivor
	for _, proba := range svr.survivors {
		if random.Peek(proba.rate) {
			return proba.survivor.Survive(parents, survivors)
		}
	}

	// Otherwise, use default survivor
	return svr.deflt.Survive(parents, survivors)
}
