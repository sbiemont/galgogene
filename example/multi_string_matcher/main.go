package main

import (
	"fmt"
	"time"

	"github.com/sbiemont/galgogene/engine"
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

// Example with multi criteria for string matcher
func main() {
	targetStr := "This is my first genetic algorithm using multi string matcher with a custom engine!"

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

	eng := engine.Engine{
		Initializer: gene.NewRandomInitializer(255),
		Selection: operator.MultiSelection{}.
			Use(0.5, operator.TournamentSelection{Fighters: 10}). // 50% chance to use tournament
			Use(0.5, operator.EliteSelection{}).                  // 50% chance to use elite
			Otherwise(operator.RouletteSelection{}),              // otherwise, use roulette
		CrossOver: operator.MultiCrossOver{}.
			Use(0.1, operator.OnePointCrossOver{}). // 10% chance to apply 1 point crossover
			Use(0.5, operator.UniformCrossOver{}),  // 50% chance to apply uniform crossover
		Mutation: operator.MultiMutation{}.
			Use(0.95, operator.UniqueMutation{}).  // 95% chance to mutate one bit
			Use(0.05, operator.UniformMutation{}), // 5% chance to mutate using uniform transformation
		Survivor: operator.MultiSurvivor{}.
			Use(0.01, operator.RankSurvivor{}).  // 1% chance to use rank survivor (newest individuals)
			Otherwise(operator.EliteSurvivor{}), // Otherwise, select elite individuals
		Termination: operator.MultiTermination{}.
			Use(&operator.GenerationTermination{K: 300}).                  // End at generation #300
			Use(&operator.FitnessTermination{Fitness: 1}).                 // End with perfect fitness
			Use(&operator.DurationTermination{Duration: 3 * time.Second}), // End after 3s
		Fitness: fitness,
		OnNewGeneration: func(pop gene.Population) {
			elite := pop.Elite()
			fmt.Printf(
				"Generation #%d, dur: %4dms fit: %f, tot: %f, str: %s\n",
				pop.Stats.GenerationNb,
				pop.Stats.TotalDuration.Milliseconds(),
				elite.Fitness,
				pop.Stats.TotalFitness,
				string(elite.Code.Raw),
			)
		},
	}

	// Run and check output
	popSize := 300
	sol, err := eng.Run(popSize, 2*popSize, len(targetStr))
	if err != nil {
		panic(err)
	}

	_, ok := sol.Termination.(*operator.FitnessTermination)
	if ok {
		fmt.Println("\nsuccess")
	} else {
		fmt.Println("\nfailure")
	}
}
