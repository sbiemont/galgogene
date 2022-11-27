package example

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"galgogene.git/engine"
	"galgogene.git/gene"
	"galgogene.git/operator"
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
				Use(0.8, operator.UniformCrossOver{}),   // 80% chance to apply uniform crossover
			Mutation: operator.MultiMutation{}.
				Use(0.1, operator.UniqueMutation{}), // 10% chance to mutate one bit
			Survivor: operator.MultiSurvivor{}.
				Use(0.9, operator.EliteSurvivor{}).
				Otherwise(operator.ChildrenSurvivor{}),
			Termination: operator.MultiTermination{}.
				Use(&operator.GenerationTermination{K: 200}).              // End at generation #200
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
		last, best, termination, err := eng.Run(popSize, bitsSize, fitness)
		So(err, ShouldBeNil)
		So(termination, ShouldNotBeNil)
		So(last.Individuals, ShouldHaveLength, popSize)
		So(best.Individuals, ShouldHaveLength, popSize)

		// Print elite of best population
		elite := best.Elite()
		bytes, _ := szr.ToBytes(elite.Code)
		fmt.Printf(
			"Best #%d, fit: %f, tot: %f, str: %s\n",
			best.Stats.GenerationNb,
			best.Elite().Fitness,
			best.Stats.TotalFitness,
			string(bytes),
		)
	})
}
