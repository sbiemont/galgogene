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
	szr, _ := gene.NewSerializer(8)

	targetStr := "This is my first genetic algorithm using multi string matcher with a custom engine!"
	targetBits := szr.ToBits([]byte(targetStr))
	bitsSize := targetBits.Len()

	// Fitness: match the input string bit by bit
	var fitness gene.FitnessFct = func(bits gene.Bits) float64 {
		var fitness float64
		for i, targetBit := range targetBits.Raw {
			if targetBit == bits.Raw[i] {
				fitness += 1
			}
		}
		return fitness / float64(bitsSize)
	}

	eng := engine.Engine{
		Initializer: gene.NewRandomInitializer(1),
		Selection: operator.MultiSelection{}.
			Use(0.005, operator.TournamentSelection{Fighters: 2}). // 0.5% chance to use tournament
			Use(0.001, operator.EliteSelection{}).                 // 0.1% chance to use elite
			Otherwise(operator.RouletteSelection{}),               // otherwise, use roulette
		CrossOver: operator.MultiCrossOver{}.
			Use(0.1, operator.OnePointCrossOver{}). // 10% chance to apply 1 point crossover
			Use(0.5, operator.UniformCrossOver{}),  // 50% chance to apply uniform crossover
		Mutation: operator.MultiMutation{}.
			Use(0.05, operator.UniqueMutation{}).  // 5% chance to mutate one bit
			Use(0.05, operator.UniformMutation{}), // 5% chance to mutate using uniform transformation
		Survivor: operator.MultiSurvivor{}.
			Use(0.1, operator.RankSurvivor{}).      // 10% chance to use rank survivor
			Use(0.6, operator.EliteSurvivor{}).     // 60 % chance to select elite individuals
			Otherwise(operator.ChildrenSurvivor{}), // Otherwise, select new individuals
		Termination: operator.MultiTermination{}.
			Use(&operator.GenerationTermination{K: 300}).                  // End at generation #300
			Use(&operator.FitnessTermination{Fitness: 1}).                 // End with perfect fitness
			Use(&operator.DurationTermination{Duration: 3 * time.Second}), // End after 3s
		OnNewGeneration: func(pop gene.Population) {
			elite := pop.Elite()
			bytes, _ := szr.ToBytes(elite.Code)
			fmt.Printf(
				"Generation #%d, dur: %4dms fit: %f, tot: %f, str: %s\n",
				pop.Stats.GenerationNb,
				pop.Stats.TotalDuration.Milliseconds(),
				elite.Fitness,
				pop.Stats.TotalFitness,
				string(bytes),
			)
		},
	}

	// Run and check output
	popSize := 300
	_, err := eng.Run(popSize, bitsSize, fitness)
	if err != nil {
		panic(err)
	}
}
