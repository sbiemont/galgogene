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
	Convey("multi string matcher", t, func() {
		targetStr := "This is my first genetic algorithm using multi string matcher!"
		targetBits := gene.NewBitsFromBytes([]byte(targetStr))
		bitsSize := len(targetBits)

		// Fitness: match the input string bit by bit
		var fitness gene.FitnessFct = func(bits gene.Bits) float64 {
			var fitness float64
			for i, targetBit := range targetBits {
				if targetBit == bits[i] {
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
			Mutator: operator.MultiMutators{
				operator.NewProbaMutator(1, operator.UniformCrossOver{}),   // 100% chance to apply uniform cross-over
				operator.NewProbaMutator(0.05, operator.Mutate{Rate: 0.5}), // 5% chance to mutate (with 50% chance of changing each bits)
			},
			Survivor: operator.MultiSurvivor{
				operator.SurvivorAddParentsElite{K: 2}, // Add first 2 best individuals from parent's population (elitism is not so good for population diversity)
				operator.SurvivorElite{},               // Then, only keey k best individual in new generation
			},
			Ender: operator.MultiEnder{
				&operator.EnderGeneration{K: 100},              // End at generation #100
				&operator.EnderImprovement{},                   // End when total fitness cannot be improved
				&operator.EnderAboveFitness{Fitness: 1},        // End with perfect fitness
				&operator.EnderDuration{Duration: time.Second}, // End after 1s
			},
			OnNewGeneration: func(pop gene.Population) {
				elite := pop.Elite()
				fmt.Printf(
					"Generation #%d, fit: %f, tot: %f, str: %s\n",
					pop.Stats.GenerationNb, elite.Fitness, pop.Stats.TotalFitness, string(elite.Code.ToBytes()),
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
