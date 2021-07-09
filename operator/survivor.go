package operator

import "genalgo.git/gene"

// Survivor defines an action to be applied on the current generation
type Survivor interface {
	// Survive allow to choose some individual from the parents population and/or update the survivors
	Survive(parents gene.Population, survivors *gene.Population)
}

// SurvivorAddAllParents adds all parents individuals to the surviving population
type SurvivorAddAllParents struct{}

func (svr SurvivorAddAllParents) Survive(parents gene.Population, survivors *gene.Population) {
	survivors.Individuals = append(survivors.Individuals, parents.Individuals...)
}

// SurvivorAddParentsElite adds the K parents elite to the surviving population
type SurvivorAddParentsElite struct {
	K int // Number of individuals
}

func (svr SurvivorAddParentsElite) Survive(parents gene.Population, survivors *gene.Population) {
	parents.Sort()
	elite := parents.First(svr.K)
	survivors.Individuals = append(survivors.Individuals, elite.Individuals...)
}

// SurvivorElite only selects the K elite from the surviving population
// If K is 0, the same population size is kept
type SurvivorElite struct {
	K int // Number of individuals
}

func (svr SurvivorElite) Survive(parents gene.Population, survivors *gene.Population) {
	k := svr.K
	if k == 0 {
		k = len(parents.Individuals) // if K is 0, keep the same number of initial population
	}
	survivors.Sort()
	tmp := (*survivors).First(k)
	survivors.Individuals = tmp.Individuals
}

// MultiSurvivor defines a list of **ordered** surviving actions
type MultiSurvivor []Survivor

func (svr MultiSurvivor) Survive(parents gene.Population, survivors *gene.Population) {
	for _, s := range svr {
		s.Survive(parents, survivors)
	}
}
