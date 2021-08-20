package example

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"genalgo.git/engine"
	"genalgo.git/gene"
	"genalgo.git/operator"
	. "github.com/smartystreets/goconvey/convey"
)

// Example with multi criteria for string matcher
func TestStringMatcher(t *testing.T) {
	szr, _ := gene.NewSerializer(8)
	rand.Seed(time.Now().Unix())

	Convey("multi string matcher", t, func() {
		targetStr := "This is my first genetic algorithm using multi string matcher!"
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
				Use(0.5, operator.TournamentSelection{Fighters: 2}). // 50% chance to use tournament
				Otherwise(operator.RouletteSelection{}),             // otherwise, use roulette
			CrossOver: operator.MultiCrossOver{}.
				Use(0.1, operator.TwoPointsCrossOver{}). // 10% chance to apply 2 points crossover
				Use(1, operator.UniformCrossOver{}),     // 100% chance to apply uniform crossover
			Mutation: operator.MultiMutation{}.
				Use(0.05, operator.UniformMutation{}), // 5% chance to mutate (with 50% chance of changing each bits)
			Survivor: operator.MultiSurvivor{}.
				Use(0.9, operator.SurvivorElite{}).
				Otherwise(operator.SurvivorChildren{}),
			Termination: operator.MultiTermination{}.
				Use(&operator.GenerationTermination{K: 100}).              // End at generation #100
				Use(&operator.ImprovementTermination{}).                   // End when total fitness cannot be improved
				Use(&operator.FitnessTermination{Fitness: 1}).             // End with perfect fitness
				Use(&operator.DurationTermination{Duration: time.Second}), // End after 1s
			OnNewGeneration: func(pop gene.Population) {
				elite := pop.Elite()
				bytes, _ := szr.ToBytes(elite.Code)
				fmt.Printf(
					"Generation #%d, fit: %f, tot: %f, str: %s\n",
					pop.Stats.GenerationNb,
					elite.Fitness,
					pop.Stats.TotalFitness,
					string(bytes),
				)
			},
		}

		// Run and check output
		popSize := 100
		last, termination, err := eng.Run(popSize, bitsSize, fitness)
		So(err, ShouldBeNil)
		So(termination, ShouldNotBeNil)
		So(last.Individuals, ShouldHaveLength, popSize)
	})
}
