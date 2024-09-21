package main

import (
	"fmt"

	"github.com/sbiemont/galgogene/engine"
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

func main() {
	targetStr := "This is my first string matcher!"

	// Fitness: match the input string char by char
	var fitness gene.Fitness = func(chr gene.Chromosome) float64 {
		chrStr := string(chr.Raw)

		var fitness float64
		for i := range len(targetStr) {
			if targetStr[i] == chrStr[i] {
				fitness += 1
			}
		}
		return fitness / float64(chr.Len())
	}

	// Engine will stop when max fitness is reached
	eng := engine.Engine{
		Initializer: gene.NewRandomInitializer(255),
		Selection:   operator.TournamentSelection{Fighters: 6},
		CrossOver:   operator.ThreePointsCrossOver{},
		Mutation:    operator.UniqueMutation{},
		Survivor:    operator.EliteSurvivor{},
		Termination: &operator.FitnessTermination{Fitness: 1},
		Fitness:     fitness,
		OnNewGeneration: func(pop gene.Population) {
			elite := pop.Elite()
			fmt.Printf(
				"Generation #%d, dur: %11s fit: %f, tot: %f, str: %s\n",
				pop.Stats.GenerationNb,
				pop.Stats.TotalDuration,
				elite.Fitness,
				pop.Stats.TotalFitness,
				string(elite.Code.Raw),
			)
		},
	}

	// Run and check output
	popSize := 50
	_, err := eng.Run(popSize, 50, len(targetStr))
	if err != nil {
		panic(err)
	}
}
