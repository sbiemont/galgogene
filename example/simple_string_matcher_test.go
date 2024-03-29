package example

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/sbiemont/galgogene/engine"
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSimpleStringMatcher(t *testing.T) {
	szr, _ := gene.NewSerializer(8)

	Convey("simple string matcher bit by bit", t, func() {
		rand.Seed(2)
		targetStr := "This is my first genetic algorithm using simple string matcher!"
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

		// Engine will stop when max fitness is reached
		perfectFitness := &operator.FitnessTermination{Fitness: 1}
		eng := engine.Engine{
			Initializer: gene.NewRandomInitializer(1),
			Selection:   operator.RouletteSelection{},
			CrossOver:   operator.UniformCrossOver{},
			Survivor:    operator.EliteSurvivor{},
			Termination: perfectFitness,
			OnNewGeneration: func(pop gene.Population) {
				elite := pop.Elite()
				bytes, _ := szr.ToBytes(elite.Code)
				fmt.Printf(
					"Generation #%d, dur: %s fit: %f, tot: %f, str: %s\n",
					pop.Stats.GenerationNb,
					pop.Stats.TotalDuration,
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
		So(termination, ShouldEqual, perfectFitness)
		So(last.Individuals, ShouldHaveLength, popSize)
		So(best.Individuals, ShouldHaveLength, popSize)
	})
}
