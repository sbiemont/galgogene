package engine

import (
	"fmt"
	"testing"
	"time"

	"genalgo.git/gene"
	"genalgo.git/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEngine(t *testing.T) {
	Convey("new", t, func() {
		Convey("new engine", func() {
			targetStr := "This is my first genetic algorithm!"
			targetBits := gene.NewBitsFromBytes([]byte(targetStr))
			bitsSize := len(targetBits)

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

			// engine := Engine{
			// 	selector: operator.MultiSelectors{
			// 		operator.NewProbaSelector(0.5, operator.SelectorRoulette{}),
			// 		operator.NewProbaSelector(1, operator.SelectorTournament{Fighters: 3}),
			// 	},
			// 	mutator: operator.MultiMutators{
			// 		operator.NewProbaMutator(1, operator.UniformCrossOver{}),
			// 		operator.NewProbaMutator(0.1, operator.Mutate{Rate: 0.5}),
			// 	},
			// 	survivor: operator.MultiSurvivor{
			// 		operator.SurvivorAddParentsElite{K: 10},
			// 		operator.SurvivorElite{K: popSize},
			// 	},
			// 	ender: operator.MultiEnder{
			// 		&operator.EnderGeneration{K: 50},
			// 		&operator.EnderImprovement{},
			// 		&operator.EnderAboveFitness{Fitness: 1},
			// 		&operator.EnderDuration{Duration: time.Second},
			// 	},
			// 	OnNewPopulation: func(pop gene.Population) {
			// 		elite := pop.Elite()
			// 		fmt.Printf(
			// 			"Generation #%d, fit: %f, tot: %f, str: %s\n",
			// 			pop.NbGeneration, elite.Fitness, pop.TotalFitness, string(elite.Code.ToBytes()),
			// 		)
			// 	},
			// }

			engine := Engine{
				Selector: operator.SelectorRoulette{},
				Mutator:  operator.UniformCrossOver{},
				Survivor: operator.SurvivorElite{K: popSize},
				Ender:    &operator.EnderDuration{Duration: 100 * time.Millisecond},
				OnNewGeneration: func(pop gene.Population) {
					elite := pop.Elite()
					fmt.Printf(
						"Generation #%d, dur: %s fit: %f, tot: %f, str: %s\n",
						pop.GenerationNb, pop.TotalDuration, elite.Fitness, pop.TotalFitness, string(elite.Code.ToBytes()),
					)
				},
			}

			last, ender, err := engine.Run(popSize, bitsSize, fitness)
			engine.OnNewGeneration(last)
			So(err, ShouldBeNil)
			So(ender, ShouldNotBeNil)
		})
	})
}
