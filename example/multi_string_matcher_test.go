package example

import (
	"fmt"
	"testing"
	"time"

	"genalgo.git/engine"
	"genalgo.git/gene"
	"genalgo.git/operator"
	. "github.com/smartystreets/goconvey/convey"
)

// Example with multi criteria for
func TestStringMatcher(t *testing.T) {
	szr, _ := gene.NewSerializer(8)

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

		popSize := 100

		eng := engine.Engine{
			Selection: operator.MultiSelection{
				operator.NewProbaSelection(0.5, operator.SelectionRoulette{}),            // 50% chance to get roulette
				operator.NewProbaSelection(1, operator.SelectionTournament{Fighters: 2}), // otherwise, use tournament
			},
			CrossOver: operator.MultiCrossOver{
				operator.NewProbaCrossOver(0.5, operator.UniformCrossOver{}),   // 50% chance to apply uniform cross-over
				operator.NewProbaCrossOver(0.9, operator.TwoPointsCrossOver{}), // 50% chance to apply uniform cross-over
			},
			Mutation: operator.MultiMutation{
				operator.NewProbaMutation(0.05, operator.UniformMutation{}), // 5% chance to mutate (with 50% chance of changing each bits)
			},
			Survivor: operator.MultiSurvivor{
				operator.SurvivorAddAllParents{}, // Add first all parents to the new generation pool
				operator.SurvivorElite{},         // Then, only keey k best individual in new generation
			},
			Termination: operator.MultiTermination{
				&operator.TerminationGeneration{K: 100},              // End at generation #100
				&operator.TerminationImprovement{},                   // End when total fitness cannot be improved
				&operator.TerminationAboveFitness{Fitness: 1},        // End with perfect fitness
				&operator.TerminationDuration{Duration: time.Second}, // End after 1s
			},
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
		last, termination, err := eng.Run(popSize, bitsSize, fitness)
		So(err, ShouldBeNil)
		So(termination, ShouldNotBeNil)
		So(last.Individuals, ShouldHaveLength, popSize)
	})
}
