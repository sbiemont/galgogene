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
			Selector: operator.MultiSelector{
				operator.NewProbaSelector(0.5, operator.SelectorRoulette{}),            // 50% chance to get roulette
				operator.NewProbaSelector(1, operator.SelectorTournament{Fighters: 2}), // otherwise, use tournament
			},
			Mutator: operator.MultiMutator{
				operator.NewProbaMutator(1, operator.UniformCrossOver{}),   // 100% chance to apply uniform cross-over
				operator.NewProbaMutator(0.05, operator.Mutate{Rate: 0.5}), // 5% chance to mutate (with 50% chance of changing each bits)
			},
			Survivor: operator.MultiSurvivor{
				operator.SurvivorAddAllParents{}, // Add first all parents to the new generation pool
				operator.SurvivorElite{},         // Then, only keey k best individual in new generation
			},
			Ender: operator.MultiEnder{
				&operator.EnderGeneration{K: 100},              // End at generation #100
				&operator.EnderImprovement{},                   // End when total fitness cannot be improved
				&operator.EnderAboveFitness{Fitness: 1},        // End with perfect fitness
				&operator.EnderDuration{Duration: time.Second}, // End after 1s
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
		last, ender, err := eng.Run(popSize, bitsSize, fitness)
		So(err, ShouldBeNil)
		So(ender, ShouldNotBeNil)
		So(last.Individuals, ShouldHaveLength, popSize)
	})
}
