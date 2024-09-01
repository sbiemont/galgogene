package main

import (
	"fmt"
	"math/rand"

	"github.com/sbiemont/galgogene/engine"
	"github.com/sbiemont/galgogene/gene"
	"github.com/sbiemont/galgogene/operator"
)

func main() {
	szr, _ := gene.NewSerializer(8)

	rand.New(rand.NewSource(2))
	targetStr := "This is my first string matcher!"
	// targetStr := "Hello world!"
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
	_, err := eng.Run(popSize, bitsSize, fitness)
	if err != nil {
		panic(err)
	}
}
