package operator

import (
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/random"
)

// Survivor defines an action to be applied on the current generation
type Survivor interface {
	// Survive allow to choose some individual from the parents population and/or update the survivors
	Survive(parents gene.Population, offsprings gene.Population) gene.Population
}

// mergePopulations creates a new population with all individuals of both populations but no stats
func mergePopulations(pop1, pop2 gene.Population) gene.Population {
	return gene.Population{
		Individuals: append(pop1.Individuals, pop2.Individuals...),
	}
}

// ------------------------------

// EliteSurvivor selects the elite from the parents + children population
type EliteSurvivor struct{}

func (svr EliteSurvivor) Survive(parents gene.Population, offsprings gene.Population) gene.Population {
	survivors := mergePopulations(parents, offsprings)
	survivors.SortByFitness()
	return survivors.First(parents.Len())
}

// ------------------------------

// RankSurvivor selects the newer individuals from the parents + children population
type RankSurvivor struct{}

func (svr RankSurvivor) Survive(parents gene.Population, offsprings gene.Population) gene.Population {
	survivors := mergePopulations(parents, offsprings)
	survivors.SortByRank()
	return survivors.First(parents.Len())
}

// ------------------------------

// RandomSurvivor selects purely random survivors in the parents + children population
type RandomSurvivor struct{}

func (svr RandomSurvivor) Survive(parents gene.Population, offsprings gene.Population) gene.Population {
	survivors := mergePopulations(parents, offsprings)
	survivors.Shuffle()
	return survivors.First(parents.Len())
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
	return append(svr, probaSurvivor{
		rate:     rate,
		survivor: survivor,
	})
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
func (svr multiSurvivor) Survive(parents gene.Population, offsprings gene.Population) gene.Population {
	// Run first survivor
	for _, proba := range svr.survivors {
		if random.Peek(proba.rate) {
			return proba.survivor.Survive(parents, offsprings)
		}
	}

	// Otherwise, use default survivor
	return svr.deflt.Survive(parents, offsprings)
}
