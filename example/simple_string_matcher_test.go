package example

import (
	"fmt"
	"math/rand"
	"testing"

	"genalgo.git/engine"
	"genalgo.git/gene"
	"genalgo.git/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func bitsToString(bits gene.Bits) string {
	bytes, _ := bits.GroupBitsBy(8)
	return string(bytes)
}

func TestSimpleStringMatcher(t *testing.T) {
	Convey("simple string matcher bit by bit", t, func() {
		rand.Seed(180)
		targetStr := "This is my first genetic algorithm using simple string matcher!"
		targetBits := gene.NewBitsFromBytes([]byte(targetStr))
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

		// Engine will stop when max fitness is reached
		perfectFitness := &operator.EnderAboveFitness{Fitness: 1.0}
		eng := engine.Engine{
			Selector: operator.SelectorRoulette{},
			Mutator:  operator.UniformCrossOver{},
			Survivor: operator.SurvivorElite{},
			Ender:    perfectFitness,
			OnNewGeneration: func(pop gene.Population) {
				elite := pop.Elite()
				fmt.Printf(
					"Generation #%d, dur: %s fit: %f, tot: %f, str: %s\n",
					pop.Stats.GenerationNb,
					pop.Stats.TotalDuration,
					elite.Fitness,
					pop.Stats.TotalFitness,
					bitsToString(elite.Code),
				)
			},
		}

		// Run and check output
		popSize := 100
		last, ender, err := eng.Run(popSize, bitsSize, fitness)
		So(err, ShouldBeNil)
		So(ender, ShouldEqual, perfectFitness)
		So(last.Individuals, ShouldHaveLength, popSize)
	})
}
