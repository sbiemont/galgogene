package operator

import (
	"errors"

	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"
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
type MultiSurvivor []probaSurvivor

// Use the given probabilistic survivor
func (svr MultiSurvivor) Use(rate float64, survivor Survivor) MultiSurvivor {
	svr = append(svr, probaSurvivor{
		rate:     rate,
		survivor: survivor,
	})
	return svr
}

// Otherwise defines the survivor to be used if no survivor have been picked
func (svr MultiSurvivor) Otherwise(survivor Survivor) multiSurvivor {
	return multiSurvivor{
		survivors: svr,
		deflt:     survivor,
	}
}

// multiSurvivor defines a list of **ordered** surviving actions ending with a default one
type multiSurvivor struct {
	survivors []probaSurvivor
	deflt     Survivor
}

// Survive applies one of the defined surviors 
func (svr multiSurvivor) Survive(parents gene.Population, survivors *gene.Population) error {
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
